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

package inventory

import (
	"github.com/nalej/derrors"
	"github.com/stronker/authx/internal/app/authx/entities"
	"sync"
	"time"
)

// MockupInventoryTTL to match the scylladb one. This value will automatically expire entries.
const MockupInventoryTTL = time.Hour * 3

type MockupInventoryProvider struct {
	// Mutex for managing mockup access.
	sync.Mutex
	// eicJoinToken with the join tokens per tokenID
	eicJoinToken map[string]entities.EICJoinToken
}

func NewMockupInventoryProvider() *MockupInventoryProvider {
	return &MockupInventoryProvider{
		eicJoinToken: make(map[string]entities.EICJoinToken, 0),
	}
}

func (m *MockupInventoryProvider) unsafeExistECJoinToken(organizationID string, token string) bool {
	retrieved, existToken := m.eicJoinToken[token]
	return existToken && (retrieved.OrganizationID == organizationID)
}

func (m *MockupInventoryProvider) AddECJoinToken(token *entities.EICJoinToken) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if m.unsafeExistECJoinToken(token.OrganizationID, token.TokenID) {
		derrors.NewAlreadyExistsError("token already exists")
	}
	m.eicJoinToken[token.TokenID] = *token
	return nil
}

func (m *MockupInventoryProvider) GetECJoinToken(organizationID string, token string) (*entities.EICJoinToken, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExistECJoinToken(organizationID, token) {
		return nil, derrors.NewNotFoundError("invalid token")
	}
	result, _ := m.eicJoinToken[token]
	return &result, nil
}

func (m *MockupInventoryProvider) Clear() derrors.Error {
	m.Lock()
	m.eicJoinToken = make(map[string]entities.EICJoinToken, 0)
	m.Unlock()
	return nil
}
