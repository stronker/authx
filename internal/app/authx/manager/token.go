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

// Issuer name of the application that generates new tokens.
const Issuer string = "authx"

// GeneratedToken is the object that defines the basic structure.
type GeneratedToken struct {
	// Token is the token generated.
	Token        string
	// RefreshToken is the id required to renew an old token.
	RefreshToken string
}

// NewGeneratedToken build a new object of this structure.
func NewGeneratedToken(token string, refreshToken string) *GeneratedToken {
	return &GeneratedToken{Token: token, RefreshToken: refreshToken}
}
// Token is a interface manages the business logic of tokens.
type Token interface {
	// Generate a new token with the personal claim.
	Generate(personalClaim *token.PersonalClaim, expirationPeriod time.Duration,
		secret string) (*GeneratedToken, derrors.Error)
	// Refresh renew an old token.
	Refresh(personalClaim *token.PersonalClaim, tokenID string, refreshToken string,
		expirationPeriod time.Duration, secret string) (*GeneratedToken, derrors.Error)
	// Clean remove all the data from the providers.
	Clean() derrors.Error
}

// JWTToken is an implementation of token using JWT.
type JWTToken struct {
	TokenProvider providers.Token
	Password      Password
}

// NewJWTToken create a new instance of JWTToken
func NewJWTToken(tokenProvider providers.Token, password Password) Token {
	return &JWTToken{TokenProvider: tokenProvider, Password: password}

}

// NewJWTTokenMockup create a new mockup of JWTToken
func NewJWTTokenMockup() Token {
	return NewJWTToken(providers.NewTokenMockup(), NewBCryptPassword())
}

// Generate a new JWT token with the personal claim.
func (m *JWTToken) Generate(personalClaim *token.PersonalClaim, expirationPeriod time.Duration,
	secret string) (*GeneratedToken, derrors.Error) {

	claim := token.NewClaim(*personalClaim, Issuer, time.Now(), expirationPeriod)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := t.SignedString([]byte(secret))
	if err != nil {
		return nil, derrors.NewGenericError("impossible generate JWT token", err)
	}

	refreshToken := token.GenerateUUID()
	hashedRefreshToken, err := m.Password.GenerateHashedPassword(refreshToken)
	if err != nil {
		return nil, derrors.NewGenericError("impossible generate RefreshToken", err)
	}
	tokenData := providers.NewTokenData(claim.UserID, claim.Id, hashedRefreshToken, claim.ExpiresAt)
	err = m.TokenProvider.Add(tokenData)
	if err != nil {
		return nil, derrors.NewGenericError("impossible store RefreshToken", err)
	}
	gToken := NewGeneratedToken(tokenString, refreshToken)
	return gToken, nil
}

// Refresh renew an old token.
func (m *JWTToken) Refresh(personalClaim *token.PersonalClaim, tokenID string, refreshToken string,
	expirationPeriod time.Duration, secret string) (*GeneratedToken, derrors.Error) {

	username := personalClaim.UserID

	tokenData, err := m.TokenProvider.Get(username, tokenID)
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

	err = m.TokenProvider.Delete(username, tokenID)
	if err != nil {
		log.Warn().Err(err).Msg("impossible delete refresh token")
	}
	return gt, nil
}

// Clean remove all the data from the providers.
func (m *JWTToken) Clean() derrors.Error {
	return m.TokenProvider.Truncate()
}
