package device

import (
	"fmt"
	"github.com/nalej/authx/internal/app/authx/entities"
	"github.com/nalej/derrors"
	"sync"
)

type MockupDeviceCredentialsProvider struct {
	sync.Mutex
	// deviceCredentials indexed by organizationId#deviceGroupId#deviceId
	deviceCredentials map[string]entities.DeviceCredentials
	deviceByApiKey    map[string]entities.DeviceCredentials
	// groupCredentials indexed by organizationId#deviceGroupId
	groupCredentials map[string]entities.DeviceGroupCredentials
	// groupCredentials indexed by device_group_api_key
	groupByApyKey map[string]entities.DeviceGroupCredentials
}

func NewMockupDeviceCredentialsProvider() Provider {
	return &MockupDeviceCredentialsProvider{
		deviceCredentials: make(map[string]entities.DeviceCredentials, 0),
		deviceByApiKey:    make(map[string]entities.DeviceCredentials, 0),
		groupCredentials:  make(map[string]entities.DeviceGroupCredentials, 0),
		groupByApyKey:     make(map[string]entities.DeviceGroupCredentials, 0),
	}
}

// ---------------------------
func GenerateGroupKey(organizationId string, deviceGroupId string) string {
	return fmt.Sprintf("%s#%s", organizationId, deviceGroupId)
}
func (m *MockupDeviceCredentialsProvider) unsafeExistsGroupCredentials(group_api_key string) bool {
	_, exists := m.groupCredentials[group_api_key]

	return exists
}
func (m *MockupDeviceCredentialsProvider) AddDeviceGroupCredentials(groupCredentials *entities.DeviceGroupCredentials) derrors.Error {

	m.Lock()
	defer m.Unlock()

	key := GenerateGroupKey(groupCredentials.OrganizationID, groupCredentials.DeviceGroupID)

	if !m.unsafeExistsGroupCredentials(key) {
		m.groupCredentials[key] = *groupCredentials
		m.groupByApyKey[groupCredentials.DeviceGroupApiKey] = *groupCredentials

	} else {
		return derrors.NewAlreadyExistsError("add device group credentials").WithParams(groupCredentials.OrganizationID, groupCredentials.DeviceGroupID)
	}

	return nil
}
func (m *MockupDeviceCredentialsProvider) UpdateDeviceGroupCredentials(groupCredentials *entities.DeviceGroupCredentials) derrors.Error {

	m.Lock()
	defer m.Unlock()

	key := GenerateGroupKey(groupCredentials.OrganizationID, groupCredentials.DeviceGroupID)

	if !m.unsafeExistsGroupCredentials(key) {
		return derrors.NewNotFoundError("device group credentials").WithParams(groupCredentials.OrganizationID, groupCredentials.DeviceGroupID)
	}
	m.groupCredentials[key] = *groupCredentials
	m.groupByApyKey[groupCredentials.DeviceGroupApiKey] = *groupCredentials

	return nil
}
func (m *MockupDeviceCredentialsProvider) ExistsDeviceGroup(organizationId string, deviceGroupId string) (bool, derrors.Error) {

	m.Lock()
	defer m.Unlock()

	key := GenerateGroupKey(organizationId, deviceGroupId)

	return m.unsafeExistsGroupCredentials(key), nil

}
func (m *MockupDeviceCredentialsProvider) GetDeviceGroup(organizationId string, deviceGroupId string) (*entities.DeviceGroupCredentials, derrors.Error) {

	m.Lock()
	defer m.Unlock()

	key := GenerateGroupKey(organizationId, deviceGroupId)

	group, exists := m.groupCredentials[key]
	if !exists {
		return nil, derrors.NewNotFoundError("device group credentials").WithParams(organizationId, deviceGroupId)
	}

	return &group, nil

}
func (m *MockupDeviceCredentialsProvider) GetDeviceGroupByApiKey(apiKey string) (*entities.DeviceGroupCredentials, derrors.Error) {

	m.Lock()
	defer m.Unlock()

	group, exists := m.groupByApyKey[apiKey]
	if !exists {
		return nil, derrors.NewNotFoundError("device group apiKey").WithParams(apiKey)
	}

	return &group, nil
}
func (m *MockupDeviceCredentialsProvider) RemoveDeviceGroup(organizationId string, deviceGroupId string) derrors.Error {

	m.Lock()
	defer m.Unlock()

	key := GenerateGroupKey(organizationId, deviceGroupId)

	group, exits := m.groupCredentials[key]
	if !exits {
		return derrors.NewNotFoundError("device group credentials").WithParams(organizationId, deviceGroupId)
	}

	delete(m.groupCredentials, key)
	delete(m.groupByApyKey, group.DeviceGroupApiKey)

	return nil
}
func (m *MockupDeviceCredentialsProvider) TruncateDeviceGroup() derrors.Error {
	m.groupCredentials = make(map[string]entities.DeviceGroupCredentials, 0)
	m.groupByApyKey = make(map[string]entities.DeviceGroupCredentials, 0)
	return nil
}

// ---------------------------

func GenerateDeviceKey(organizationId string, deviceGroupId string, deviceId string) string {
	return fmt.Sprintf("%s#%s#%s", organizationId, deviceGroupId, deviceId)
}
func (m *MockupDeviceCredentialsProvider) unsafeExistsDeviceCredentials(key string) bool {
	_, exits := m.deviceCredentials[key]
	return exits
}
func (m *MockupDeviceCredentialsProvider) AddDeviceCredentials(credentials *entities.DeviceCredentials) derrors.Error {
	m.Lock()
	defer m.Unlock()

	groupKey := GenerateGroupKey(credentials.OrganizationID, credentials.DeviceGroupID)
	if !m.unsafeExistsGroupCredentials(groupKey) {
		return derrors.NewNotFoundError("device group").WithParams(credentials.OrganizationID, credentials.DeviceGroupID)
	}

	deviceKey := GenerateDeviceKey(credentials.OrganizationID, credentials.DeviceGroupID, credentials.DeviceID)

	if !m.unsafeExistsDeviceCredentials(deviceKey) {
		m.deviceCredentials[deviceKey] = *credentials
		m.deviceByApiKey[credentials.DeviceApiKey] = *credentials
	} else {
		return derrors.NewAlreadyExistsError("device credentials").WithParams(credentials.OrganizationID,
			credentials.DeviceGroupID, credentials.DeviceID)
	}
	return nil

}
func (m *MockupDeviceCredentialsProvider) UpdateDeviceCredentials(credentials *entities.DeviceCredentials) derrors.Error {
	m.Lock()
	defer m.Unlock()

	deviceKey := GenerateDeviceKey(credentials.OrganizationID, credentials.DeviceGroupID, credentials.DeviceID)

	if !m.unsafeExistsDeviceCredentials(deviceKey) {
		return derrors.NewNotFoundError("device credentials").WithParams(credentials.OrganizationID,
			credentials.DeviceGroupID, credentials.DeviceID)
	} else {
		m.deviceCredentials[deviceKey] = *credentials
		m.deviceByApiKey[credentials.DeviceApiKey] = *credentials
	}
	return nil
}
func (m *MockupDeviceCredentialsProvider) ExistsDevice(organizationId string, deviceGroupId string, deviceId string) (bool, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	deviceKey := GenerateDeviceKey(organizationId, deviceGroupId, deviceId)

	return m.unsafeExistsDeviceCredentials(deviceKey), nil
}
func (m *MockupDeviceCredentialsProvider) GetDevice(organizationId string, deviceGroupId string, deviceId string) (*entities.DeviceCredentials, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	deviceKey := GenerateDeviceKey(organizationId, deviceGroupId, deviceId)

	device, exists := m.deviceCredentials[deviceKey]
	if !exists {
		return nil, derrors.NewNotFoundError("device group credentials").WithParams(organizationId, deviceGroupId, deviceId)
	}

	return &device, nil

}
func (m *MockupDeviceCredentialsProvider) GetDeviceByApiKey(apiKey string) (*entities.DeviceCredentials, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	device, exits := m.deviceByApiKey[apiKey]
	if !exits {
		return nil, derrors.NewNotFoundError("device credentials by api Key").WithParams(apiKey)
	}

	return &device, nil
}
func (m *MockupDeviceCredentialsProvider) RemoveDevice(organizationId string, deviceGroupId string, deviceId string) derrors.Error {
	m.Lock()
	defer m.Unlock()

	deviceKey := GenerateDeviceKey(organizationId, deviceGroupId, deviceId)

	device, exits := m.deviceCredentials[deviceKey]
	if !exits {
		return derrors.NewNotFoundError("device credentials").WithParams(organizationId, deviceGroupId, deviceId)
	}

	delete(m.deviceCredentials, deviceKey)
	delete(m.deviceByApiKey, device.DeviceApiKey)
	return nil
}
func (m *MockupDeviceCredentialsProvider) TruncateDevice() {
	m.deviceCredentials = make(map[string]entities.DeviceCredentials, 0)
	m.deviceByApiKey = make(map[string]entities.DeviceCredentials, 0)
}

// -------------

func (m *MockupDeviceCredentialsProvider) Truncate() derrors.Error {
	m.TruncateDevice()
	m.TruncateDeviceGroup()
	return nil
}
