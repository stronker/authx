/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package credentials

import (
	"github.com/nalej/authx/internal/app/authx/entities"
	"github.com/nalej/derrors"
)


// BasicCredentials is the interface that define how we are store the basic credential information.
type BasicCredentials interface {
	// Delete remove a specific user credentials.
	Delete(username string) derrors.Error
	// Add adds a new basic credentials.
	Add(credentials *entities.BasicCredentialsData) derrors.Error
	// Get recover a user credentials.
	Get(username string) (*entities.BasicCredentialsData, derrors.Error)
	// Edit update a specific user credentials.
	Edit(username string, edit *entities.EditBasicCredentialsData) derrors.Error
	// Exist check if exists a specific credentials.
	Exist(username string) (*bool,derrors.Error)
	// Truncate removes all credentials.
	Truncate() derrors.Error
}
