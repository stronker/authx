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
		DeviceId: 		deviceID,
		TokenID:        tokenID,
		RefreshToken:   refreshToken,
		ExpirationDate: expirationDate,
		OrganizationId: organizationID,
		DeviceGroupId: 	deviceGroupID,
	}
}

