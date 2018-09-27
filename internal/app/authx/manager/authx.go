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

package manager

import (
	"github.com/nalej/authx/internal/app/authx/providers"
	"github.com/nalej/derrors"
)

type Authx struct {
	Password Password
	Provider providers.BasicCredentials
}

func (m *Authx) DeleteCredentials(username string) derrors.Error {
	return m.Provider.Delete(username)
}

func (m *Authx) AddBasicCredentials(username string, organizationID string, roleID string, password string) derrors.Error {
	hashedPassword, err := m.Password.GenerateHashedPassword(password)
	if err != nil {
		return err
	}

	entity := providers.NewBasicCredentialsData(username, hashedPassword, roleID, organizationID)
	return m.Provider.Add(entity)
}

func (m *Authx) LoginWithBasicCredentials(username string, password string) derrors.Error {
	credentials, err := m.Provider.Get(username)
	if err != nil {
		return err
	}
	return m.Password.CompareHashAndPassword(credentials.Password, password)
}
