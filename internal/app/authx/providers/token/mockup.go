/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package token

import (
	"fmt"
	"github.com/nalej/authx/internal/app/authx/entities"
	"github.com/nalej/derrors"
	"sync"
	"time"
)

// TokenMockup is an in-memory mockup.
type TokenMockup struct {
	sync.Mutex
	data map[string]entities.TokenData
}

// NewTokenMockup create a new instance of TokenMockup.
func NewTokenMockup() Token {
	return &TokenMockup{data: make(map[string]entities.TokenData, 0)}
}

func (p *TokenMockup) unsafeExists(username string, tokenID string) bool {
	_, ok := p.data[p.generateID(tokenID, username)]
	return ok
}

func (p *TokenMockup) unsafeGet(username string, tokenID string) (*entities.TokenData, derrors.Error) {

	data, ok := p.data[p.generateID(tokenID, username)]
	if !ok {
		return nil, derrors.NewNotFoundError("token not found").WithParams(username, tokenID)
	}
	return &data, nil
}

// Delete an existing token.
func (p *TokenMockup) Delete(username string, tokenID string) derrors.Error {

	p.Lock()
	defer p.Unlock()

	id := p.generateID(tokenID, username)
	_, err := p.unsafeGet(username, tokenID)
	if err != nil {
		return derrors.NewNotFoundError("username not found").WithParams(username)
	}

	delete(p.data, id)
	return nil
}

// Add a token.
func (p *TokenMockup) Add(token *entities.TokenData) derrors.Error {
	p.Lock()
	defer p.Unlock()
	if p.unsafeExists(token.Username, token.TokenID) {
		return derrors.NewAlreadyExistsError("token").WithParams(token.Username, token.TokenID)
	}
	p.data[p.generateID(token.TokenID, token.Username)] = *token
	return nil
}

// Get an existing token.
func (p *TokenMockup) Get(username string, tokenID string) (*entities.TokenData, derrors.Error) {
	p.Lock()
	defer p.Unlock()
	data, ok := p.data[p.generateID(tokenID, username)]
	if !ok {
		return nil, derrors.NewNotFoundError("token not found").WithParams(username, tokenID)
	}
	return &data, nil
}

// Exist checks if the token was added.
func (p *TokenMockup) Exist(username string, tokenID string) (*bool, derrors.Error) {
	p.Lock()
	defer p.Unlock()
	_, ok := p.data[p.generateID(tokenID, username)]
	return &ok, nil
}

func (p *TokenMockup) Update(token *entities.TokenData) derrors.Error {

	p.Lock()
	defer p.Unlock()

	if !p.unsafeExists(token.Username, token.TokenID) {
		return derrors.NewNotFoundError("token").WithParams(token.Username, token.TokenID)
	}
	p.data[p.generateID(token.TokenID, token.Username)] = *token
	return nil
}

// Truncate cleans all data.
func (p *TokenMockup) Truncate() derrors.Error {
	p.Lock()
	defer p.Unlock()
	p.data = make(map[string]entities.TokenData, 0)
	return nil
}

func (p *TokenMockup) generateID(tokenID string, username string) string {
	return fmt.Sprintf("%s:%s", username, tokenID)
}

func (p *TokenMockup) DeleteExpiredTokens() derrors.Error {
	p.Lock()
	defer p.Unlock()

	idBorrow := make([]string, 0)

	for _, token := range p.data {
		if token.ExpirationDate < time.Now().Unix() {
			id := p.generateID(token.TokenID, token.Username)
			idBorrow = append(idBorrow, id)
		}

	}
	for _, id := range idBorrow {
		delete(p.data, id)
	}
	return nil
}
