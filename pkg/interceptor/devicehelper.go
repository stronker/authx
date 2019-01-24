/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package interceptor

import (
	"context"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-authx-go"
	"google.golang.org/grpc/metadata"
)

const DeviceIdField = "device_id"
const DeviceGroupIdField = "device_group_id"

type DeviceRequestMetadata struct{
	OrganizationID         string
	DeviceGroupID string
	DeviceID string
	DevicePrimitive bool
}

// GetDeviceRequestMetadata extracts the request metadata from the context so that it
// can be easily consumed by upper layers.
func GetDeviceRequestMetadata(ctx context.Context) (*DeviceRequestMetadata, derrors.Error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, derrors.NewInvalidArgumentError("expecting JWT metadata")
	}
	deviceID, found := md[DeviceIdField]
	if !found {
		return nil, derrors.NewUnauthenticatedError("device_id not found")
	}
	deviceGroupID, found := md[DeviceGroupIdField]
	if !found {
		return nil, derrors.NewUnauthenticatedError("device_group_id not found")
	}
	organizationID, found := md[OrganizationIdField]
	if !found {
		return nil, derrors.NewUnauthenticatedError("organizationID not found")
	}
	_, devicePrimitive := md[grpc_authx_go.AccessPrimitive_DEVICE.String()]

	return &DeviceRequestMetadata{
		OrganizationID:  organizationID[0],
		DeviceGroupID:   deviceGroupID[0],
		DeviceID:        deviceID[0],
		DevicePrimitive: devicePrimitive,
	}, nil
}
