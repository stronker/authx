/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package providers

import (
	"github.com/nalej/derrors"
	pbAuthx "github.com/nalej/grpc-authx-go"
)

type RoleData struct {
	OrganizationId string
	RoleId         string
	Name           string
	Primitives     []pbAuthx.AccessPrimitive
}

func NewRoleData(organizationID string, roleID string, name string, primitives []pbAuthx.AccessPrimitive) *RoleData {
	return &RoleData{
		OrganizationId: organizationID,
		RoleId:         roleID,
		Name:           name,
		Primitives:     primitives,
	}
}

type EditRoleData struct {
	Name       *string
	Primitives *[]pbAuthx.AccessPrimitive
}

func (d *EditRoleData) WithName(name string) *EditRoleData {
	d.Name = &name
	return d
}

func (d *EditRoleData) WithPrimitives(primitives []pbAuthx.AccessPrimitive) *EditRoleData {
	d.Primitives = &primitives
	return d
}

func NewEditRoleData() *EditRoleData {
	return &EditRoleData{}
}

type Role interface {
	Delete(roleID string) derrors.Error
	Add(role *RoleData) derrors.Error
	Get(roleID string) (*RoleData, derrors.Error)
	Edit(roleID string, edit EditRoleData) derrors.Error
}

type RoleMockup struct {
	data map[string]RoleData
}

func (p *RoleMockup) Delete(roleID string) derrors.Error {
	_, ok := p.data[roleID]
	if !ok {
		return derrors.NewOperationError("username not found")
	}
	delete(p.data, roleID)
	return nil
}

func (p *RoleMockup) Add(role *RoleData) derrors.Error {
	p.data[role.RoleId] = *role
	return nil
}

func (p *RoleMockup) Get(roleID string) (*RoleData, derrors.Error) {
	data := p.data[roleID]
	return &data, nil
}

func (p *RoleMockup) Edit(roleID string, edit EditRoleData) derrors.Error {
	data, ok := p.data[roleID]
	if !ok {
		return derrors.NewOperationError("username not found")
	}
	if edit.Name != nil {
		data.Name = *edit.Name
	}
	if edit.Primitives != nil {
		data.Primitives = *edit.Primitives
	}
	return nil
}
