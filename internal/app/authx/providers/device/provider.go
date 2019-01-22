package device

import (
	"github.com/nalej/authx/internal/app/authx/entities"
	"github.com/nalej/derrors"
)

type Provider interface {

	AddDeviceGroupCredentials (*entities.DeviceGroupCredentials) derrors.Error
	UpdateDeviceGroupCredentials (*entities.DeviceGroupCredentials) derrors.Error
	ExistsDeviceGroup (organizationId string, deviceGroupId string) (bool, derrors.Error)
	GetDeviceGroup (organizationId string, deviceGroupId string) (* entities.DeviceGroupCredentials, derrors.Error)
	GetDeviceGroupByApiKey(deviceApiKey string)  (* entities.DeviceGroupCredentials, derrors.Error)
	RemoveDeviceGroup (organizationId string, deviceGroupId string) derrors.Error
	Clear() derrors.Error

	AddDeviceCredentials (* entities.DeviceCredentials) derrors.Error
	UpdateDeviceCredentials (* entities.DeviceCredentials) derrors.Error
	ExistsDevice (organizationId string, deviceGroupId string, deviceId string ) (bool, derrors.Error)
	GetDevice (organizationId string, deviceGroupId string, deviceId string) (* entities.DeviceCredentials, derrors.Error)
	GetDeviceByApiKey(deviceApiKey string)  (* entities.DeviceCredentials, derrors.Error)
	RemoveDevice (organizationId string, deviceGroupId string, deviceId string) derrors.Error

}
