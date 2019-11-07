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

package entities

// TokenData is the information that the system stores.
type DeviceTokenData struct {
	DeviceId       string `cql:"device_id"`
	TokenID        string `cql:"token_id"`
	RefreshToken   string `cql:"refresh_token"`
	ExpirationDate int64  `cql:"expiration_date"`
	OrganizationId string `cql:"organization_id"`
	DeviceGroupId  string `cql:"e¡device_group_id"`
}

// NewTokenData creates an instance of the structure
func NewDeviceTokenData(deviceID string, tokenID string, refreshToken string,
	expirationDate int64, organizationID string, deviceGroupID string) *DeviceTokenData {

	return &DeviceTokenData{
		DeviceId:       deviceID,
		TokenID:        tokenID,
		RefreshToken:   refreshToken,
		ExpirationDate: expirationDate,
		OrganizationId: organizationID,
		DeviceGroupId:  deviceGroupID,
	}
}
