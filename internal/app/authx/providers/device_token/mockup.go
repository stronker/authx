package device_token

import (
	"fmt"
	"github.com/nalej/authx/internal/app/authx/entities"
	"github.com/nalej/derrors"
	"sync"
	"time"
)

// DeviceTokenMockup is an in-memory mockup.
type DeviceTokenMockup struct {
	sync.Mutex
	data map[string]entities.DeviceTokenData
}

// NewTokenMockup create a new instance of TokenMockup.
func NewDeviceTokenMockup() Provider {
	return &DeviceTokenMockup{
		data: make(map[string]entities.DeviceTokenData,0),
	}
}
func (m *DeviceTokenMockup) unsafeExists (deviceID string, tokenID string) bool {
	_, ok := m.data[m.generateID(deviceID, tokenID)]
	return ok
}

func (m *DeviceTokenMockup) unsafeGet(deviceID string, tokenID string) (*entities.DeviceTokenData, derrors.Error) {

	data, ok := m.data[m.generateID(deviceID, tokenID)]
	if !ok {
		return nil, derrors.NewNotFoundError("device token not found").WithParams(deviceID, tokenID)
	}
	return &data, nil
}

func (m *DeviceTokenMockup) generateID(deviceID string, tokenID string) string {
	return fmt.Sprintf("%s:%s", deviceID, tokenID)
}

// Delete an existing token.
func (m *DeviceTokenMockup) Delete(deviceID string, tokenID string) derrors.Error{
	m.Lock()
	defer m.Unlock()

	id := m.generateID(deviceID, tokenID)
	_, err := m.unsafeGet(deviceID, tokenID)
	if err != nil {
		return derrors.NewNotFoundError("device not found").WithParams(deviceID)
	}

	delete(m.data, id)
	return nil
}
// Add a token.
func (m *DeviceTokenMockup) Add(token *entities.DeviceTokenData) derrors.Error {
	m.Lock()
	defer m.Unlock()

	if m.unsafeExists(token.DeviceId, token.TokenID){
		return derrors.NewAlreadyExistsError("device token").WithParams(token.DeviceId, token.TokenID)
	}
	m.data[m.generateID(token.DeviceId, token.TokenID)] = *token
	return nil
}
// Get an existing token.
func (m *DeviceTokenMockup) Get(deviceID string, tokenID string) (*entities.DeviceTokenData, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	data, ok := m.data[m.generateID(deviceID, tokenID)]
	if !ok {
		return nil, derrors.NewNotFoundError("device token not found").WithParams(deviceID, tokenID)
	}
	return &data, nil
}
// Exist checks if the token was added.
func (m *DeviceTokenMockup) Exist(deviceID string, tokenID string) (*bool, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	_, ok := m.data[m.generateID(deviceID, tokenID)]
	return &ok, nil
}
// Update an existing token
func (m *DeviceTokenMockup) Update(token *entities.DeviceTokenData) derrors.Error {
	m.Lock()
	defer m.Unlock()

	if ! m.unsafeExists(token.DeviceId, token.TokenID){
		return  derrors.NewNotFoundError("device token").WithParams(token.DeviceId, token.TokenID)
	}
	m.data[m.generateID(token.DeviceId, token.TokenID)] = *token
	return nil
}
// Truncate cleans all data.
func (m *DeviceTokenMockup) Truncate() derrors.Error{
	m.Lock()
	defer m.Unlock()

	m.data = make(map[string]entities.DeviceTokenData,0)
	return nil
}

func (m *DeviceTokenMockup) DeleteExpiredTokens() derrors.Error{
	m.Lock()
	defer m.Unlock()

	idBorrow := make([]string, 0)

	for _, token := range m.data{
		if token.ExpirationDate < time.Now().Unix() {
			id := m.generateID(token.DeviceId, token.TokenID)
			idBorrow = append(idBorrow, id)
		}

	}
	for _, id := range idBorrow{
		delete(m.data, id)
	}
	return nil
}
