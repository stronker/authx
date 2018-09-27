/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package manager

import (
"github.com/nalej/authx/internal/app/authx/providers"
"github.com/nalej/derrors"
pbAuthx "github.com/nalej/grpc-authx-go"
)

type Authx struct {
	Password            Password
	CredentialsProvider providers.BasicCredentials
	RoleProvider        providers.Role
	TokenProvider		providers.Token
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

func (m *Authx) LoginWithBasicCredentials(username string, password string) derrors.Error {
	credentials, err := m.CredentialsProvider.Get(username)
	if err != nil {
		return err
	}
	err=m.Password.CompareHashAndPassword(credentials.Password, password)
	if err != nil{
		return err
	}
	//role,err:=m.RoleProvider.Get(credentials.RoleID)

	panic("implement me")
}

func (m *Authx) AddRole(role *pbAuthx.Role) derrors.Error {
	entity := providers.NewRoleData(role.OrganizationId, role.RoleId, role.Name, role.Primitives)
	return m.RoleProvider.Add(entity)
}

func (m *Authx) EditUserRole(username string, roleID string) derrors.Error {
	edit := providers.NewEditBasicCredentialsData().WithRoleID(roleID)
	return m.CredentialsProvider.Edit(username, edit)
}
