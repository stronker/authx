/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package providers

import "github.com/nalej/derrors"

type BasicCredentialsData struct {
	Username       string
	Password       [] byte
	RoleID         string
	OrganizationID string
}

func NewBasicCredentialsData(username string, password [] byte, roleID string, organizationID string) *BasicCredentialsData {
	return &BasicCredentialsData{
		Username:       username,
		Password:       password,
		RoleID:         roleID,
		OrganizationID: organizationID,
	}
}

type BasicCredentials interface {
	Delete(username string) derrors.Error
	Add(credentials *BasicCredentialsData) derrors.Error
	Get(username string) (*BasicCredentialsData, derrors.Error)
}

type BasicCredentialsMockup struct {
	data map[string]BasicCredentialsData
}

func NewBasicCredentialMockup() *BasicCredentialsMockup {
	return &BasicCredentialsMockup{data: map[string]BasicCredentialsData{}}
}

func (p *BasicCredentialsMockup) Delete(username string) derrors.Error {
	_, ok := p.data[username]
	if !ok {
		return derrors.NewOperationError("Not found username")
	}
	delete(p.data, username)
	return nil
}

func (p *BasicCredentialsMockup) Add(credentials *BasicCredentialsData) derrors.Error {
	p.data[credentials.Username] = *credentials
	return nil
}

func (p *BasicCredentialsMockup) Get(username string) (*BasicCredentialsData, derrors.Error) {
	data := p.data[username]
	return &data, nil
}
