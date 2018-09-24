/*
 * Copyright 2018 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
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
