/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package manager

import (
	"github.com/nalej/authx/internal/app/authx/providers"
	"github.com/nalej/derrors"
)

type Authx struct {
	Password Password
	Provider providers.BasicCredentials
}

func (m *Authx) DeleteCredentials(username string) derrors.Error {
	return m.Provider.Delete(username)
}

func (m *Authx) AddBasicCredentials(username string, organizationID string, roleID string, password string) derrors.Error {
	hashedPassword, err := m.Password.GenerateHashedPassword(password)
	if err != nil {
		return err
	}

	entity := providers.NewBasicCredentialsData(username, hashedPassword, roleID, organizationID)
	return m.Provider.Add(entity)
}

func (m *Authx) LoginWithBasicCredentials(username string, password string) derrors.Error {
	credentials, err := m.Provider.Get(username)
	if err != nil {
		return err
	}
	return m.Password.CompareHashAndPassword(credentials.Password, password)
}
