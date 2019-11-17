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

package token

import (
	"github.com/nalej/derrors"
	"github.com/stronker/authx/internal/app/authx/entities"
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
	// Update an existing token
	Update(token *entities.TokenData) derrors.Error
	// Truncate cleans all data.
	Truncate() derrors.Error
	
	DeleteExpiredTokens() derrors.Error
}
