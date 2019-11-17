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

package role

import (
	"github.com/gocql/gocql"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"github.com/stronker/authx/internal/app/authx/entities"
	"sync"
)

const table = "roles"
const tablePK_1 = "organization_id"
const tablePK_2 = "role_id"

const rowNotFound = "not found"

type ScyllaRoleProvider struct {
	Address  string
	Port     int
	KeySpace string
	sync.Mutex
	Session *gocql.Session
}

func NewScyllaRoleProvider(address string, port int, keyspace string) *ScyllaRoleProvider {
	provider := ScyllaRoleProvider{Address: address, Port: port, KeySpace: keyspace}
	provider.connect()
	return &provider
}

func (sp *ScyllaRoleProvider) connect() derrors.Error {
	
	// connect to the cluster
	conf := gocql.NewCluster(sp.Address)
	conf.Keyspace = sp.KeySpace
	conf.Port = sp.Port
	
	session, err := conf.CreateSession()
	if err != nil {
		log.Error().Str("provider", "ScyllaRolesProvider").Str("trace", conversions.ToDerror(err).DebugReport()).Msg("unable to connect")
		return derrors.AsError(err, "cannot connect")
	}
	
	sp.Session = session
	return nil
}

func (sp *ScyllaRoleProvider) Disconnect() {
	
	sp.Lock()
	defer sp.Unlock()
	
	if sp.Session != nil {
		sp.Session.Close()
		sp.Session = nil
	}
	
}

func (sp *ScyllaRoleProvider) checkConnectionAndConnect() derrors.Error {
	
	if sp.Session != nil {
		return nil
	}
	log.Info().Str("provider", "ScyllaRolesProvider").Msg("session not connected, trying to connect it!")
	err := sp.connect()
	if err != nil {
		return err
	}
	
	return nil
}

// --------------------------------------------------------------------------------------------------------------------

func (sp *ScyllaRoleProvider) unsafeGet(organizationID string, roleID string) (*entities.RoleData, derrors.Error) {
	
	var role entities.RoleData
	stmt, names := qb.Select(table).Where(qb.Eq(tablePK_1)).Where(qb.Eq(tablePK_2)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		tablePK_1: organizationID,
		tablePK_2: roleID})
	
	err := q.GetRelease(&role)
	if err != nil {
		if err.Error() == rowNotFound {
			return nil, derrors.NewNotFoundError("role").WithParams(organizationID, roleID)
		} else {
			return nil, derrors.AsError(err, "cannot get role")
		}
	}
	
	return &role, nil
}

func (sp *ScyllaRoleProvider) unsafeExist(organizationID string, roleID string) (*bool, derrors.Error) {
	
	ok := false
	var returnedId string
	
	stmt, names := qb.Select(table).Columns(tablePK_1).Where(qb.Eq(tablePK_1)).Where(qb.Eq(tablePK_2)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		tablePK_1: organizationID,
		tablePK_2: roleID})
	
	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound {
			return &ok, nil
		} else {
			return &ok, derrors.AsError(err, "cannot determine if role exists")
		}
	}
	ok = true
	return &ok, nil
}

// --------------------------------------------------------------------------------------------------------------------

func (sp *ScyllaRoleProvider) Delete(organizationID string, roleID string) derrors.Error {
	
	sp.Lock()
	defer sp.Unlock()
	
	// check connection
	if err := sp.checkConnectionAndConnect(); err != nil {
		return err
	}
	
	exists, err := sp.unsafeExist(organizationID, roleID)
	
	if err != nil {
		return err
	}
	if !*exists {
		return derrors.NewNotFoundError("role").WithParams(organizationID, roleID)
	}
	
	stmt, _ := qb.Delete(table).Where(qb.Eq(tablePK_1)).Where(qb.Eq(tablePK_2)).ToCql()
	cqlErr := sp.Session.Query(stmt, organizationID, roleID).Exec()
	
	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot delete role")
	}
	
	return nil
}

// Add a new role.
func (sp *ScyllaRoleProvider) Add(role *entities.RoleData) derrors.Error {
	
	sp.Lock()
	defer sp.Unlock()
	
	if err := sp.checkConnectionAndConnect(); err != nil {
		return err
	}
	
	exists, err := sp.unsafeExist(role.OrganizationID, role.RoleID)
	if err != nil {
		return err
	}
	if *exists {
		return derrors.NewAlreadyExistsError("role").WithParams(role.OrganizationID, role.RoleID)
	}
	
	// add new basic credential
	stmt, names := qb.Insert(table).Columns("organization_id", "role_id", "name", "internal", "primitives").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(role)
	cqlErr := q.ExecRelease()
	
	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot add new role")
	}
	
	return nil
}

// Get recovers an existing role.
func (sp *ScyllaRoleProvider) Get(organizationID string, roleID string) (*entities.RoleData, derrors.Error) {
	
	sp.Lock()
	defer sp.Unlock()
	
	if err := sp.checkConnectionAndConnect(); err != nil {
		return nil, err
	}
	
	var role entities.RoleData
	stmt, names := qb.Select(table).Where(qb.Eq(tablePK_1)).Where(qb.Eq(tablePK_2)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		tablePK_1: organizationID,
		tablePK_2: roleID})
	
	err := q.GetRelease(&role)
	if err != nil {
		if err.Error() == rowNotFound {
			return nil, derrors.NewNotFoundError("role").WithParams(organizationID, roleID)
		} else {
			return nil, derrors.AsError(err, "cannot get role")
		}
	}
	
	return &role, nil
}

// Edit updates an existing role.
func (sp *ScyllaRoleProvider) Edit(organizationID string, roleID string, edit *entities.EditRoleData) derrors.Error {
	
	sp.Lock()
	defer sp.Unlock()
	
	data, err := sp.unsafeGet(organizationID, roleID)
	
	if err != nil {
		return err
	}
	if edit.Name != nil {
		data.Name = *edit.Name
	}
	if edit.Primitives != nil {
		data.Primitives = *edit.Primitives
	}
	// update
	stmt, names := qb.Update(table).Set("name", "primitives").Where(qb.Eq(tablePK_1)).Where(qb.Eq(tablePK_2)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(data)
	cqlErr := q.ExecRelease()
	
	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot edit role")
	}
	
	return nil
}

// Exist checks if a role exists.
func (sp *ScyllaRoleProvider) Exist(organizationID string, roleID string) (*bool, derrors.Error) {
	
	sp.Lock()
	defer sp.Unlock()
	
	ok := false
	
	if err := sp.checkConnectionAndConnect(); err != nil {
		return &ok, err
	}
	
	var returnedId string
	
	stmt, names := qb.Select(table).Columns(tablePK_1).Where(qb.Eq(tablePK_1)).Where(qb.Eq(tablePK_2)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		tablePK_1: organizationID,
		tablePK_2: roleID})
	
	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound {
			return &ok, nil
		} else {
			return &ok, derrors.AsError(err, "cannot determine if role exists")
		}
	}
	ok = true
	return &ok, nil
}

// List the roles associated with an organization.
func (sp *ScyllaRoleProvider) List(organizationID string) ([]entities.RoleData, derrors.Error) {
	
	sp.Lock()
	defer sp.Unlock()
	
	if err := sp.checkConnectionAndConnect(); err != nil {
		return nil, err
	}
	
	result := make([]entities.RoleData, 0)
	
	stmt, names := qb.Select(table).Where(qb.Eq(tablePK_1)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		tablePK_1: organizationID,
	})
	
	cqlErr := gocqlx.Select(&result, q.Query)
	
	if cqlErr != nil {
		return nil, derrors.AsError(cqlErr, "cannot list roles")
	}
	
	return result, nil
}

// Truncate clears the provider.
func (sp *ScyllaRoleProvider) Truncate() derrors.Error {
	sp.Lock()
	defer sp.Unlock()
	
	if err := sp.checkConnectionAndConnect(); err != nil {
		return err
	}
	
	err := sp.Session.Query("TRUNCATE TABLE roles").Exec()
	if err != nil {
		dErr := derrors.AsError(err, "cannot truncate role table")
		log.Error().Str("trace", dErr.DebugReport()).Msg("failed to truncate the table")
		return dErr
	}
	
	return nil
}
