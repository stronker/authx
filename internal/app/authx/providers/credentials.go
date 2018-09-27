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

type BasicCredentialsData struct {
	Username       string
	Password       [] byte
	RoleID         string
	OrganizationID string
}

func NewBasicCredentialsData(username string, password [] byte, roleID string, organizationID string) *BasicCredentialsData {
	return &BasicCredentialsData{
		Username:       username,
		Password:       password,
		RoleID:         roleID,
		OrganizationID: organizationID,
	}
}

type BasicCredentials interface {
	Delete(username string) derrors.Error
	Add(credentials *BasicCredentialsData) derrors.Error
	Get(username string) (*BasicCredentialsData, derrors.Error)
}

type BasicCredentialsMockup struct {
	data map[string]BasicCredentialsData
}

func NewBasicCredentialMockup() *BasicCredentialsMockup {
	return &BasicCredentialsMockup{data: map[string]BasicCredentialsData{}}
}

func (p *BasicCredentialsMockup) Delete(username string) derrors.Error {
	_, ok := p.data[username]
	if !ok {
		return derrors.NewOperationError("Not found username")
	}
	delete(p.data, username)
	return nil
}

func (p *BasicCredentialsMockup) Add(credentials *BasicCredentialsData) derrors.Error {
	p.data[credentials.Username] = *credentials
	return nil
}

func (p *BasicCredentialsMockup) Get(username string) (*BasicCredentialsData, derrors.Error) {
	data := p.data[username]
	return &data, nil
}
