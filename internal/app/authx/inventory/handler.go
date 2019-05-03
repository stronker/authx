/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package inventory

import (
	"context"
	"github.com/nalej/authx/internal/app/entities"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
)

// Handler for inventory related operations.
type Handler struct{
	manager Manager
}

// NewHandler creates a new handler with a given manager
func NewHandler(manager Manager) * Handler{
	return &Handler{
		manager: manager,
	}
}

// CreateEICJoinToken creates an EICJoinToken for new controllers to join the system.
func (h * Handler) CreateEICJoinToken(ctx context.Context, organizationID *grpc_organization_go.OrganizationId) (*grpc_authx_go.EICJoinToken, error) {
	vErr := entities.ValidOrganizationID(organizationID)
	if vErr != nil{
		return nil, conversions.ToGRPCError(vErr)
	}
	token, err := h.manager.CreateEICJoinToken(organizationID)
	if err != nil{
		log.Error().Str("trace", err.DebugReport()).Msg("cannot generate join token")
		return nil, conversions.ToGRPCError(derrors.NewInternalError("cannot generate join token"))
	}
	log.Debug().Msg("EC join token has been generated")
	return token, nil
}

// ValidEICJoinToken checks if the Token provide is still valid to join a EIC.
func (h * Handler) ValidEICJoinToken(ctx context.Context, token *grpc_authx_go.EICJoinRequest) (*grpc_common_go.Success, error) {
	if token.Token == "" {
		return nil, conversions.ToGRPCError(derrors.NewInvalidArgumentError("token is mandatory"))
	}
	valid, err := h.manager.ValidEICJoinToken(token)
	if err != nil{
		return nil, conversions.ToGRPCError(err)
	}
	if valid {
		return &grpc_common_go.Success{}, nil
	}
	return nil, conversions.ToGRPCError(derrors.NewUnauthenticatedError("invalid token"))
}

