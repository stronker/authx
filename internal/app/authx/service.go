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

//Run launch the Authx service.
func (s *Service) Run() {
	s.Config.Print()
	p := s.GetProviders()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		log.Fatal().Errs("failed to listen: %v", []error{err})
		return
	}

	//roleProvider := role.NewRoleMockup()
	//credProvider := credentials.NewBasicCredentialMockup()
	passwordMgr := manager.NewBCryptPassword()
	tokenMgr := manager.NewJWTTokenMockup()
	authxMgr := manager.NewAuthx(passwordMgr, tokenMgr, p.credProvider, p.roleProvider, p.devProvider, s.Secret, s.ExpirationTime)

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
