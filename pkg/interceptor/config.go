/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package interceptor

// AuthorizationConfig is structure that contains a set of permissions. The key of the map is the method name.
type AuthorizationConfig struct {
	AllowsAll bool `json:"allowsAll"`
	// Permission is a map of permissions the key is the method name.
	Permissions map[string]Permission `json:"permissions"`
}

// Config is the complete configuration file.
type Config struct {
	Authorization *AuthorizationConfig
	Secret        string
	Header        string
}

// NewConfig creates a new instance of the structure.
func NewConfig(config *AuthorizationConfig,
	secret string, header string) *Config {

	return &Config{Authorization: config, Secret: secret, Header: header}
}
