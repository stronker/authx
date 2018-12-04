/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package role

import (
	"github.com/nalej/authx/internal/app/authx/entities"
	"github.com/nalej/derrors"
)

// Role is the interface to store role entities in the system.
type Role interface {
	// Delete an existing role.
	Delete(organizationID string, roleID string) derrors.Error
	// Add a new role.
	Add(role *entities.RoleData) derrors.Error
	// Get recovers an existing role.
	Get(organizationID string, roleID string) (*entities.RoleData, derrors.Error)
	// Edit updates an existing role.
	Edit(organizationID string, roleID string, edit *entities.EditRoleData) derrors.Error
	// Exist checks if a role exists.
	Exist(organizationID string, roleID string) (*bool, derrors.Error)
	// List the roles associated with an organization.
	List(organizationID string) ([]entities.RoleData, derrors.Error)
	// Truncate clears the provider.
	Truncate() derrors.Error
}

