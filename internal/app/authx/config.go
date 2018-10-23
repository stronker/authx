/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package authx

import (
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

// Config is the set of required configuration parameters.
type Config struct {
	Port       int
	Secret     string
	ExpirationTime time.Duration
}

// Print information about the configuration of the application.
func (conf * Config) Print() {
	log.Info().Int("port", conf.Port).Msg("gRPC port")
	log.Info().Str("secret", strings.Repeat("*", len(conf.Secret))).Msg("Token secret")
	log.Info().Str("duration", conf.ExpirationTime.String()).Msg("Expiration time")
}