/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package manager

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/nalej/authx/internal/app/authx/providers"
	"github.com/nalej/authx/pkg/token"
	"github.com/nalej/derrors"
	"github.com/rs/zerolog/log"
	"time"
)

const Issuer string = "authx"

type GeneratedToken struct {
	Token        string
	RefreshToken string
}

func NewGeneratedToken(token string, refreshToken string) *GeneratedToken {
	return &GeneratedToken{Token: token, RefreshToken: refreshToken}
}

type Token interface {
	Generate(personalClaim *token.PersonalClaim, expirationPeriod time.Duration, secret string) (*GeneratedToken, derrors.Error)
	Refresh(personalClaim *token.PersonalClaim, tokenID string, refreshToken string, expirationPeriod time.Duration, secret string) (*GeneratedToken, derrors.Error)
}

type JWTToken struct {
	Provider providers.Token
	Password Password
}

func (m *JWTToken) Generate(personalClaim *token.PersonalClaim, expirationPeriod time.Duration, secret string) (*GeneratedToken, derrors.Error) {
	claim := token.NewClaim(*personalClaim, Issuer, time.Now(), expirationPeriod)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := t.SignedString(secret)
	if err != nil {
		return nil, derrors.NewGenericError("impossible generate JWT token", err)
	}

	refreshToken := token.GenerateUUID()
	hashedRefreshToken, err := m.Password.GenerateHashedPassword(refreshToken)
	if err != nil {
		return nil, derrors.NewGenericError("impossible generate RefreshToken", err)
	}
	tokenData := providers.NewTokenData(claim.UserID, claim.Id, hashedRefreshToken, claim.ExpiresAt)
	err = m.Provider.Add(tokenData)
	if err != nil {
		return nil, derrors.NewGenericError("impossible store RefreshToken", err)
	}
	gToken := NewGeneratedToken(tokenString, refreshToken)
	return gToken, nil
}

func (m *JWTToken) Refresh(personalClaim *token.PersonalClaim, tokenID string, refreshToken string,
	expirationPeriod time.Duration, secret string) (*GeneratedToken, derrors.Error) {
	username := personalClaim.UserID

	tokenData, err := m.Provider.Get(username, tokenID)
	if err != nil {
		return nil, derrors.NewGenericError("impossible recover RefreshToken", err)
	}
	ts := time.Now().Unix()
	if tokenData == nil || ts > tokenData.ExpirationDate {
		return nil, derrors.NewGenericError("the refresh token is expired")
	}

	err = m.Password.CompareHashAndPassword(tokenData.RefreshToken, refreshToken)
	if err != nil {
		return nil, derrors.NewGenericError("the refresh token is not valid", err)
	}

	gt, err := m.Generate(personalClaim, expirationPeriod, secret)
	if err != nil {
		return nil, derrors.NewGenericError("impossible create new token", err)
	}

	err = m.Provider.Delete(username, tokenID)
	if err != nil {
		log.Warn().Err(err).Msg("impossible delete refresh token")
	}
	return gt, nil
}
