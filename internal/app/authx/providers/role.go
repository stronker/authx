/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package providers

import (
	"github.com/nalej/derrors"
)

type RoleData struct {
	OrganizationId string
	RoleId         string
	Name           string
	Primitives     []string
}

func NewRoleData(organizationID string, roleID string, name string, primitives []string) *RoleData {
	return &RoleData{
		OrganizationId: organizationID,
		RoleId:         roleID,
		Name:           name,
		Primitives:     primitives,
	}
}

type EditRoleData struct {
	Name       *string
	Primitives *[]string
}

func (d *EditRoleData) WithName(name string) *EditRoleData {
	d.Name = &name
	return d
}

func (d *EditRoleData) WithPrimitives(primitives []string) *EditRoleData {
	d.Primitives = &primitives
	return d
}

func NewEditRoleData() *EditRoleData {
	return &EditRoleData{}
}

type Role interface {
	Delete(organizationID string, roleID string) derrors.Error
	Add(role *RoleData) derrors.Error
	Get(organizationID string, roleID string) (*RoleData, derrors.Error)
	Edit(organizationID string, roleID string, edit *EditRoleData) derrors.Error
	Exist(username string, tokenID string) (*bool, derrors.Error)
	Truncate() derrors.Error
}

type RoleMockup struct {
	data map[string]RoleData
}

func NewRoleMockup() Role {
	return &RoleMockup{data: map[string]RoleData{}}
}

func (p *RoleMockup) Delete(organizationID string, roleID string) derrors.Error {
	data, ok := p.data[roleID]
	if !ok || data.OrganizationId != organizationID {
		return derrors.NewNotFoundError("role not found").WithParams(roleID)
	}
	delete(p.data, roleID)
	return nil
}

func (p *RoleMockup) Add(role *RoleData) derrors.Error {
	p.data[role.RoleId] = *role
	return nil
}

func (p *RoleMockup) Get(organizationID string, roleID string) (*RoleData, derrors.Error) {
	data, ok := p.data[roleID]
	if !ok || data.OrganizationId != organizationID {
		return nil, derrors.NewNotFoundError("role not found").WithParams(organizationID, roleID)
	}

	return &data, nil
}

func (p *RoleMockup) Edit(organizationID string, roleID string, edit *EditRoleData) derrors.Error {
	data, err := p.Get(organizationID, roleID)
	if err != nil {
		return err
	}
	if edit.Name != nil {
		data.Name = *edit.Name
	}
	if edit.Primitives != nil {
		data.Primitives = *edit.Primitives
	}
	p.data[roleID] = *data
	return nil
}
func (p *RoleMockup) Exist(organizationID string, roleID string) (*bool, derrors.Error) {
	result := true
	data, ok := p.data[roleID]
	if !ok || data.OrganizationId != organizationID {
		result = false
		return &result, nil
	}
	return &result, nil
}

func (p *RoleMockup) Truncate() derrors.Error {
	p.data = map[string]RoleData{}
	return nil
}
