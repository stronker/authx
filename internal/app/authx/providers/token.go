/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package providers

import (
	"fmt"
	"github.com/nalej/derrors"
	"sync"
)

// TokenData is the information that the system stores.
type TokenData struct {
	Username       string
	TokenID        string
	RefreshToken   []byte
	ExpirationDate int64
}

// NewTokenData creates an instance of the structure
func NewTokenData(username string, tokenID string, refreshToken []byte,
	expirationDate int64) *TokenData {

	return &TokenData{Username: username,
		TokenID:        tokenID,
		RefreshToken:   refreshToken,
		ExpirationDate: expirationDate}
}
// Token is the interface to store the token information.
type Token interface {
	// Delete an existing token.
	Delete(username string, tokenID string) derrors.Error
	// Add a token.
	Add(token *TokenData) derrors.Error
	// Get an existing token.
	Get(username string, tokenID string) (*TokenData, derrors.Error)
	// Exist checks if the token was added.
	Exist(username string, tokenID string) (*bool, derrors.Error)
	// Truncate cleans all data.
	Truncate() derrors.Error
}
// TokenMockup is an in-memory mockup.
type TokenMockup struct {
	sync.Mutex
	data map[string]TokenData
}

// NewTokenMockup create a new instance of TokenMockup.
func NewTokenMockup() Token {
	return &TokenMockup{data: map[string]TokenData{}}
}

// Delete an existing token.
func (p *TokenMockup) Delete(username string, tokenID string) derrors.Error {

	id := p.generateID(tokenID, username)
	_, err := p.Get(username, tokenID)
	if err != nil {
		return derrors.NewNotFoundError("username not found").WithParams(username)
	}
	p.Lock()
	defer p.Unlock()
	delete(p.data, id)
	return nil
}

// Add a token.
func (p *TokenMockup) Add(token *TokenData) derrors.Error {
	p.Lock()
	defer p.Unlock()
	p.data[p.generateID(token.TokenID, token.Username)] = *token
	return nil
}

// Get an existing token.
func (p *TokenMockup) Get(username string, tokenID string) (*TokenData, derrors.Error) {
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

// Truncate cleans all data.
func (p *TokenMockup) Truncate() derrors.Error {
	p.Lock()
	defer p.Unlock()
	p.data = map[string]TokenData{}
	return nil
}

func (p *TokenMockup) generateID(tokenID string, username string) string {
	return fmt.Sprintf("%s:%s", username, tokenID)
}
