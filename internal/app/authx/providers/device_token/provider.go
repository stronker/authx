package device_token

import (
	"github.com/nalej/authx/internal/app/authx/entities"
	"github.com/nalej/derrors"
)

type Provider interface {
	// Delete an existing token.
	Delete(deviceID string, tokenID string) derrors.Error
	// Add a token.
	Add(token *entities.DeviceTokenData) derrors.Error
	// Get an existing token.
	Get(deviceID string, tokenID string) (*entities.DeviceTokenData, derrors.Error)
	// Exist checks if the token was added.
	Exist(deviceID string, tokenID string) (*bool, derrors.Error)
	// Update an existing token
	Update(token *entities.DeviceTokenData) derrors.Error
	// Truncate cleans all data.
	Truncate() derrors.Error

	// Get an existing token.
	GetByRefreshToken(refreshToken string) (*entities.DeviceTokenData, derrors.Error)

	DeleteExpiredTokens() derrors.Error

}
