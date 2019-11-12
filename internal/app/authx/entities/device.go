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

import (
	"github.com/google/uuid"
	"github.com/nalej/grpc-authx-go"
)

type DeviceGroupCredentials struct {
	OrganizationID            string
	DeviceGroupID             string
	DeviceGroupApiKey         string
	Enabled                   bool
	DefaultDeviceConnectivity bool
	Secret                    string
}

func NewDeviceGroupCredentials(organizationId string, deviceGroupId string, deviceGroupApiKey string,
	enabled bool, defaultConnectivity bool) *DeviceGroupCredentials {
	return &DeviceGroupCredentials{
		OrganizationID:            organizationId,
		DeviceGroupID:             deviceGroupId,
		DeviceGroupApiKey:         deviceGroupApiKey,
		Enabled:                   enabled,
		DefaultDeviceConnectivity: defaultConnectivity,
	}
}

func NewDeviceGroupCredentialsFromGRPC(addRequest *grpc_authx_go.AddDeviceGroupCredentialsRequest) *DeviceGroupCredentials {
	return &DeviceGroupCredentials{
		OrganizationID:            addRequest.OrganizationId,
		DeviceGroupID:             addRequest.DeviceGroupId,
		DeviceGroupApiKey:         uuid.New().String(),
		Enabled:                   addRequest.Enabled,
		DefaultDeviceConnectivity: addRequest.DefaultDeviceConnectivity,
		Secret:                    uuid.New().String(), // secret of the device_group
	}
}

func (dg *DeviceGroupCredentials) ToGRPC() *grpc_authx_go.DeviceGroupCredentials {
	return &grpc_authx_go.DeviceGroupCredentials{
		OrganizationId:            dg.OrganizationID,
		DeviceGroupId:             dg.DeviceGroupID,
		DeviceGroupApiKey:         dg.DeviceGroupApiKey,
		Enabled:                   dg.Enabled,
		DefaultDeviceConnectivity: dg.DefaultDeviceConnectivity,
	}
}

// ----------------------- //
// -- DeviceCredentials -- //
// ----------------------- //

type DeviceCredentials struct {
	OrganizationID string
	DeviceGroupID  string
	DeviceID       string
	DeviceApiKey   string
	Enabled        bool
}

func NewDeviceCredentials(organizationId string, deviceGroupId string, deviceId string,
	enabled bool, deviceGroupApiKey string) *DeviceCredentials {
	return &DeviceCredentials{
		OrganizationID: organizationId,
		DeviceGroupID:  deviceGroupId,
		DeviceID:       deviceId,
		DeviceApiKey:   deviceGroupApiKey,
		Enabled:        enabled,
	}
}

func NewDeviceCredentialsFromGRPC(addRequest *grpc_authx_go.AddDeviceCredentialsRequest) *DeviceCredentials {
	return &DeviceCredentials{
		OrganizationID: addRequest.OrganizationId,
		DeviceGroupID:  addRequest.DeviceGroupId,
		DeviceID:       addRequest.DeviceId,
		DeviceApiKey:   uuid.New().String(),
	}
}

func (dg *DeviceCredentials) ToGRPC() *grpc_authx_go.DeviceCredentials {
	return &grpc_authx_go.DeviceCredentials{
		OrganizationId: dg.OrganizationID,
		DeviceGroupId:  dg.DeviceGroupID,
		DeviceId:       dg.DeviceID,
		DeviceApiKey:   dg.DeviceApiKey,
		Enabled:        dg.Enabled,
	}
}
