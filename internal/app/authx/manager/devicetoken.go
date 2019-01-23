/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package manager

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/nalej/authx/internal/app/authx/entities"
	tokenProvider "github.com/nalej/authx/internal/app/authx/providers/token"
	"github.com/nalej/authx/pkg/token"
	"github.com/nalej/derrors"
	"time"
)

// Token is a interface manages the business logic of tokens.
type DeviceToken interface {
	// Generate a new token with the device claim.
	Generate(deviceClaim *token.DeviceClaim, expirationPeriod time.Duration, secret string, update bool) (*GeneratedToken, derrors.Error)
	// Refresh renew an old token.
	Refresh(oldToken string, refreshToken string,
		expirationPeriod time.Duration, secret string) (*GeneratedToken, derrors.Error)
	// Clean remove all the data from the providers.
	Clean() derrors.Error
}


type JWTDeviceToken struct {
	TokenProvider tokenProvider.Token
	Password      Password
}

// NewJWTToken create a new instance of JWTToken
func NewJWTDeviceToken(tokenProvider tokenProvider.Token, password Password) DeviceToken {
	return &JWTDeviceToken{TokenProvider: tokenProvider, Password: password}

}

// NewJWTTokenMockup create a new mockup of JWTToken
func NewJWTDeviceTokenMockup() DeviceToken {
	return &JWTDeviceToken{TokenProvider: tokenProvider.NewTokenMockup(), Password: NewBCryptPassword()}
}

// Generate a new JWT token with the personal claim.
func (m *JWTDeviceToken) Generate(deviceClaim *token.DeviceClaim, expirationPeriod time.Duration,
	secret string, update bool) (*GeneratedToken, derrors.Error) {

	deviceClaim.ExpiresAt = time.Now().Add(expirationPeriod).Unix()

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, deviceClaim)
	tokenString, err := t.SignedString([]byte(secret))
	if err != nil {
		return nil, derrors.NewInternalError("impossible generate JWT Device token", err)
	}

	refreshToken := token.GenerateUUID()
	hashedRefreshToken, err := m.Password.GenerateHashedPassword(refreshToken)
	if err != nil {
		return nil, derrors.NewInternalError("impossible generate RefreshToken", err)
	}
	tokenData := entities.NewTokenData(deviceClaim.DeviceID, deviceClaim.Id, hashedRefreshToken, deviceClaim.ExpiresAt)
	if ! update {
		err = m.TokenProvider.Add(tokenData)
	}else{
		err = m.TokenProvider.Update(tokenData)
	}
	if err != nil {
		return nil, derrors.NewInternalError("impossible store RefreshToken", err)
	}
	gToken := NewGeneratedToken(tokenString, refreshToken)
	return gToken, nil
}

// Refresh renew an old token.
func (m *JWTDeviceToken) Refresh(oldToken string, refreshToken string,
	expirationPeriod time.Duration, secret string) (*GeneratedToken, derrors.Error) {

	tk, jwtErr := jwt.ParseWithClaims(oldToken, &token.DeviceClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if jwtErr != nil {
		return nil, derrors.NewUnauthenticatedError("impossible recover RefreshToken", jwtErr)
	}

	cl, ok := tk.Claims.(*token.DeviceClaim)
	if !ok {
		return nil, derrors.NewUnauthenticatedError("impossible recover device token")
	}
	deviceID := cl.DeviceID
	tokenID:= cl.Id

	tokenData, err := m.TokenProvider.Get(deviceID, tokenID)
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

	gt, err := m.Generate(cl, expirationPeriod, secret, true)
	if err != nil {
		return nil, derrors.NewInternalError("impossible create new token", err)
	}

	return gt, nil
}

// Clean remove all the data from the providers.
func (m *JWTDeviceToken) Clean() derrors.Error {
	return m.TokenProvider.Truncate()
}

