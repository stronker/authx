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

package device_token

import (
	"fmt"
	"github.com/nalej/authx/internal/app/authx/entities"
	"github.com/nalej/derrors"
	"sync"
	"time"
)

// DeviceTokenMockup is an in-memory mockup.
type DeviceTokenMockup struct {
	sync.Mutex
	data               map[string]entities.DeviceTokenData
	dataByRefreshToken map[string]entities.DeviceTokenData
}

// NewTokenMockup create a new instance of TokenMockup.
func NewDeviceTokenMockup() Provider {
	return &DeviceTokenMockup{
		data:               make(map[string]entities.DeviceTokenData, 0),
		dataByRefreshToken: make(map[string]entities.DeviceTokenData, 0),
	}
}
func (m *DeviceTokenMockup) unsafeExists(deviceID string, tokenID string) bool {
	_, ok := m.data[m.generateID(deviceID, tokenID)]
	return ok
}

func (m *DeviceTokenMockup) unsafeGet(deviceID string, tokenID string) (*entities.DeviceTokenData, derrors.Error) {

	data, ok := m.data[m.generateID(deviceID, tokenID)]
	if !ok {
		return nil, derrors.NewNotFoundError("device token not found").WithParams(deviceID, tokenID)
	}
	return &data, nil
}

func (m *DeviceTokenMockup) generateID(deviceID string, tokenID string) string {
	return fmt.Sprintf("%s:%s", deviceID, tokenID)
}

// Delete an existing token.
func (m *DeviceTokenMockup) Delete(deviceID string, tokenID string) derrors.Error {
	m.Lock()
	defer m.Unlock()

	id := m.generateID(deviceID, tokenID)
	token, err := m.unsafeGet(deviceID, tokenID)
	if err != nil {
		return derrors.NewNotFoundError("device not found").WithParams(deviceID)
	}

	delete(m.data, id)
	delete(m.dataByRefreshToken, token.RefreshToken)
	return nil
}

// Add a token.
func (m *DeviceTokenMockup) Add(token *entities.DeviceTokenData) derrors.Error {
	m.Lock()
	defer m.Unlock()

	if m.unsafeExists(token.DeviceId, token.TokenID) {
		return derrors.NewAlreadyExistsError("device token").WithParams(token.DeviceId, token.TokenID)
	}
	m.data[m.generateID(token.DeviceId, token.TokenID)] = *token
	m.dataByRefreshToken[token.RefreshToken] = *token
	return nil
}

// Get an existing token.
func (m *DeviceTokenMockup) Get(deviceID string, tokenID string) (*entities.DeviceTokenData, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	data, ok := m.data[m.generateID(deviceID, tokenID)]
	if !ok {
		return nil, derrors.NewNotFoundError("device token not found").WithParams(deviceID, tokenID)
	}
	return &data, nil
}

// Exist checks if the token was added.
func (m *DeviceTokenMockup) Exist(deviceID string, tokenID string) (*bool, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	_, ok := m.data[m.generateID(deviceID, tokenID)]
	return &ok, nil
}

// Update an existing token
func (m *DeviceTokenMockup) Update(token *entities.DeviceTokenData) derrors.Error {
	m.Lock()
	defer m.Unlock()

	oldToken, err := m.unsafeGet(token.DeviceId, token.TokenID)
	if err != nil {
		return derrors.NewNotFoundError("device token").WithParams(token.DeviceId, token.TokenID)
	}
	delete(m.dataByRefreshToken, oldToken.RefreshToken)
	m.data[m.generateID(token.DeviceId, token.TokenID)] = *token
	m.dataByRefreshToken[token.RefreshToken] = *token
	return nil
}

// Truncate cleans all data.
func (m *DeviceTokenMockup) Truncate() derrors.Error {
	m.Lock()
	defer m.Unlock()

	m.data = make(map[string]entities.DeviceTokenData, 0)
	m.dataByRefreshToken = make(map[string]entities.DeviceTokenData, 0)
	return nil
}

func (m *DeviceTokenMockup) DeleteExpiredTokens() derrors.Error {
	m.Lock()
	defer m.Unlock()

	idBorrow := make([]string, 0)

	for _, token := range m.data {
		if token.ExpirationDate < time.Now().Unix() {
			id := m.generateID(token.DeviceId, token.TokenID)
			idBorrow = append(idBorrow, id)
			delete(m.dataByRefreshToken, token.RefreshToken)
		}

	}
	for _, id := range idBorrow {
		delete(m.data, id)
	}
	return nil
}

func (m *DeviceTokenMockup) GetByRefreshToken(refreshToken string) (*entities.DeviceTokenData, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	data, ok := m.dataByRefreshToken[refreshToken]
	if !ok {
		return nil, derrors.NewNotFoundError("device token not found").WithParams(refreshToken)
	}
	return &data, nil
}
