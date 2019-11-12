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

package entities

import (
	"github.com/google/uuid"
	"time"
)

type EICJoinToken struct {
	OrganizationID string `json:"organization_id"`
	TokenID        string `json:"token_id"`
	ExpiresOn      int64  `json:"expires_on"`
}

func NewEICJoinToken(organizationID string, ttl time.Duration) *EICJoinToken {
	return &EICJoinToken{
		OrganizationID: organizationID,
		TokenID:        uuid.New().String(),
		ExpiresOn:      time.Now().Add(ttl).Unix(),
	}
}
