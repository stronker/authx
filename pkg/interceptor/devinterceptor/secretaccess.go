/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package devinterceptor

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-device-go"
)

// SecretAccess defines an interface to facilitate accessing group secrets.
type SecretAccess interface {
	// Connect to the appropriate backend.
	Connect() derrors.Error
	// AddDeviceGroupCredentials adds credentials of a device group
	RetrieveSecret(id *grpc_device_go.DeviceGroupId) (string, derrors.Error)
}
