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
	RefreshToken   []byte
	ExpirationDate int64
}

func NewTokenData(username string, tokenID string, refreshToken []byte,
	expirationDate int64) *TokenData {

	return &TokenData{Username: username,
		TokenID:        tokenID,
		RefreshToken:   refreshToken,
		ExpirationDate: expirationDate}
}

type Token interface {
	Delete(username string, tokenID string) derrors.Error
	Add(token *TokenData) derrors.Error
	Get(username string, tokenID string) (*TokenData, derrors.Error)
	Truncate() derrors.Error
}

type TokenMockup struct {
	data map[string]TokenData
}

func NewTokenMockup() Token {
	return &TokenMockup{data: map[string]TokenData{}}
}

func (p *TokenMockup) Delete(username string, tokenID string) derrors.Error {
	id:=p.generateID(tokenID, username)
	_, ok := p.data[id]
	if !ok {
		return derrors.NewOperationError("username not found")
	}
	delete(p.data, id)
	return nil
}

func (p *TokenMockup) Add(token *TokenData) derrors.Error {
	p.data[p.generateID(token.TokenID, token.Username)] = *token
	return nil
}

func (p *TokenMockup) Get(username string, tokenID string) (*TokenData, derrors.Error) {
	data, ok := p.data[p.generateID(tokenID, username)]
	if !ok {
		return nil, nil
	}
	return &data, nil
}

func (p *TokenMockup) Truncate() derrors.Error {
	p.data = map[string]TokenData{}
	return nil
}

func (p *TokenMockup) generateID(tokenID string, username string) string {
	return fmt.Sprintf("%s:%s", username, tokenID)
}
