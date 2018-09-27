/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package providers

import (
	"fmt"
	"github.com/nalej/derrors"
)

type TokenData struct {
	Username       string
	TokenID        string
	Token          []byte
	ExpirationDate int64
}

type Token interface {
	Delete(username string) derrors.Error
	Add(token *TokenData) derrors.Error
	Get(username string) (*TokenData, derrors.Error)
}

type TokenMockup struct {
	data map[string]TokenData
}

func (p *TokenMockup) Delete(username string, tokenID string) derrors.Error {
	_, ok := p.data[p.GenerateID(tokenID, username)]
	if !ok {
		return derrors.NewOperationError("username not found")
	}
	delete(p.data, username)
	return nil
}

func (p *TokenMockup) Add(token *TokenData) derrors.Error {
	p.data[p.GenerateID(token.TokenID, token.Username)] = *token
	return nil
}

func (p *TokenMockup) Get(username string, tokenID string) (*TokenData, derrors.Error) {
	data := p.data[p.GenerateID(tokenID, username)]
	return &data, nil
}

func (p *TokenMockup) GenerateID(tokenID string, username string) string {
	return fmt.Sprintf("%s:%s", username, tokenID)
}
