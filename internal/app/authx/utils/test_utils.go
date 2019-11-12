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

package utils

import (
	"github.com/google/uuid"
	"github.com/nalej/authx/internal/app/authx/entities"
	"os"
)

func RunIntegrationTests() bool {
	var runIntegration = os.Getenv("RUN_INTEGRATION_TEST")
	return runIntegration == "true"
}

type DeviceTestHelper struct {
}

func NewDeviceTestHepler() *DeviceTestHelper {
	return &DeviceTestHelper{}
}

func (d *DeviceTestHelper) CreateDeviceGroupCredentials() *entities.DeviceGroupCredentials {

	return &entities.DeviceGroupCredentials{
		OrganizationID:            uuid.New().String(),
		DeviceGroupID:             uuid.New().String(),
		DeviceGroupApiKey:         uuid.New().String(),
		Enabled:                   true,
		DefaultDeviceConnectivity: false,
		Secret:                    uuid.New().String(),
	}
}

func (d *DeviceTestHelper) CreateDeviceCredentials(group entities.DeviceGroupCredentials) *entities.DeviceCredentials {

	return &entities.DeviceCredentials{
		OrganizationID: group.OrganizationID,
		DeviceGroupID:  group.DeviceGroupID,
		DeviceID:       uuid.New().String(),
		DeviceApiKey:   uuid.New().String(),
		Enabled:        true,
	}
}
