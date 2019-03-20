/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package authx

import (
	"fmt"
	"github.com/nalej/authx/internal/app/authx/handler"
	"github.com/nalej/authx/internal/app/authx/manager"
	"github.com/nalej/authx/internal/app/authx/providers/credentials"
	"github.com/nalej/authx/internal/app/authx/providers/device"
	"github.com/nalej/authx/internal/app/authx/providers/device_token"
	"github.com/nalej/authx/internal/app/authx/providers/token"
	"github.com/nalej/authx/internal/app/authx/providers/role"
	pbAuthx "github.com/nalej/grpc-authx-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

// Service is the Authx service instance.
type Service struct {
	// Config is the required parameters.
	Config
}

type Providers struct {
	roleProvider role.Role
	credProvider credentials.BasicCredentials
	devProvider device.Provider
	tokenProvider token.Token
	devTokenProvider device_token.Provider
}

type TokenManagers struct {
	tokenManager manager.Token
	deviceTokenManager manager.DeviceToken
}

// NewService create a new service instance.
func NewService(config Config) *Service {
	return &Service{Config: config}
}

func (s *Service) CreateInMemoryProviders() * Providers {
	return &Providers {
		roleProvider: role.NewRoleMockup(),
		credProvider: credentials.NewBasicCredentialMockup(),
		devProvider: device.NewMockupDeviceCredentialsProvider(),
		tokenProvider: token.NewTokenMockup(),
		devTokenProvider: device_token.NewDeviceTokenMockup(),
	}
}

func (s *Service) CreateDBScyllaProviders() * Providers {
	return &Providers {
		roleProvider: role.NewScyllaRoleProvider(
			s.Config.ScyllaDBAddress, s.Config.ScyllaDBPort, s.Config.KeySpace),
		credProvider: credentials.NewScyllaCredentialsProvider(
			s.Config.ScyllaDBAddress, s.Config.ScyllaDBPort, s.Config.KeySpace),
		devProvider: device.NewScyllaDeviceCredentialsProvider(
			s.Config.ScyllaDBAddress, s.Config.ScyllaDBPort, s.Config.KeySpace),
		tokenProvider: token.NewScyllaTokenProvider(
			s.Config.ScyllaDBAddress, s.Config.ScyllaDBPort, s.Config.KeySpace),
		devTokenProvider:device_token.NewScyllaDeviceTokenProvider(
			s.Config.ScyllaDBAddress, s.Config.ScyllaDBPort, s.Config.KeySpace),
	}
}

// GetProviders builds the providers according to the selected backend.
func (s *Service) GetProviders() * Providers {
	if s.Config.UseInMemoryProviders {
		return s.CreateInMemoryProviders()
	} else if s.Config.UseDBScyllaProviders {
		return s.CreateDBScyllaProviders()
	}
	log.Fatal().Msg("unsupported type of provider")
	return nil
}

func (s *Service) createInMemoryManagers() * TokenManagers {
	return &TokenManagers{
		tokenManager: manager.NewJWTTokenMockup(),
		deviceTokenManager: manager.NewJWTDeviceTokenMockup(),
	}
}
func (s *Service) createDBScyllaManagers(tokenProvider token.Token, password manager.Password,
	deviceProvider device.Provider, deviceTokenProvider device_token.Provider) * TokenManagers {
	return &TokenManagers{
		tokenManager: manager.NewJWTToken(tokenProvider, password),
		deviceTokenManager: manager.NewJWTDeviceToken(deviceProvider, deviceTokenProvider ),
	}
}


func (s *Service) getTokenManager(tokenProvider token.Token, password manager.Password,
	deviceProvider device.Provider, deviceTokenProvider device_token.Provider) * TokenManagers {
	if s.Config.UseInMemoryProviders {
		return s.createInMemoryManagers()
	} else if s.Config.UseDBScyllaProviders {
		return s.createDBScyllaManagers(tokenProvider, password, deviceProvider, deviceTokenProvider)
	}
	log.Fatal().Msg("unsupported type of provider")
	return nil
}


//Run launch the Authx service.
func (s *Service) Run() {
	s.Config.Print()
	p := s.GetProviders()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		log.Fatal().Errs("failed to listen: %v", []error{err})
		return
	}

	passwordMgr := manager.NewBCryptPassword()

	// Create the token manager (memory/scylla)
	t := s.getTokenManager(p.tokenProvider, passwordMgr, p.devProvider, p.devTokenProvider)
	tokenMgr := t.tokenManager
	deviceMgr := t.deviceTokenManager

	authxMgr := manager.NewAuthx(passwordMgr, tokenMgr,deviceMgr, p.credProvider, p.roleProvider, p.devProvider,
		s.Secret, s.ExpirationTime, s.DeviceExpirationTime, p.devTokenProvider)

	h := handler.NewAuthx(authxMgr)


	grpcServer := grpc.NewServer()

	pbAuthx.RegisterAuthxServer(grpcServer, h)

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)
	log.Info().Int("Port", s.Port).Msg("Launching gRPC server")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Errs("failed to serve: %v", []error{err})
	}
}
