/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package inventory

import (
	"github.com/nalej/authx/internal/app/authx/entities"
	"github.com/nalej/derrors"
)

type Provider interface{
	// AddECJoinToken stores a join token for edge controllers.
	AddECJoinToken(token *entities.EICJoinToken) derrors.Error
	// IsJoinTokenValidForEC checks if a token is still valid for joining new EC
	GetECJoinToken(organizationID string, token string) (*entities.EICJoinToken, derrors.Error)
	// Clear all elements
	Clear() derrors.Error
}
