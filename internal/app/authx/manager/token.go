/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package manager

import (
	"github.com/nalej/authx/internal/app/authx/providers"
	"github.com/nalej/derrors"
)

type GeneratedToken struct {
	Token        string
	RefreshToken string
}

func NewGeneratedToken(token string, refreshToken string) *GeneratedToken {
	return &GeneratedToken{Token: token, RefreshToken: refreshToken}
}

type Token interface {
	Generate(username string, data map[string]string) (GeneratedToken, derrors.Error)
	Refresh(username string, refreshToken string) (GeneratedToken, derrors.Error)
}

type JWTToken struct {
	provider providers.Token
}

func (m *JWTToken) Generate(username string, data map[string]string) (GeneratedToken, derrors.Error) {
	panic("implement me")
}

func (m *JWTToken) Refresh(username string, refreshToken string) (GeneratedToken, derrors.Error) {
	panic("implement me")
}
