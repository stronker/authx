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

package device

import (
	"github.com/nalej/derrors"
	"github.com/stronker/authx/internal/app/authx/entities"
)

type Provider interface {
	// AddDeviceGroupCredentials adds credentials of a device group
	AddDeviceGroupCredentials(*entities.DeviceGroupCredentials) derrors.Error
	// UpdateDeviceGroupCredentials updates a device group (Enabled and/or DefaultDeviceConnectivity" flags)
	UpdateDeviceGroupCredentials(*entities.DeviceGroupCredentials) derrors.Error
	// ExistsDeviceGroup checks if a group exists
	ExistsDeviceGroup(organizationId string, deviceGroupId string) (bool, derrors.Error)
	// GetDeviceGroup retrieves a device group credentials
	GetDeviceGroup(organizationId string, deviceGroupId string) (*entities.DeviceGroupCredentials, derrors.Error)
	// GetDeviceGroupByApiKey retrieves a device group credentials by GroupApiKey
	GetDeviceGroupByApiKey(deviceApiKey string) (*entities.DeviceGroupCredentials, derrors.Error)
	// RemoveDeviceGroup removes a device group
	RemoveDeviceGroup(organizationId string, deviceGroupId string) derrors.Error
	
	// Truncate removes all stored devices and device groups
	Truncate() derrors.Error
	
	// AddDeviceCredentials adds credentials of a device
	AddDeviceCredentials(*entities.DeviceCredentials) derrors.Error
	// UpdateDeviceCredentials updates a device credentials (Enable flag)
	UpdateDeviceCredentials(*entities.DeviceCredentials) derrors.Error
	// ExistsDevice checks if a device exists
	ExistsDevice(organizationId string, deviceGroupId string, deviceId string) (bool, derrors.Error)
	// GetDevice retrieves a device credentials
	GetDevice(organizationId string, deviceGroupId string, deviceId string) (*entities.DeviceCredentials, derrors.Error)
	// GetDeviceByApiKey retrieves a device credentials by apiKey
	GetDeviceByApiKey(deviceApiKey string) (*entities.DeviceCredentials, derrors.Error)
	// RemoveDevice removes credentials from a device
	RemoveDevice(organizationId string, deviceGroupId string, deviceId string) derrors.Error
}
