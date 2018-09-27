/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package authx

import (
	"fmt"
	"github.com/nalej/authx/internal/app/authx/handler"
	pbAuthx "github.com/nalej/grpc-authx-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		log.Fatal().Errs("failed to listen: %v", []error{err})
		return
	}

	h := handler.NewAuthxServer()

	// Create the TLS credentials
	creds, err := credentials.NewServerTLSFromFile(s.certPath, s.keyPath)
	if err != nil {
		log.Fatal().Errs("could not load TLS keys: %s", []error{err})
		return
	}
	grpcServer := grpc.NewServer(grpc.Creds(creds))

	pbAuthx.RegisterAuthxServer(grpcServer, h)

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)
	log.Info().Int("port", s.Port).Msg("Launching gRPC server")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Errs("failed to serve: %v", []error{err})
	}
}
