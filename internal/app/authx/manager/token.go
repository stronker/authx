/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package manager

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/nalej/authx/internal/app/authx/entities"
	nalejToken "github.com/nalej/authx/internal/app/authx/providers/token"
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
		secret string, update bool) (*GeneratedToken, derrors.Error)
	// Refresh renew an old token.
	Refresh(oldToken string, refreshToken string,
		expirationPeriod time.Duration, secret string) (*GeneratedToken, derrors.Error)
	// Clean remove all the data from the providers.
	Clean() derrors.Error
}

// JWTToken is an implementation of token using JWT.
type JWTToken struct {
	TokenProvider nalejToken.Token
	Password      Password
}

// NewJWTToken create a new instance of JWTToken
func NewJWTToken(tokenProvider nalejToken.Token, password Password) Token {
	return &JWTToken{TokenProvider: tokenProvider, Password: password}

}

// NewJWTTokenMockup create a new mockup of JWTToken
func NewJWTTokenMockup() Token {
	return NewJWTToken(nalejToken.NewTokenMockup(), NewBCryptPassword())
}

// Generate a new JWT token with the personal claim.
func (m *JWTToken) Generate(personalClaim *token.PersonalClaim, expirationPeriod time.Duration,
	secret string, update bool) (*GeneratedToken, derrors.Error) {

	claim := token.NewClaim(*personalClaim, Issuer, time.Now(), expirationPeriod)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := t.SignedString([]byte(secret))
	if err != nil {
		return nil, derrors.NewInternalError("impossible generate JWT token", err)
	}

	refreshToken := token.GenerateUUID()
	hashedRefreshToken, err := m.Password.GenerateHashedPassword(refreshToken)
	if err != nil {
		return nil, derrors.NewInternalError("impossible generate RefreshToken", err)
	}
	tokenData := entities.NewTokenData(claim.UserID, claim.Id, hashedRefreshToken, claim.ExpiresAt)
	err = m.TokenProvider.Add(tokenData)

	if err != nil {
		return nil, derrors.NewInternalError("impossible store RefreshToken", err)
	}
	gToken := NewGeneratedToken(tokenString, refreshToken)
	return gToken, nil
}

// Refresh renew an old token.
func (m *JWTToken) Refresh(oldToken string, refreshToken string,
	expirationPeriod time.Duration, secret string) (*GeneratedToken, derrors.Error) {

	tk, jwtErr := jwt.ParseWithClaims(oldToken, &token.Claim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if jwtErr != nil {
		return nil, derrors.NewUnauthenticatedError("impossible recover RefreshToken", jwtErr)
	}

	cl, ok := tk.Claims.(*token.Claim)
	if !ok {
		return nil, derrors.NewUnauthenticatedError("impossible recover token")
	}
	username := cl.UserID
	tokenID:= cl.Id

	tokenData, err := m.TokenProvider.Get(username, tokenID)
	if err != nil {
		return nil, derrors.NewUnauthenticatedError("impossible recover RefreshToken", err)
	}
	ts := time.Now().Unix()
	if tokenData == nil || ts > tokenData.ExpirationDate {
		return nil, derrors.NewUnauthenticatedError("the refresh token is expired")
	}

	err = m.Password.CompareHashAndPassword(tokenData.RefreshToken, refreshToken)
	if err != nil {
		return nil, derrors.NewUnauthenticatedError("the refresh token is not valid", err)
	}

	gt, err := m.Generate(&cl.PersonalClaim, expirationPeriod, secret, true)
	if err != nil {
		return nil, derrors.NewInternalError("impossible create new token", err)
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
