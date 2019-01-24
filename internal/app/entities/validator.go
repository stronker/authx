/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package entities

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-device-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-user-go"
)

const emptyOrganizationId = "organization_id cannot be empty"
const emptyEmail = "email cannot be empty"
const emptyName = "name cannot be empty"
const emptyRoleID = "role_id cannot be empty"
const emptyDeviceGroupId = "device_group_id cannot be empty"
const emptyDeviceId = "device_id cannot be empty"
const emptyRefreshToken = "refreshToken is mandatory"
const emptyToken = "token is mandatory"


func ValidOrganizationID(organizationID *grpc_organization_go.OrganizationId) derrors.Error {
	if organizationID.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	return nil
}

func ValidUserID(userID * grpc_user_go.UserId) derrors.Error {
	if userID.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if userID.Email == "" {
		return derrors.NewInvalidArgumentError(emptyEmail)
	}
	return nil
}

// -- Device Credentials -- //

func ValidAddDeviceGroupCredentials (addRequest * grpc_authx_go.AddDeviceGroupCredentialsRequest) derrors.Error {

	if addRequest.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if addRequest.DeviceGroupId == "" {
		return derrors.NewInvalidArgumentError(emptyDeviceGroupId)
	}
	return nil
}

func ValidUpdateDeviceGroupCredentialsRequest (request * grpc_authx_go.UpdateDeviceGroupCredentialsRequest) derrors.Error{
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.DeviceGroupId == "" {
		return derrors.NewInvalidArgumentError(emptyDeviceGroupId)
	}
	if !request.UpdateEnabled && !request.UpdateDeviceConnectivity {
		return derrors.NewInvalidArgumentError("enabled or default_device_connectivity flags must change")
	}
	return nil
}

func ValidDeviceGroupLoginRequest (request * grpc_authx_go.DeviceGroupLoginRequest) derrors.Error{

	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.DeviceGroupApiKey == "" {
		return derrors.NewInvalidArgumentError("device_group_api_key cannot be empty")
	}
	return nil

}

func ValidAddDeviceCredentials (addRequest * grpc_authx_go.AddDeviceCredentialsRequest) derrors.Error {

	if addRequest.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if addRequest.DeviceGroupId == "" {
		return derrors.NewInvalidArgumentError(emptyDeviceGroupId)
	}
	if addRequest.DeviceId == "" {
		return derrors.NewInvalidArgumentError("device_id cannot be empty")
	}
	return nil
}

func ValidUpdateDeviceCredentialsRequest (request * grpc_authx_go.UpdateDeviceCredentialsRequest) derrors.Error{

	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.DeviceGroupId == "" {
		return derrors.NewInvalidArgumentError(emptyDeviceGroupId)
	}
	if request.DeviceId == "" {
		return derrors.NewInvalidArgumentError(emptyDeviceId)
	}

	return nil
}

func ValidDeviceLoginRequest (request * grpc_authx_go.DeviceLoginRequest) derrors.Error{

	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.DeviceApiKey == "" {
		return derrors.NewInvalidArgumentError("device_group_api_key cannot be empty")
	}
	return nil

}

func ValidDeviceID (request * grpc_device_go.DeviceId) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.DeviceGroupId == "" {
		return derrors.NewInvalidArgumentError(emptyDeviceGroupId)
	}
	if request.DeviceId == "" {
		return derrors.NewInvalidArgumentError(emptyDeviceId)
	}
	return nil
}

func ValidDeviceGroupID (request * grpc_device_go.DeviceGroupId) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.DeviceGroupId == "" {
		return derrors.NewInvalidArgumentError(emptyDeviceGroupId)
	}
	return nil
}

func ValidRefreshToken (request * grpc_authx_go.RefreshTokenRequest) derrors.Error {
	if request.RefreshToken == "" {
		return derrors.NewInvalidArgumentError(emptyRefreshToken)
	}

	if request.Token == "" {
		return derrors.NewInvalidArgumentError(emptyToken)
	}
	return nil
}
