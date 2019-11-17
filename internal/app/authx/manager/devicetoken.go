/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package manager

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/nalej/authx/pkg/token"
	"github.com/nalej/derrors"
	"github.com/rs/zerolog/log"
	"github.com/stronker/authx/internal/app/authx/entities"
	"github.com/stronker/authx/internal/app/authx/providers/device"
	"github.com/stronker/authx/internal/app/authx/providers/device_token"
	"time"
)

// Token is a interface manages the business logic of tokens.
type DeviceToken interface {
	// Generate a new token with the device claim.
	Generate(deviceClaim *token.DeviceClaim, expirationPeriod time.Duration, secret string) (*GeneratedToken, derrors.Error)
	// Refresh renew an old token.
	Refresh(oldToken string, refreshToken string,
		expirationPeriod time.Duration, secret string) (*GeneratedToken, derrors.Error)
	// Gets the deviceClaim of a deviceToken
	GetTokenInfo(tokenInfo string, secret string) (*token.DeviceClaim, derrors.Error)
	// Clean remove all the data from the providers.
	Clean() derrors.Error
}

type JWTDeviceToken struct {
	DeviceProvider      device.Provider // device Provider
	DeviceTokenProvider device_token.Provider
}

// NewJWTToken create a new instance of JWTToken
func NewJWTDeviceToken(deviceProvider device.Provider, tokenProvider device_token.Provider) DeviceToken {
	return &JWTDeviceToken{
		DeviceProvider:      deviceProvider,
		DeviceTokenProvider: tokenProvider}
	
}

// NewJWTTokenMockup create a new mockup of JWTToken
func NewJWTDeviceTokenMockup() DeviceToken {
	return NewJWTDeviceToken(device.NewMockupDeviceCredentialsProvider(),
		device_token.NewDeviceTokenMockup())
}

// Generate a new JWT token with the personal claim.
func (m *JWTDeviceToken) Generate(deviceClaim *token.DeviceClaim, expirationPeriod time.Duration,
	secret string) (*GeneratedToken, derrors.Error) {
	
	newClaim := token.NewDeviceClaim(deviceClaim.OrganizationID, deviceClaim.DeviceGroupID, deviceClaim.DeviceID, expirationPeriod)
	
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaim)
	tokenString, err := t.SignedString([]byte(secret))
	if err != nil {
		return nil, derrors.NewInternalError("impossible generate JWT Device token", err)
	}
	
	refreshToken := token.GenerateUUID()
	
	tokenData := entities.NewDeviceTokenData(newClaim.DeviceID, newClaim.Id, refreshToken,
		newClaim.ExpiresAt, newClaim.OrganizationID, newClaim.DeviceGroupID)
	
	err = m.DeviceTokenProvider.Add(tokenData)
	if err != nil {
		return nil, derrors.NewInternalError("impossible store RefreshToken", err)
	}
	gToken := NewGeneratedToken(tokenString, refreshToken)
	return gToken, nil
}

func (m *JWTDeviceToken) GetTokenInfo(tokenInfo string, secret string) (*token.DeviceClaim, derrors.Error) {
	
	tk, jwtErr := jwt.ParseWithClaims(tokenInfo, &token.DeviceClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if jwtErr != nil {
		return nil, derrors.NewUnauthenticatedError("impossible recover token", jwtErr)
	}
	
	cl, ok := tk.Claims.(*token.DeviceClaim)
	if !ok {
		return nil, derrors.NewUnauthenticatedError("impossible recover device token")
	}
	return cl, nil
	
}

// Refresh renew an old token.
func (m *JWTDeviceToken) Refresh(oldToken string, refreshToken string,
	expirationPeriod time.Duration, secret string) (*GeneratedToken, derrors.Error) {
	
	dToken, err := m.DeviceTokenProvider.GetByRefreshToken(refreshToken)
	if err != nil {
		return nil, derrors.NewInternalError("error getting token info", err)
	}
	
	group, err := m.DeviceProvider.GetDeviceGroup(dToken.OrganizationId, dToken.DeviceGroupId)
	if err != nil {
		return nil, err
	}
	
	tk, jwtErr := jwt.ParseWithClaims(oldToken, &token.DeviceClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(group.Secret), nil
	})
	if jwtErr != nil {
		return nil, derrors.NewUnauthenticatedError("impossible recover RefreshToken", jwtErr)
	}
	
	cl, ok := tk.Claims.(*token.DeviceClaim)
	if !ok {
		return nil, derrors.NewUnauthenticatedError("impossible recover device token")
	}
	
	deviceID := cl.DeviceID
	tokenID := cl.Id
	
	tokenData, err := m.DeviceTokenProvider.Get(deviceID, tokenID)
	if err != nil {
		return nil, derrors.NewUnauthenticatedError("impossible recover RefreshToken", err)
	}
	ts := time.Now().Unix()
	if tokenData == nil || ts > tokenData.ExpirationDate {
		return nil, derrors.NewUnauthenticatedError("the refresh token is expired")
	}
	
	gt, err := m.Generate(cl, expirationPeriod, secret)
	if err != nil {
		return nil, derrors.NewInternalError("impossible create new token", err)
	}
	
	err = m.DeviceTokenProvider.Delete(deviceID, tokenID)
	if err != nil {
		log.Warn().Err(err).Msg("impossible delete refresh token")
	}
	
	return gt, nil
}

// Clean remove all the data from the providers.
func (m *JWTDeviceToken) Clean() derrors.Error {
	return m.DeviceTokenProvider.Truncate()
}
