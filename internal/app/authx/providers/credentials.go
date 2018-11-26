/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package providers

import (
	"github.com/nalej/derrors"
	"sync"
)
// BasicCredentialsData is the struct that is store in the database.
type BasicCredentialsData struct {
	// Username is the credential id.
	Username       string
	// Password is the user defined password.
	Password       [] byte
	// RoleID is the assigned role.
	RoleID         string
	// OrganizationID is the assigned organization.
	OrganizationID string
}

// NewBasicCredentialsData creates an instance of BasicCredentialsData.
func NewBasicCredentialsData(username string, password [] byte, roleID string, organizationID string) *BasicCredentialsData {
	return &BasicCredentialsData{
		Username:       username,
		Password:       password,
		RoleID:         roleID,
		OrganizationID: organizationID,
	}
}

// EditBasicCredentialsData is an object that allows to edit the credetentials record.
type EditBasicCredentialsData struct {
	Password *[] byte
	RoleID   *string
}

// WithPassword allows to change the password.
func (d *EditBasicCredentialsData) WithPassword(password [] byte) *EditBasicCredentialsData {
	d.Password = &password
	return d
}

// WithRoleID allows to change the roleID
func (d *EditBasicCredentialsData) WithRoleID(roleID string) *EditBasicCredentialsData {
	d.RoleID = &roleID
	return d
}

// NewEditBasicCredentialsData create a new instance of EditBasicCredentialsData.
func NewEditBasicCredentialsData() *EditBasicCredentialsData {
	return &EditBasicCredentialsData{}
}

// BasicCredentials is the interface that define how we are store the basic credential information.
type BasicCredentials interface {
	// Delete remove a specific user credentials.
	Delete(username string) derrors.Error
	// Add adds a new basic credentials.
	Add(credentials *BasicCredentialsData) derrors.Error
	// Get recover a user credentials.
	Get(username string) (*BasicCredentialsData, derrors.Error)
	// Edit update a specific user credentials.
	Edit(username string, edit *EditBasicCredentialsData) derrors.Error
	// Exist check if exists a specific credentials.
	Exist(username string) (*bool,derrors.Error)
	// Truncate removes all credentials.
	Truncate() derrors.Error
}

// BasicCredentialsMockup is an implementation of this provider only for testing
type BasicCredentialsMockup struct {
	sync.Mutex
	data map[string]BasicCredentialsData
}

// NewBasicCredentialMockup create new mockup.
func NewBasicCredentialMockup() *BasicCredentialsMockup {
	return &BasicCredentialsMockup{data: map[string]BasicCredentialsData{}}
}

// Delete remove a specific user credentials.
func (p *BasicCredentialsMockup) Delete(username string) derrors.Error {
	//p.Lock()
	//defer p.Unlock()
	_, ok := p.data[username]
	if !ok {
		return derrors.NewNotFoundError("username not found").WithParams(username)
	}
	delete(p.data, username)
	return nil
}

// Add adds a new basic credentials.
func (p *BasicCredentialsMockup) Add(credentials *BasicCredentialsData) derrors.Error {
	p.Lock()
	defer p.Unlock()
	p.data[credentials.Username] = *credentials
	return nil
}

// Get recover a user credentials.
func (p *BasicCredentialsMockup) Get(username string) (*BasicCredentialsData, derrors.Error) {
	p.Lock()
	defer p.Unlock()
	data, ok := p.data[username]
	if !ok {
		return nil, derrors.NewNotFoundError("credentials not found").WithParams(username)
	}
	return &data, nil
}

// Exist check if exists a specific credentials.
func (p *BasicCredentialsMockup) Exist(username string) (*bool,derrors.Error){
	p.Lock()
	defer p.Unlock()
	_, ok := p.data[username]
	return &ok,nil
}

// Edit update a specific user credentials.
func (p *BasicCredentialsMockup) Edit(username string, edit *EditBasicCredentialsData) derrors.Error {

	data, err := p.Get(username)
	if err != nil {
		return err
	}
	if edit.RoleID != nil {
		data.RoleID = *edit.RoleID
	}
	if edit.Password != nil {
		data.Password = *edit.Password
	}
	p.data[username] = *data
	return nil
}

// Truncate removes all credentials.
func (p *BasicCredentialsMockup) Truncate() derrors.Error {
	p.Lock()
	p.data = map[string]BasicCredentialsData{}
	defer p.Unlock()
	return nil
}
