/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package manager

import (
	"github.com/nalej/authx/internal/app/authx/providers"
	"github.com/nalej/authx/pkg/token"
	"github.com/nalej/derrors"
	pbAuthx "github.com/nalej/grpc-authx-go"
	"time"
)

type Authx struct {
	Password            Password
	Token               Token
	CredentialsProvider providers.BasicCredentials
	RoleProvider        providers.Role

	secret             string
	expirationDuration time.Duration
}

func (m *Authx) DeleteCredentials(username string) derrors.Error {
	return m.CredentialsProvider.Delete(username)
}

func (m *Authx) AddBasicCredentials(username string, organizationID string, roleID string, password string) derrors.Error {
	hashedPassword, err := m.Password.GenerateHashedPassword(password)
	if err != nil {
		return err
	}

	entity := providers.NewBasicCredentialsData(username, hashedPassword, roleID, organizationID)
	return m.CredentialsProvider.Add(entity)
}

func (m *Authx) LoginWithBasicCredentials(username string, password string) (*pbAuthx.LoginResponse, derrors.Error) {
	credentials, err := m.CredentialsProvider.Get(username)
	if err != nil {
		return nil, err
	}
	err = m.Password.CompareHashAndPassword(credentials.Password, password)
	if err != nil {
		return nil, err
	}
	role, err := m.RoleProvider.Get(credentials.RoleID)
	if err != nil {
		return nil, err
	}
	personalClaim := NewPersonalClaim(username, role.Primitives)
	gToken, err := m.Token.Generate(personalClaim, m.expirationDuration, m.secret)
	if err != nil {
		return nil, err
	}
	response := &pbAuthx.LoginResponse{Token: gToken.Token, RefreshToken: gToken.RefreshToken}
	return response, nil
}

func (m *Authx) RefreshToken(username string, tokenID string, refreshToken string) (*pbAuthx.LoginResponse, derrors.Error) {
	credentials, err := m.CredentialsProvider.Get(username)
	if err != nil {
		return nil, err
	}
	role, err := m.RoleProvider.Get(credentials.RoleID)
	if err != nil {
		return nil, err
	}
	personalClaim := NewPersonalClaim(username, role.Primitives)

	gToken, err := m.Token.Refresh(personalClaim, tokenID, refreshToken, m.expirationDuration, m.secret)
	if err != nil {
		return nil, err
	}
	response := &pbAuthx.LoginResponse{Token: gToken.Token, RefreshToken: gToken.RefreshToken}
	return response, nil
}

func (m *Authx) AddRole(role *pbAuthx.Role) derrors.Error {
	entity := providers.NewRoleData(role.OrganizationId, role.RoleId, role.Name, role.Primitives)
	return m.RoleProvider.Add(entity)
}

func (m *Authx) EditUserRole(username string, roleID string) derrors.Error {
	edit := providers.NewEditBasicCredentialsData().WithRoleID(roleID)
	return m.CredentialsProvider.Edit(username, edit)
}

func NewPersonalClaim(username string, primitives [] pbAuthx.AccessPrimitive) *token.PersonalClaim {
	strPrimitives := make([] string, 0, len(primitives))
	for _, p := range primitives {
		strPrimitives = append(strPrimitives, p.String())
	}
	return token.NewPersonalClaim(username, strPrimitives)
}
