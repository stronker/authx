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

package interceptor

import (
	"encoding/json"
	"github.com/nalej/derrors"
	"github.com/rs/zerolog/log"
	"io/ioutil"
)

const DefaultCacheEntries = 100

// AuthorizationConfig is structure that contains a set of permissions. The key of the map is the method name.
type AuthorizationConfig struct {
	// AllowsAll If the header is not found, allow access depending on this parameter.
	AllowsAll bool `json:"allows_all"`
	// Permission is a map of permissions the key is the method name.
	Permissions map[string]Permission `json:"permissions"`
}

func LoadAuthorizationConfig(path string) (*AuthorizationConfig, derrors.Error) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, derrors.NewInvalidArgumentError("impossible read config file", err)
	}

	authCfg := &AuthorizationConfig{}
	jErr := json.Unmarshal(dat, authCfg)
	if jErr != nil {
		return nil, derrors.NewInternalError("impossible unmarshal file", jErr)
	}
	log.Debug().Int("permissions", len(authCfg.Permissions)).Msg("Authorization matrix loaded")
	return authCfg, nil
}

// Config is the complete configuration file.
type Config struct {
	Authorization *AuthorizationConfig
	// Secret contains the shared secret with the authx component to sign the JWT token.
	Secret string
	// Name of the header where the token is found.
	Header string
	// Number of cached entries for group secrets
	NumCacheEntries int
}

// NewConfig creates a new instance of the structure.
func NewConfig(config *AuthorizationConfig,
	secret string, header string) *Config {

	return &Config{Authorization: config, Secret: secret, Header: header, NumCacheEntries: DefaultCacheEntries}
}
