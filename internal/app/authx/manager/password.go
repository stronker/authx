/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package manager

import (
	"github.com/nalej/derrors"
	"golang.org/x/crypto/bcrypt"
)

//Password is an interface to generate hash of passwords.
type Password interface {
	// GenerateHashedPassword generates a password with a random salt.
	GenerateHashedPassword(password string) ([]byte, derrors.Error)
	// CompareHashAndPassword compare a hashed password with a specif password.
	CompareHashAndPassword(hashedPassword []byte, password string) derrors.Error
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
func (m *BCryptPassword) GenerateHashedPassword(password string) ([]byte, derrors.Error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, derrors.NewInternalError("error hashing the password", err)
	}
	return hashedPassword, nil
}

// CompareHashAndPassword compare a hashed password with a specif password.
func (m *BCryptPassword) CompareHashAndPassword(hashedPassword []byte, password string) derrors.Error {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return derrors.NewUnauthenticatedError("password is not valid", err)
	}
	return nil
}
