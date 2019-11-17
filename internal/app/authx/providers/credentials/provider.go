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

package credentials

import (
	"github.com/nalej/derrors"
	"github.com/stronker/authx/internal/app/authx/entities"
)

// BasicCredentials is the interface that define how we are store the basic credential information.
type BasicCredentials interface {
	// Delete remove a specific user credentials.
	Delete(username string) derrors.Error
	// Add adds a new basic credentials.
	Add(credentials *entities.BasicCredentialsData) derrors.Error
	// Get recover a user credentials.
	Get(username string) (*entities.BasicCredentialsData, derrors.Error)
	// Edit update a specific user credentials.
	Edit(username string, edit *entities.EditBasicCredentialsData) derrors.Error
	// Exist check if exists a specific credentials.
	Exist(username string) (*bool, derrors.Error)
	// Truncate removes all credentials.
	Truncate() derrors.Error
}
