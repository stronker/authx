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

package certificates

import (
	"context"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
	"github.com/stronker/authx/internal/app/authx/entities"
)

type Handler struct {
	manager Manager
}

func NewHandler(manager Manager) *Handler {
	return &Handler{
		manager,
	}
}

// CreateControllerCert creates a certificate for an edge controller.
func (h *Handler) CreateControllerCert(ctx context.Context, request *grpc_authx_go.EdgeControllerCertRequest) (*grpc_authx_go.PEMCertificate, error) {
	vErr := entities.ValidEdgeControllerCertRequest(request)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	pem, err := h.manager.CreateControllerCert(request)
	if err != nil {
		log.Warn().Str("trace", err.DebugReport()).Msg("cannot create ")
		return nil, conversions.ToGRPCError(err)
	}
	return pem, nil
}
