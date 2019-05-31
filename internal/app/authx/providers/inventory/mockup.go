/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package inventory

import (
	"github.com/nalej/authx/internal/app/authx/entities"
	"github.com/nalej/derrors"
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
