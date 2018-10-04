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

type EditBasicCredentialsData struct {
	Password *[] byte
	RoleID   *string
}

func (d *EditBasicCredentialsData) WithPassword(password [] byte) *EditBasicCredentialsData {
	d.Password = &password
	return d
}

func (d *EditBasicCredentialsData) WithRoleID(roleID string) *EditBasicCredentialsData {
	d.RoleID = &roleID
	return d
}

func NewEditBasicCredentialsData() *EditBasicCredentialsData {
	return &EditBasicCredentialsData{}
}

type BasicCredentials interface {
	Delete(username string) derrors.Error
	Add(credentials *BasicCredentialsData) derrors.Error
	Get(username string) (*BasicCredentialsData, derrors.Error)
	Edit(username string, edit *EditBasicCredentialsData) derrors.Error
	Truncate() derrors.Error
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
		return derrors.NewOperationError("username not found")
	}
	delete(p.data, username)
	return nil
}

func (p *BasicCredentialsMockup) Add(credentials *BasicCredentialsData) derrors.Error {
	p.data[credentials.Username] = *credentials
	return nil
}

func (p *BasicCredentialsMockup) Get(username string) (*BasicCredentialsData, derrors.Error) {
	data, ok := p.data[username]
	if !ok {
		return nil, nil
	}
	return &data, nil
}

func (p *BasicCredentialsMockup) Edit(username string, edit *EditBasicCredentialsData) derrors.Error {
	data, ok := p.data[username]
	if !ok {
		return derrors.NewOperationError("username not found")
	}
	if edit.RoleID != nil {
		data.RoleID = *edit.RoleID
	}
	if edit.Password != nil {
		data.Password = *edit.Password
	}
	p.data[username] = data
	return nil
}

func (p *BasicCredentialsMockup) Truncate() derrors.Error {
	p.data = map[string]BasicCredentialsData{}
	return nil
}
