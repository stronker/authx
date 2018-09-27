/*
 * Copyright 2018 Nalej
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


