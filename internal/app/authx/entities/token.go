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

// TokenData is the information that the system stores.
type TokenData struct {
	Username       string `cql:"user_name"`
	TokenID        string `cql:"token_id"`
	RefreshToken   []byte `cql:"refresh_token"`
	ExpirationDate int64  `cql:"expiration_date"`
}

// NewTokenData creates an instance of the structure
func NewTokenData(username string, tokenID string, refreshToken []byte,
	expirationDate int64) *TokenData {

	return &TokenData{Username: username,
		TokenID:        tokenID,
		RefreshToken:   refreshToken,
		ExpirationDate: expirationDate}
}
