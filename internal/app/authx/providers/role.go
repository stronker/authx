/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package providers

import (
	"github.com/nalej/derrors"
)

// RoleData is the structure that is stored in the provider.
type RoleData struct {
	OrganizationID string
	RoleID         string
	Name           string
	Primitives     []string
}

// NewRoleData create a new instance of the structure.
func NewRoleData(organizationID string, roleID string, name string, primitives []string) *RoleData {
	return &RoleData{
		OrganizationID: organizationID,
		RoleID:         roleID,
		Name:           name,
		Primitives:     primitives,
	}
}

// EditRoleData is the structure that is used to edit the data in the provider.
type EditRoleData struct {
	Name       *string
	Primitives *[]string
}

//WithName update the name of the role.
func (d *EditRoleData) WithName(name string) *EditRoleData {
	d.Name = &name
	return d
}

//WithPrimitives update the primitives.
func (d *EditRoleData) WithPrimitives(primitives []string) *EditRoleData {
	d.Primitives = &primitives
	return d
}

//NewEditRoleData create a new instance of the structure.
func NewEditRoleData() *EditRoleData {
	return &EditRoleData{}
}

// Role is the interface to store role entities in the system.
type Role interface {
	// Delete an existing role.
	Delete(organizationID string, roleID string) derrors.Error
	// Add a new role.
	Add(role *RoleData) derrors.Error
	// Get recovers an existing role.
	Get(organizationID string, roleID string) (*RoleData, derrors.Error)
	// Edit updates an existing role.
	Edit(organizationID string, roleID string, edit *EditRoleData) derrors.Error
	// Exist checks if a role exists.
	Exist(username string, tokenID string) (*bool, derrors.Error)
	// Truncate clears the provider.
	Truncate() derrors.Error
}

// RoleMockup is a in-memory provider.
type RoleMockup struct {
	data map[string]RoleData
}

// NewRoleMockup create a new instance of the RoleMockup structure.
func NewRoleMockup() Role {
	return &RoleMockup{data: map[string]RoleData{}}
}

// Delete an existing role.
func (p *RoleMockup) Delete(organizationID string, roleID string) derrors.Error {
	data, ok := p.data[roleID]
	if !ok || data.OrganizationID != organizationID {
		return derrors.NewNotFoundError("role not found").WithParams(roleID)
	}
	delete(p.data, roleID)
	return nil
}

// Add a new role.
func (p *RoleMockup) Add(role *RoleData) derrors.Error {
	p.data[role.RoleID] = *role
	return nil
}

// Get recovers an existing role.
func (p *RoleMockup) Get(organizationID string, roleID string) (*RoleData, derrors.Error) {
	data, ok := p.data[roleID]
	if !ok || data.OrganizationID != organizationID {
		return nil, derrors.NewNotFoundError("role not found").WithParams(organizationID, roleID)
	}

	return &data, nil
}

// Edit updates an existing role.
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

// Exist checks if a role exists.
func (p *RoleMockup) Exist(organizationID string, roleID string) (*bool, derrors.Error) {
	result := true
	data, ok := p.data[roleID]
	if !ok || data.OrganizationID != organizationID {
		result = false
		return &result, nil
	}
	return &result, nil
}

// Truncate clears the provider.
func (p *RoleMockup) Truncate() derrors.Error {
	p.data = map[string]RoleData{}
	return nil
}
