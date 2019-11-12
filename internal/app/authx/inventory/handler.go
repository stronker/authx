/*
 * Copyright 2019 Nalej
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
 *
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
type Handler struct {
	manager Manager
}

// NewHandler creates a new handler with a given manager
func NewHandler(manager Manager) *Handler {
	return &Handler{
		manager: manager,
	}
}

// CreateEICJoinToken creates an EICJoinToken for new controllers to join the system.
func (h *Handler) CreateEICJoinToken(ctx context.Context, organizationID *grpc_organization_go.OrganizationId) (*grpc_authx_go.EICJoinToken, error) {
	vErr := entities.ValidOrganizationID(organizationID)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	token, err := h.manager.CreateEICJoinToken(organizationID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot generate join token")
		return nil, conversions.ToGRPCError(derrors.NewInternalError("cannot generate join token"))
	}
	log.Debug().Msg("EC join token has been generated")
	return token, nil
}

// ValidEICJoinToken checks if the Token provide is still valid to join a EIC.
func (h *Handler) ValidEICJoinToken(ctx context.Context, token *grpc_authx_go.EICJoinRequest) (*grpc_common_go.Success, error) {
	if token.Token == "" {
		return nil, conversions.ToGRPCError(derrors.NewInvalidArgumentError("token is mandatory"))
	}
	valid, err := h.manager.ValidEICJoinToken(token)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if valid {
		return &grpc_common_go.Success{}, nil
	}
	return nil, conversions.ToGRPCError(derrors.NewUnauthenticatedError("invalid token"))
}
