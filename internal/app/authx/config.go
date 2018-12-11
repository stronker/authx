/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package authx

import (
	"github.com/nalej/derrors"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

// Config is the set of required configuration parameters.
type Config struct {
	Port       int
	Secret     string
	ExpirationTime time.Duration
	// Use in-memory providers
	UseInMemoryProviders bool
	// Use scyllaDBProviders
	UseDBScyllaProviders bool
	// Database Address
	ScyllaDBAddress string
	ScyllaDBPort int
	// DataBase KeySpace
	KeySpace string
}

func (conf * Config) Validate() derrors.Error {
	if conf.Port <= 0 {
		return derrors.NewInvalidArgumentError("port must be specified")
	}
	if conf.UseDBScyllaProviders {
		if conf.ScyllaDBAddress == "" {
			return derrors.NewInvalidArgumentError("address must be specified to use dbScylla Providers")
		}
		if  conf.KeySpace == "" {
			return derrors.NewInvalidArgumentError("keyspace must be specified to use dbScylla Providers")
		}
		if conf.ScyllaDBPort <= 0 {
			return derrors.NewInvalidArgumentError("port must be specified to use dbScylla Providers ")
		}
	}
	if !conf.UseDBScyllaProviders && !conf.UseInMemoryProviders {
		return derrors.NewInvalidArgumentError("a type of provider must be selected")
	}
	return nil
}

// Print information about the configuration of the application.
func (conf * Config) Print() {
	log.Info().Int("port", conf.Port).Msg("gRPC port")
	log.Info().Str("secret", strings.Repeat("*", len(conf.Secret))).Msg("Token secret")
	log.Info().Str("duration", conf.ExpirationTime.String()).Msg("Expiration time")
	if conf.UseInMemoryProviders {
		log.Info().Bool("UseInMemoryProviders", conf.UseInMemoryProviders).Msg("Using in-memory providers")
	}
	if conf.UseDBScyllaProviders {
		log.Info().Bool("UseDBScyllaProviders", conf.UseDBScyllaProviders).Msg("using dbScylla providers")
		log.Info().Str("URL", conf.ScyllaDBAddress).Str("KeySpace", conf.KeySpace).Int("Port", conf.ScyllaDBPort).Msg("ScyllaDB")
	}
}