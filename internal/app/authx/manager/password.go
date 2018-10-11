/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package manager

import (
	"github.com/nalej/derrors"
	"golang.org/x/crypto/bcrypt"
)

type Password interface {
	GenerateHashedPassword(password string) ([] byte, derrors.Error)
	CompareHashAndPassword(hashedPassword [] byte, password string) derrors.Error
}

func NewBCryptPassword() Password {
	return &BCryptPassword{cost: bcrypt.DefaultCost}
}


type BCryptPassword struct {
	cost int
}

func (m *BCryptPassword) GenerateHashedPassword(password string) ([] byte, derrors.Error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return hashedPassword, derrors.AsError(err, "error hashing the password")
}

func (m *BCryptPassword) CompareHashAndPassword(hashedPassword [] byte, password string) derrors.Error {
	return derrors.AsError(bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)), "error comparing passwords")
}


