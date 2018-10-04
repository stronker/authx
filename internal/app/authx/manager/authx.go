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

const DefaultExpirationDuration = "10h"
const DefaultSecret = "MyLittleSecret"

type Authx struct {
	Password            Password
	Token               Token
	CredentialsProvider providers.BasicCredentials
	RoleProvider        providers.Role

	secret             string
	expirationDuration time.Duration
}

func NewAuthx(password Password, tokenManager Token, credentialsProvider providers.BasicCredentials,
	roleProvide providers.Role, secret string, expirationDuration time.Duration) *Authx {

	return &Authx{Password: password, Token: tokenManager,
		CredentialsProvider: credentialsProvider, RoleProvider: roleProvide,
		secret: secret, expirationDuration: expirationDuration}

}

func NewAuthxMockup() *Authx {
	d, _ := time.ParseDuration(DefaultExpirationDuration)
	return NewAuthx(NewBCryptPassword(), NewJWTTokenMockup(),
		providers.NewBasicCredentialMockup(), providers.NewRoleMockup(),
		DefaultSecret, d)
}

func (m *Authx) DeleteCredentials(username string) derrors.Error {
	return m.CredentialsProvider.Delete(username)
}

func (m *Authx) AddBasicCredentials(username string, organizationID string, roleID string, password string) derrors.Error {
	role, err := m.RoleProvider.Get(organizationID, roleID)
	if err != nil {
		return err
	}
	if role == nil{
		return derrors.NewOperationError("role not found")
	}

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
	role, err := m.RoleProvider.Get(credentials.OrganizationID, credentials.RoleID)
	if err != nil {
		return nil, err
	}
	personalClaim := token.NewPersonalClaim(username, role.Name, role.Primitives)
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
	role, err := m.RoleProvider.Get(credentials.OrganizationID, credentials.RoleID)
	if err != nil {
		return nil, err
	}
	personalClaim := token.NewPersonalClaim(username, role.Name, role.Primitives)

	gToken, err := m.Token.Refresh(personalClaim, tokenID, refreshToken, m.expirationDuration, m.secret)
	if err != nil {
		return nil, err
	}
	response := &pbAuthx.LoginResponse{Token: gToken.Token, RefreshToken: gToken.RefreshToken}
	return response, nil
}

func (m *Authx) AddRole(role *pbAuthx.Role) derrors.Error {
	entity := providers.NewRoleData(role.OrganizationId, role.RoleId, role.Name, PrimitivesToString(role.Primitives))
	return m.RoleProvider.Add(entity)
}

func (m *Authx) EditUserRole(username string, roleID string) derrors.Error {
	credentials, err := m.CredentialsProvider.Get(username)
	if err != nil {
		return err
	}
	role, err := m.RoleProvider.Get(credentials.OrganizationID, roleID)
	if err != nil {
		return err
	}

	if role == nil{
		return derrors.NewOperationError("role not found")
	}

	edit := providers.NewEditBasicCredentialsData().WithRoleID(roleID)
	return m.CredentialsProvider.Edit(username, edit)
}

func (m *Authx) Clean() derrors.Error {
	err := m.Token.Clean()
	if err != nil {
		return err
	}
	err = m.CredentialsProvider.Truncate()
	if err != nil {
		return err
	}
	err = m.RoleProvider.Truncate()
	if err != nil {
		return err
	}
	return nil
}

func PrimitivesToString(primitives [] pbAuthx.AccessPrimitive) [] string {
	strPrimitives := make([] string, 0, len(primitives))
	for _, p := range primitives {
		strPrimitives = append(strPrimitives, p.String())
	}
	return strPrimitives
}
