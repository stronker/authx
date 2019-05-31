/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package token

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"time"
)

// PersonalClaim is the claim that include system information.
type PersonalClaim struct {
	UserID         string   `json:"userID,omitempty"`
	Primitives     []string `json:"access,omitempty"`
	RoleName       string   `json:"role,omitempty"`
	OrganizationID string   `json:"organizationID,omitempty"`
}

// NewPersonalClaim creates a new instance of the structure.
func NewPersonalClaim(userID string, roleName string, primitives []string, organizationID string) *PersonalClaim {
	return &PersonalClaim{UserID: userID, RoleName: roleName, Primitives: primitives, OrganizationID: organizationID}
}

// Claim joins the personal claim and the standard JWT claim.
type Claim struct {
	jwt.StandardClaims
	PersonalClaim
}

// NewClaim create a new instance of the structure.
func NewClaim(personalClaim PersonalClaim, issuer string, creationTime time.Time, expirationPeriod time.Duration) *Claim {
	stdClaim := jwt.StandardClaims{
		Issuer:    issuer,
		Id:        GenerateUUID(),
		ExpiresAt: creationTime.Add(expirationPeriod).Unix(),
		NotBefore: creationTime.Unix(),
		IssuedAt:  creationTime.Unix(),
	}

	return &Claim{StandardClaims: stdClaim, PersonalClaim: personalClaim}
}

// GenerateUUID creates a new random UUID.
func GenerateUUID() string {
	return uuid.New().String()
}
