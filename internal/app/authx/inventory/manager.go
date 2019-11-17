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
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/stronker/authx/internal/app/authx/config"
	"github.com/stronker/authx/internal/app/authx/entities"
	"github.com/stronker/authx/internal/app/authx/providers/inventory"
	"time"
)

type Manager struct {
	provider          inventory.Provider
	serverCert        string
	joinTokenDuration time.Duration
}

func NewManager(provider inventory.Provider, cfg config.Config) Manager {
	return Manager{
		provider:          provider,
		serverCert:        cfg.ManagementClusterCert,
		joinTokenDuration: cfg.EdgeControllerExpTime,
	}
}

func (m *Manager) CreateEICJoinToken(organizationID *grpc_organization_go.OrganizationId) (*grpc_authx_go.EICJoinToken, derrors.Error) {
	token := entities.NewEICJoinToken(organizationID.OrganizationId, m.joinTokenDuration)
	// Store the token
	err := m.provider.AddECJoinToken(token)
	if err != nil {
		return nil, err
	}
	
	return &grpc_authx_go.EICJoinToken{
		OrganizationId: organizationID.OrganizationId,
		Token:          token.TokenID,
		Cacert:         m.serverCert,
	}, nil
}

func (m *Manager) ValidEICJoinToken(token *grpc_authx_go.EICJoinRequest) (bool, derrors.Error) {
	stored, err := m.provider.GetECJoinToken(token.OrganizationId, token.Token)
	if err != nil {
		return false, err
	}
	if stored.ExpiresOn > time.Now().Unix() {
		return true, nil
	}
	return false, nil
}
