/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package authx

import (
	"fmt"
	"github.com/nalej/authx/internal/app/authx/handler"
	"github.com/nalej/authx/internal/app/authx/manager"
	"github.com/nalej/authx/internal/app/authx/providers"
	pbAuthx "github.com/nalej/grpc-authx-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type Service struct {
	Config
}

func NewService(config Config) *Service {
	return &Service{Config: config}
}

func (s *Service) Run() {
	s.Config.Print()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		log.Fatal().Errs("failed to listen: %v", []error{err})
		return
	}

	roleProvider := providers.NewRoleMockup()
	credProvider := providers.NewBasicCredentialMockup()
	passwordMgr := manager.NewBCryptPassword()
	tokenMgr := manager.NewJWTTokenMockup()
	authxMgr := manager.NewAuthx(passwordMgr, tokenMgr, credProvider, roleProvider, s.Secret, s.ExpirationTime)

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
