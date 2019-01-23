/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package manager

import (
	"github.com/nalej/authx/pkg/token"
	"github.com/nalej/derrors"
	"time"
)

// Token is a interface manages the business logic of tokens.
type DeviceToken interface {
	// Generate a new token with the device claim.
	Generate(deviceClaim *token.DeviceClaim, expirationPeriod time.Duration, secret string) (*GeneratedToken, derrors.Error)
	// Refresh renew an old token.
	Refresh(oldToken string, refreshToken string,
		expirationPeriod time.Duration, secret string) (*GeneratedToken, derrors.Error)
	// Clean remove all the data from the providers.
	Clean() derrors.Error
}



