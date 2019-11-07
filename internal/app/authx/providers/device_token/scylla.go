/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package device_token

import (
	"github.com/gocql/gocql"
	"github.com/nalej/authx/internal/app/authx/entities"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"sync"
	"time"
)

const rowNotFound = "not found"
const table = "devicetokens"

const ttlExpired = time.Duration(3) * time.Hour

type ScyllaDeviceTokenProvider struct {
	Address  string
	Port     int
	KeySpace string
	sync.Mutex
	Session *gocql.Session
}

func NewScyllaDeviceTokenProvider(address string, port int, keyspace string) *ScyllaDeviceTokenProvider {
	provider := ScyllaDeviceTokenProvider{Address: address, Port: port, KeySpace: keyspace}
	provider.connect()
	return &provider
}

func (sp *ScyllaDeviceTokenProvider) connect() derrors.Error {

	// connect to the cluster
	conf := gocql.NewCluster(sp.Address)
	conf.Keyspace = sp.KeySpace
	conf.Port = sp.Port

	session, err := conf.CreateSession()
	if err != nil {
		log.Error().Str("provider", "ScyllaDeviceTokenProvider").Str("trace", conversions.ToDerror(err).DebugReport()).Msg("unable to connect")
		return derrors.AsError(err, "cannot connect")
	}

	sp.Session = session
	return nil
}

func (sp *ScyllaDeviceTokenProvider) Disconnect() {

	sp.Lock()
	defer sp.Unlock()

	if sp.Session != nil {
		sp.Session.Close()
		sp.Session = nil
	}

}

// -------------------- //
// -- unsafe methods -- //
// -------------------- //
func (sp *ScyllaDeviceTokenProvider) unsafeGet(deviceID string, tokenID string) (*entities.TokenData, derrors.Error) {

	if err := sp.checkConnectionAndConnect(); err != nil {
		return nil, err
	}

	var token entities.TokenData
	stmt, names := qb.Select(table).Where(qb.Eq("device_id")).Where(qb.Eq("token_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"device_id": deviceID,
		"token_id":  tokenID})

	err := q.GetRelease(&token)
	if err != nil {
		if err.Error() == rowNotFound {
			return nil, derrors.NewNotFoundError("device token").WithParams(deviceID, tokenID)
		} else {
			return nil, derrors.AsError(err, "cannot get device token")
		}
	}

	return &token, nil
}

// Exist checks if the token was added.
func (sp *ScyllaDeviceTokenProvider) unsafeExist(deviceID string, tokenID string) (*bool, derrors.Error) {

	ok := false

	if err := sp.checkConnectionAndConnect(); err != nil {
		return &ok, err
	}

	var count int

	stmt, names := qb.Select(table).Count("device_id").Where(qb.Eq("device_id")).Where(qb.Eq("token_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"device_id": deviceID,
		"token_id":  tokenID})

	err := q.GetRelease(&count)
	if err != nil {
		if err.Error() == rowNotFound { // TODO: mirar si lo puedo quitar
			return &ok, nil
		} else {
			return &ok, derrors.AsError(err, "cannot determinate if device token exists")
		}
	}
	if count > 0 {
		ok = true
	}
	return &ok, nil
}

func (sp *ScyllaDeviceTokenProvider) checkConnectionAndConnect() derrors.Error {

	if sp.Session != nil {
		return nil
	}
	log.Info().Str("provider", "ScyllaDeviceTokeProvider").Msg("session not connected, trying to connect it!")
	err := sp.connect()
	if err != nil {
		return err
	}

	return nil
}

// ----------------------- //
// -- interface methods -- //
// ----------------------- //
// Delete an existing token.
func (sp *ScyllaDeviceTokenProvider) Delete(deviceID string, tokenID string) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkConnectionAndConnect(); err != nil {
		return err
	}

	exists, err := sp.unsafeExist(deviceID, tokenID)

	if err != nil {
		return err
	}
	if !*exists {
		return derrors.NewNotFoundError("device token").WithParams(deviceID, tokenID)
	}

	stmt, _ := qb.Delete(table).Where(qb.Eq("device_id")).Where(qb.Eq("token_id")).ToCql()
	cqlErr := sp.Session.Query(stmt, deviceID, tokenID).Exec()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot delete device token")
	}

	return nil
}

// Add a token.
func (sp *ScyllaDeviceTokenProvider) Add(token *entities.DeviceTokenData) derrors.Error {
	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkConnectionAndConnect(); err != nil {
		return err
	}

	exists, err := sp.unsafeExist(token.DeviceId, token.TokenID)
	if err != nil {
		return err
	}
	if *exists {
		return derrors.NewAlreadyExistsError("device token").WithParams(token.DeviceId, token.TokenID)
	}

	// add new basic credential
	stmt, names := qb.Insert(table).Columns("device_id", "token_id", "refresh_token", "expiration_date", "organization_id", "device_group_id").TTL(ttlExpired).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(token)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot add device token")
	}

	return nil
}

// Get an existing token.
func (sp *ScyllaDeviceTokenProvider) Get(deviceID string, tokenID string) (*entities.DeviceTokenData, derrors.Error) {
	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkConnectionAndConnect(); err != nil {
		return nil, err
	}

	var token entities.DeviceTokenData
	stmt, names := qb.Select(table).Where(qb.Eq("device_id")).Where(qb.Eq("token_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"device_id": deviceID,
		"token_id":  tokenID})

	err := q.GetRelease(&token)
	if err != nil {
		if err.Error() == rowNotFound {
			return nil, derrors.NewNotFoundError("device token").WithParams(deviceID, tokenID)
		} else {
			return nil, derrors.AsError(err, "cannot get device token")
		}
	}

	return &token, nil
}

// Exist checks if the token was added.
func (sp *ScyllaDeviceTokenProvider) Exist(deviceID string, tokenID string) (*bool, derrors.Error) {
	sp.Lock()
	defer sp.Unlock()

	return sp.unsafeExist(deviceID, tokenID)
}

// Update an existing token
func (sp *ScyllaDeviceTokenProvider) Update(token *entities.DeviceTokenData) derrors.Error {
	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkConnectionAndConnect(); err != nil {
		return err
	}

	exists, err := sp.unsafeExist(token.DeviceId, token.TokenID)
	if err != nil {
		return err
	}
	if !*exists {
		return derrors.NewNotFoundError("device token").WithParams(token.DeviceId, token.TokenID)
	}

	// add new basic credential
	stmt, names := qb.Update(table).Set("expiration_date", "refresh_token").
		Where(qb.Eq("device_id")).Where(qb.Eq("token_id")).TTL(ttlExpired).
		ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(token)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot update device token")
	}

	return nil
}

// Truncate cleans all data.
func (sp *ScyllaDeviceTokenProvider) Truncate() derrors.Error {
	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkConnectionAndConnect(); err != nil {
		return err
	}

	err := sp.Session.Query("TRUNCATE TABLE deviceTokens").Exec()
	if err != nil {
		log.Info().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("failed to truncate the device Tokens table")
		return derrors.AsError(err, "cannot truncate device token table")
	}

	return nil
}

func (sp *ScyllaDeviceTokenProvider) DeleteExpiredTokens() derrors.Error {
	// nothing to do, ttl used to delete expired tokens
	return nil
}

func (m *ScyllaDeviceTokenProvider) GetByRefreshToken(refreshToken string) (*entities.DeviceTokenData, derrors.Error) {
	return nil, nil
}
