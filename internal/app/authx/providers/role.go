/*
 * Copyright 2018 Nalej
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





