/*
* Copyright (C) 2019 Nalej - All Rights Reserved
*/

package certificates

import (
	"context"
	"github.com/nalej/authx/internal/app/authx/entities"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	manager Manager
}

func NewHandler(manager Manager) *Handler {
	return &Handler{
		manager,
	}
}

func (h*Handler) CreateControllerCert(ctx context.Context, request *grpc_authx_go.EdgeControllerCertRequest) (*grpc_authx_go.PEMCertificate, error) {
	vErr := entities.ValidEdgeControllerCertRequest(request)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	pem, err := h.manager.CreateControllerCert(request)
	if err != nil{
		log.Warn().Str("trace", err.DebugReport()).Msg("cannot create ")
		return nil, conversions.ToGRPCError(err)
	}
	return pem, nil
}

