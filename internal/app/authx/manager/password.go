/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package manager

import (
	"github.com/nalej/derrors"
	"golang.org/x/crypto/bcrypt"
)

//Password is an interface to generate hash of passwords.
type Password interface {
	// GenerateHashedPassword generates a password with a random salt.
	GenerateHashedPassword(password string) ([] byte, derrors.Error)
	// CompareHashAndPassword compare a hashed password with a specif password.
	CompareHashAndPassword(hashedPassword [] byte, password string) derrors.Error
}

// NewBCryptPassword build a object that uses BCrypt to implement the Password interface.
func NewBCryptPassword() Password {
	return &BCryptPassword{cost: bcrypt.DefaultCost}
}

// BCryptPassword implementation of Password using BCrypt
type BCryptPassword struct {
	cost int
}

// GenerateHashedPassword generates a password with a random salt.
func (m *BCryptPassword) GenerateHashedPassword(password string) ([] byte, derrors.Error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return hashedPassword, derrors.AsError(err, "error hashing the password")
}

// CompareHashAndPassword compare a hashed password with a specif password.
func (m *BCryptPassword) CompareHashAndPassword(hashedPassword [] byte, password string) derrors.Error {
	return derrors.AsError(bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)), "error comparing passwords")
}


