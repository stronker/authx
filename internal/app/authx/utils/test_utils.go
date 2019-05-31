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
