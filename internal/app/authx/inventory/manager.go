/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package inventory

import (
	"github.com/nalej/authx/internal/app/authx/config"
	"github.com/nalej/authx/internal/app/authx/entities"
	"github.com/nalej/authx/internal/app/authx/providers/inventory"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-organization-go"
	"time"
)

type Manager struct {
	provider inventory.Provider
	serverCert string
	joinTokenDuration time.Duration
}

func NewManager(provider inventory.Provider, cfg config.Config) Manager{
	return Manager{
		provider:provider,
		serverCert: cfg.ManagementClusterCert,
		joinTokenDuration: cfg.EdgeControllerExpTime,
	}
}

func (m * Manager) CreateEICJoinToken(organizationID *grpc_organization_go.OrganizationId) (*grpc_authx_go.EICJoinToken, derrors.Error) {
	token := entities.NewEICJoinToken(organizationID.OrganizationId, m.joinTokenDuration)
	// Store the token
	err := m.provider.AddECJoinToken(token)
	if err != nil{
		return nil, err
	}

	return &grpc_authx_go.EICJoinToken{
		OrganizationId: organizationID.OrganizationId,
		Token: token.TokenID,
		Cacert: m.serverCert,
	}, nil
}

func (m * Manager) ValidEICJoinToken(token *grpc_authx_go.EICJoinRequest) (bool, derrors.Error) {
	stored, err := m.provider.GetECJoinToken(token.OrganizationId, token.Token)
	if err != nil{
		return false, err
	}
	if stored.ExpiresOn > time.Now().Unix(){
		return true, nil
	}
	return false, nil
}