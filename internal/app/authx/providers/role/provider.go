/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package role

import (
	"github.com/nalej/derrors"
	"github.com/stronker/authx/internal/app/authx/entities"
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
