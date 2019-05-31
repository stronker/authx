/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package role

import (
	"github.com/nalej/authx/internal/app/authx/entities"
	"github.com/nalej/derrors"
	"sync"
)

// RoleMockup is a in-memory provider.
type RoleMockup struct {
	sync.Mutex
	data map[string]entities.RoleData
}

// NewRoleMockup create a new instance of the RoleMockup structure.
func NewRoleMockup() Role {
	return &RoleMockup{data: make(map[string]entities.RoleData, 0)}
}

func (p *RoleMockup) unsafeGet(organizationID string, roleID string) (*entities.RoleData, derrors.Error) {
	data, ok := p.data[roleID]
	if !ok || data.OrganizationID != organizationID {
		return nil, derrors.NewNotFoundError("role not found").WithParams(organizationID, roleID)
	}

	return &data, nil
}

// Delete an existing role.
func (p *RoleMockup) Delete(organizationID string, roleID string) derrors.Error {
	p.Lock()
	defer p.Unlock()

	data, ok := p.data[roleID]
	if !ok || data.OrganizationID != organizationID {
		return derrors.NewNotFoundError("role not found").WithParams(roleID)
	}
	delete(p.data, roleID)
	return nil
}

// Add a new role.
func (p *RoleMockup) Add(role *entities.RoleData) derrors.Error {
	p.Lock()
	defer p.Unlock()
	p.data[role.RoleID] = *role
	return nil
}

// Get recovers an existing role.
func (p *RoleMockup) Get(organizationID string, roleID string) (*entities.RoleData, derrors.Error) {
	p.Lock()
	defer p.Unlock()
	data, ok := p.data[roleID]
	if !ok || data.OrganizationID != organizationID {
		return nil, derrors.NewNotFoundError("role not found").WithParams(organizationID, roleID)
	}

	return &data, nil
}

// Edit updates an existing role.
func (p *RoleMockup) Edit(organizationID string, roleID string, edit *entities.EditRoleData) derrors.Error {

	p.Lock()
	defer p.Unlock()

	data, err := p.unsafeGet(organizationID, roleID)
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
	p.Lock()
	defer p.Unlock()
	result := true
	data, ok := p.data[roleID]
	if !ok || data.OrganizationID != organizationID {
		result = false
		return &result, nil
	}
	return &result, nil
}

// List the roles associated with an organization.
func (p *RoleMockup) List(organizationID string) ([]entities.RoleData, derrors.Error) {
	p.Lock()
	defer p.Unlock()
	result := make([]entities.RoleData, 0)
	for _, r := range p.data {
		if r.OrganizationID == organizationID {
			result = append(result, r)
		}
	}
	return result, nil
}

// Truncate clears the provider.
func (p *RoleMockup) Truncate() derrors.Error {
	p.Lock()
	defer p.Unlock()
	p.data = make(map[string]entities.RoleData, 0)
	return nil
}
