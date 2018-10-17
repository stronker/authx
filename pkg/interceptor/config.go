/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package interceptor

type AuthorizationConfig struct {
	Permissions map[string]Permission `json:"permissions"`
}

type Config struct {
	Authorization *AuthorizationConfig
	Secret        string
	Header        string
	AllowsAll     bool
}

func NewConfig(config *AuthorizationConfig,
	secret string, header string, allowsAll bool) *Config {

	return &Config{Authorization: config, Secret: secret, Header: header, AllowsAll: allowsAll}
}
