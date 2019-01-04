/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package token

import (
	"github.com/nalej/authx/internal/app/authx/entities"
	"github.com/nalej/derrors"
)


// Token is the interface to store the token information.
type Token interface {
	// Delete an existing token.
	Delete(username string, tokenID string) derrors.Error
	// Add a token.
	Add(token *entities.TokenData) derrors.Error
	// Get an existing token.
	Get(username string, tokenID string) (*entities.TokenData, derrors.Error)
	// Exist checks if the token was added.
	Exist(username string, tokenID string) (*bool, derrors.Error)
	// Truncate cleans all data.
	Truncate() derrors.Error

	DeleteExpiredTokens() derrors.Error
}

