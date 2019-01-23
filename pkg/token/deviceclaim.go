/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package token

import "github.com/dgrijalva/jwt-go"

// Claim joins the personal claim and the standard JWT claim.
type DeviceClaim struct {
	jwt.StandardClaims
	OrganizationID string    `json:"organizationID,omitempty"`
	DeviceID string `json:"device_id,omitempty"`
	DeviceGroupID string `json:"device_group_id,omitempty"`
	Primitives     [] string `json:"access,omitempty"`
}
