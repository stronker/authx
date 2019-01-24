package device

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/nalej/authx/internal/app/authx/entities"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"sync"
)


const (
	deviceGroupCredentialsTable = "devicegroupcredentials"
	deviceCredentialsTable = "devicecredentials"
	rowNotFound = "not found"
)


type ScyllaDeviceCredentialsProvider struct {
	Address string
	Port int
	KeySpace string
	sync.Mutex
	Session *gocql.Session

}

func (sp *ScyllaDeviceCredentialsProvider) connect() derrors.Error {

	// connect to the cluster
	conf := gocql.NewCluster(sp.Address)
	conf.Keyspace = sp.KeySpace
	conf.Port = sp.Port

	session, err := conf.CreateSession()
	if err != nil {
		log.Error().Str("provider", "ScyllaDeviceCredentialsProvider").Str("trace", conversions.ToDerror(err).DebugReport()).Msg("unable to connect")
		return derrors.AsError(err, "cannot create session")
	}

	sp.Session = session
	return nil
}

func (sp *ScyllaDeviceCredentialsProvider) Disconnect()  {

	sp.Lock()
	defer sp.Unlock()

	if sp.Session != nil {
		sp.Session.Close()
		sp.Session = nil
	}
}

func (sp *ScyllaDeviceCredentialsProvider) checkConnectionAndConnect () derrors.Error {

	if sp.Session != nil {
		return nil
	}
	log.Info().Str("provider", "ScyllaDeviceCredentialsProvider"). Msg("session not connected, trying to connect it!")
	err := sp.connect()
	if err != nil {
		return err
	}

	return nil
}

func NewScyllaDeviceCredentialsProvider (address string, port int, keyspace string) *ScyllaDeviceCredentialsProvider{
	provider := ScyllaDeviceCredentialsProvider{Address:address, Port:port, KeySpace:keyspace}
	provider.connect()
	return &provider
}

// -------------------------------

func (sp * ScyllaDeviceCredentialsProvider) unsafeExistsGroupCredentials (organizationId string, deviceGroupId string) (bool, derrors.Error) {
	if err := sp.checkConnectionAndConnect(); err != nil{
		return false, err
	}

	var returnedId string

	stmt, names := qb.Select(deviceGroupCredentialsTable).Columns("organization_id").
		Where(qb.Eq("organization_id")).Where(qb.Eq("device_group_id")).
		ToCql()

	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationId,
		"device_group_id": deviceGroupId})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		}else{
			return false, derrors.AsError(err, "cannot determinate if device group credentials exists")
		}
	}

	return true, nil
}
func (sp * ScyllaDeviceCredentialsProvider) AddDeviceGroupCredentials (groupCredentials * entities.DeviceGroupCredentials) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkConnectionAndConnect(); err != nil {
		return err
	}

	exists, err := sp.unsafeExistsGroupCredentials(groupCredentials.OrganizationID, groupCredentials.DeviceGroupID)
	if err != nil {
		return err
	}
	if exists {
		return  derrors.NewAlreadyExistsError("device group credentials").WithParams(groupCredentials.OrganizationID, groupCredentials.DeviceGroupID)
	}

	// add new basic credential
	stmt, names := qb.Insert(deviceGroupCredentialsTable).Columns("organization_id", "device_group_id",
		"device_group_api_key","enabled","default_device_connectivity").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(groupCredentials)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot add new device group credentials")
	}

	return nil

}
func (sp * ScyllaDeviceCredentialsProvider) UpdateDeviceGroupCredentials (groupCredentials * entities.DeviceGroupCredentials) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkConnectionAndConnect(); err != nil {
		return err
	}

	exists, err := sp.unsafeExistsGroupCredentials(groupCredentials.OrganizationID, groupCredentials.DeviceGroupID)
	if err != nil {
		return err
	}
	if ! exists {
		return  derrors.NewNotFoundError("device group credentials").WithParams(groupCredentials.OrganizationID, groupCredentials.DeviceGroupID)
	}

	// add new basic credential
	stmt, names := qb.Update(deviceGroupCredentialsTable).Set("enabled","default_device_connectivity").
		Where(qb.Eq("organization_id")).Where(qb.Eq("device_group_id")).
	 	ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(groupCredentials)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot update  device group credentials")
	}

	return  nil
}
func (sp * ScyllaDeviceCredentialsProvider) ExistsDeviceGroup (organizationId string, deviceGroupId string) (bool, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	return sp.unsafeExistsGroupCredentials(organizationId, deviceGroupId)

}
func (sp * ScyllaDeviceCredentialsProvider) GetDeviceGroup (organizationId string, deviceGroupId string) (* entities.DeviceGroupCredentials, derrors.Error){

	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkConnectionAndConnect(); err != nil{
		return nil, err
	}

	var deviceGroup entities.DeviceGroupCredentials

	stmt, names := qb.Select(deviceGroupCredentialsTable).
		Where(qb.Eq("organization_id")).Where(qb.Eq("device_group_id")).
		ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationId,
		"device_group_id": deviceGroupId})

	err := q.GetRelease(&deviceGroup)
	if err != nil {
		if err.Error() == rowNotFound {
			return nil, derrors.NewNotFoundError("device group credentials").WithParams(organizationId, deviceGroupId)
		} else {
			return nil, derrors.AsError(err, "cannot get device group credentials")
		}
	}

	return &deviceGroup, nil

}
func (sp * ScyllaDeviceCredentialsProvider) GetDeviceGroupByApiKey (apiKey string) (* entities.DeviceGroupCredentials, derrors.Error){

	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkConnectionAndConnect(); err != nil{
		return nil, err
	}

	var deviceGroup entities.DeviceGroupCredentials

	stmt, names := qb.Select(deviceGroupCredentialsTable).
		Where(qb.Eq("device_group_api_key")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"device_group_api_key": apiKey})

	err := q.GetRelease(&deviceGroup)
	if err != nil {
		if err.Error() == rowNotFound {
			return nil, derrors.NewNotFoundError("device group credentials apiKey").WithParams(apiKey)
		} else {
			return nil, derrors.AsError(err, "cannot get device group credentials")
		}
	}

	return &deviceGroup, nil

}
func (sp * ScyllaDeviceCredentialsProvider) RemoveDeviceGroup (organizationId string, deviceGroupId string) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkConnectionAndConnect(); err != nil{
		return err
	}

	// check if the group exists
	exists, err := sp.unsafeExistsGroupCredentials(organizationId, deviceGroupId)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("device group credentials").WithParams(organizationId, deviceGroupId)
	}

	stmt, _ := qb.Delete(deviceGroupCredentialsTable).
		Where(qb.Eq("organization_id")).Where(qb.Eq("device_group_id")).ToCql()
	cqlErr := sp.Session.Query(stmt, organizationId, deviceGroupId).Exec()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot delete device group credentials")
	}

	return nil
}
func (sp * ScyllaDeviceCredentialsProvider) TruncateDeviceGroup() derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.checkConnectionAndConnect(); err != nil {
		return err
	}

	// delete clusters table
	err := sp.Session.Query(fmt.Sprintf("TRUNCATE TABLE %s", deviceGroupCredentialsTable)).Exec()
	if err != nil {
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("failed to truncate the device group credentials table")
		return derrors.AsError(err, "cannot truncate device group table")
	}
	return nil
}

// ---------------------------
func (sp * ScyllaDeviceCredentialsProvider) unsafeExistsDeviceCredentials(organizationId string, deviceGroupId string, deviceId string) (bool, derrors.Error) {
	if err := sp.checkConnectionAndConnect(); err != nil{
		return false, err
	}

	var returnedId string

	stmt, names := qb.Select(deviceCredentialsTable).Columns("organization_id").
		Where(qb.Eq("organization_id")).Where(qb.Eq("device_group_id")).Where(qb.Eq("device_id")).
		ToCql()

	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationId,
		"device_group_id": deviceGroupId,
		"device_id": deviceId,
	})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		}else{
			return false, derrors.AsError(err, "cannot determinate if device credentials exists")
		}
	}

	return true, nil
}
func (sp * ScyllaDeviceCredentialsProvider) AddDeviceCredentials (credentials * entities.DeviceCredentials) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkConnectionAndConnect(); err != nil {
		return err
	}

	exists, err := sp.unsafeExistsGroupCredentials(credentials.OrganizationID, credentials.DeviceGroupID)
	if err != nil {
		return err
	}
	if !exists {
		return  derrors.NewNotFoundError("device group credentials").WithParams(credentials.OrganizationID, credentials.DeviceGroupID)
	}

	exists, err = sp.unsafeExistsDeviceCredentials(credentials.OrganizationID, credentials.DeviceGroupID, credentials.DeviceID)
	if err != nil {
		return err
	}
	if exists {
		return  derrors.NewAlreadyExistsError("device credentials").WithParams(credentials.OrganizationID, credentials.DeviceGroupID, credentials.DeviceID)
	}

	// add new basic credential
	stmt, names := qb.Insert(deviceCredentialsTable).Columns("organization_id", "device_group_id",
		"device_id", "device_api_key","enabled").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(credentials)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot add new device group credentials")
	}

	return nil
}
func (sp * ScyllaDeviceCredentialsProvider) UpdateDeviceCredentials (credentials * entities.DeviceCredentials) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkConnectionAndConnect(); err != nil {
		return err
	}

	exists, err := sp.unsafeExistsDeviceCredentials(credentials.OrganizationID, credentials.DeviceGroupID, credentials.DeviceID)
	if err != nil {
		return err
	}
	if ! exists {
		return  derrors.NewNotFoundError("device credentials").WithParams(credentials.OrganizationID, credentials.DeviceGroupID, credentials.DeviceID)
	}

	// add new basic credential
	stmt, names := qb.Update(deviceCredentialsTable).Set("enabled").
		Where(qb.Eq("organization_id")).Where(qb.Eq("device_group_id")).Where(qb.Eq("device_id")).
		ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(credentials)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot update device credentials")
	}

	return  nil
}
func (sp * ScyllaDeviceCredentialsProvider) ExistsDevice (organizationId string, deviceGroupId string, deviceId string) (bool, derrors.Error) {
	sp.Lock()
	defer sp.Unlock()

	return sp.unsafeExistsDeviceCredentials(organizationId, deviceGroupId, deviceId)
}
func (sp * ScyllaDeviceCredentialsProvider) GetDevice (organizationId string, deviceGroupId string, deviceId string) (* entities.DeviceCredentials, derrors.Error) {
	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkConnectionAndConnect(); err != nil{
		return nil, err
	}

	var device entities.DeviceCredentials

	stmt, names := qb.Select(deviceCredentialsTable).
		Where(qb.Eq("organization_id")).Where(qb.Eq("device_group_id")).Where(qb.Eq("device_id")).
		ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationId,
		"device_group_id": deviceGroupId,
		"device_id": deviceId,
	})

	err := q.GetRelease(&device)
	if err != nil {
		if err.Error() == rowNotFound {
			return nil, derrors.NewNotFoundError("device credentials").WithParams(organizationId, deviceGroupId, deviceId)
		} else {
			return nil, derrors.AsError(err, "cannot get device credentials")
		}
	}

	return &device, nil
}
func (sp * ScyllaDeviceCredentialsProvider) GetDeviceByApiKey (apiKey string) (* entities.DeviceCredentials, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkConnectionAndConnect(); err != nil{
		return nil, err
	}

	var device entities.DeviceCredentials

	stmt, names := qb.Select(deviceCredentialsTable).
		Where(qb.Eq("device_api_key")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"device_api_key": apiKey})

	err := q.GetRelease(&device)
	if err != nil {
		if err.Error() == rowNotFound {
			return nil, derrors.NewNotFoundError("device credentials apiKey").WithParams(apiKey)
		} else {
			return nil, derrors.AsError(err, "cannot get device credentials")
		}
	}

	return &device, nil
}
func (sp * ScyllaDeviceCredentialsProvider) RemoveDevice (organizationId string, deviceGroupId string, deviceId string) derrors.Error {
	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkConnectionAndConnect(); err != nil{
		return err
	}

	// check if the group exists
	exists, err := sp.unsafeExistsDeviceCredentials(organizationId, deviceGroupId, deviceId)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("device credentials").WithParams(organizationId, deviceGroupId, deviceId)
	}

	stmt, _ := qb.Delete(deviceCredentialsTable).
		Where(qb.Eq("organization_id")).Where(qb.Eq("device_group_id")).Where(qb.Eq("device_id")).ToCql()
	cqlErr := sp.Session.Query(stmt, organizationId, deviceGroupId, deviceId).Exec()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot delete device credentials")
	}

	return nil
}
func (sp * ScyllaDeviceCredentialsProvider) TruncateDevice() derrors.Error {
	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.checkConnectionAndConnect(); err != nil {
		return err
	}

	// delete clusters table
	err := sp.Session.Query(fmt.Sprintf("TRUNCATE TABLE %s", deviceCredentialsTable)).Exec()
	if err != nil {
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("failed to truncate the device credentials table")
		return derrors.AsError(err, "cannot truncate device table")
	}
	return nil
}

func (sp * ScyllaDeviceCredentialsProvider) Truncate()  derrors.Error{

	err := sp.TruncateDeviceGroup()
	if err != nil {
		return err
	}

	err = sp.TruncateDevice()
	if err != nil {
		return err
	}

	return nil
}

