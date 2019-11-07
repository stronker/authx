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

package token

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/nalej/grpc-authx-go"
	"time"
)

// Claim joins the personal claim and the standard JWT claim.
type DeviceClaim struct {
	jwt.StandardClaims
	OrganizationID string   `json:"organizationID,omitempty"`
	DeviceGroupID  string   `json:"device_group_id,omitempty"`
	DeviceID       string   `json:"device_id,omitempty"`
	Primitives     []string `json:"access,omitempty"`
}

func NewDeviceClaim(organizationId string, deviceGroupId string, deviceId string, expirationPeriod time.Duration) *DeviceClaim {
	stdClaim := jwt.StandardClaims{
		Issuer:    "authx",
		Id:        GenerateUUID(),
		NotBefore: time.Now().Unix(),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(expirationPeriod).Unix(),
	}

	return &DeviceClaim{
		StandardClaims: stdClaim,
		OrganizationID: organizationId,
		DeviceGroupID:  deviceGroupId,
		DeviceID:       deviceId,
		Primitives:     []string{grpc_authx_go.AccessPrimitive_name[int32(grpc_authx_go.AccessPrimitive_DEVICE)]},
	}
}
