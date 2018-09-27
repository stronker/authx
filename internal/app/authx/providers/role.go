/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package providers

import "github.com/nalej/derrors"

type RoleData struct {
	Username       string
	Password       [] byte
	RoleID         string
	OrganizationID string
}

func NewBasicRoleData(username string, password [] byte, roleID string, organizationID string) *BasicCredentialsData {
	return &BasicCredentialsData{
		Username:       username,
		Password:       password,
		RoleID:         roleID,
		OrganizationID: organizationID,
	}
}

type Role interface {
	Delete(roleID string) derrors.Error
	Add(role *RoleData) derrors.Error
	Get(roleID string) (*RoleData, derrors.Error)
}

type RoleMockup struct {
	data map[string]RoleData
}

func (p *RoleMockup) Delete(roleID string) derrors.Error {
	_, ok := p.data[roleID]
	if !ok {
		return derrors.NewOperationError("Not found username")
	}
	delete(p.data, roleID)
	return nil
}

func (p *RoleMockup) Add(role *RoleData) derrors.Error {
	p.data[role.RoleID] = *role
	return nil
}

func (p *RoleMockup) Get(roleID string) (*RoleData, derrors.Error) {
	data := p.data[roleID]
	return &data, nil
}





