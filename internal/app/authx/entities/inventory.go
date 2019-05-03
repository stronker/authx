/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package entities

import (
	"github.com/google/uuid"
	"time"
)

type EICJoinToken struct{
	OrganizationID string `json:"organization_id"`
	TokenID string `json:"token_id"`
	ExpiresOn int64 `json:"expires_on"`
}

func NewEICJoinToken(organizationID string, ttl time.Duration) * EICJoinToken{
	return &EICJoinToken{
		OrganizationID: organizationID,
		TokenID:        uuid.New().String(),
		ExpiresOn:      time.Now().Add(ttl).Unix(),
	}
}