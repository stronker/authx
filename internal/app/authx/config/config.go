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

package config

import (
	"crypto/md5"
	"fmt"
	"github.com/nalej/authx/version"
	"github.com/nalej/derrors"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"strings"
	"time"
)

// Scylla has an expiration time of 3 hours
const ttlExpirationTime = 3

// Config is the set of required configuration parameters.
type Config struct {
	// Debug level is active.
	Debug bool
	// Port where the Authx components listens for incoming connections.
	Port int
	// Secret used to sign JWT tokens.
	Secret string
	// ManagementClusterCertPath with the path of the management cluster certificate.
	ManagementClusterCertPath string
	// ManagementClusterCert with the Management cluster certificate.
	ManagementClusterCert string
	// ExpirationTime for JWT tokens.
	ExpirationTime time.Duration
	// DeviceExpirationTime for device JWT tokens.
	DeviceExpirationTime time.Duration
	// EdgeControllerExpTime with the expiration time for Edge Controller join tokens.
	EdgeControllerExpTime time.Duration
	// Use in-memory providers
	UseInMemoryProviders bool
	// Use scyllaDBProviders
	UseDBScyllaProviders bool
	// ScyllaDBAddress with the database URL
	ScyllaDBAddress string
	// ScyllaDBPort with the database port.
	ScyllaDBPort int
	// DataBase KeySpace
	KeySpace string
	// CACertPath with the path of the CA.
	CACertPath string
	// CAPrivateKeyPath with the path of the private key for the CA.
	CAPrivateKeyPath string
}

func (conf *Config) Validate() derrors.Error {
	if conf.Port <= 0 {
		return derrors.NewInvalidArgumentError("port must be specified")
	}
	if conf.UseDBScyllaProviders {
		if conf.ScyllaDBAddress == "" {
			return derrors.NewInvalidArgumentError("address must be specified to use dbScylla Providers")
		}
		if conf.KeySpace == "" {
			return derrors.NewInvalidArgumentError("keyspace must be specified to use dbScylla Providers")
		}
		if conf.ScyllaDBPort <= 0 {
			return derrors.NewInvalidArgumentError("port must be specified to use dbScylla Providers ")
		}
	}
	if !conf.UseDBScyllaProviders && !conf.UseInMemoryProviders {
		return derrors.NewInvalidArgumentError("a type of provider must be selected")
	}

	if conf.ExpirationTime.Hours() > ttlExpirationTime {
		return derrors.NewInvalidArgumentError("currently the duration can not be longer than 3h. Scylla has a 3 hours TTL")
	}
	if conf.DeviceExpirationTime.Hours() > ttlExpirationTime {
		return derrors.NewInvalidArgumentError("currently the duration of device tokens can not be longer than 3h. Scylla has a 3 hours TTL")
	}
	if conf.EdgeControllerExpTime.Hours() > ttlExpirationTime {
		return derrors.NewInvalidArgumentError("currently the duration of edge controller join tokens cannot be longer than 3h. Scylla has a 3 hours TTL")
	}

	// Load server certificate
	if conf.ManagementClusterCertPath != "" {
		err := conf.loadCert()
		if err != nil {
			return err
		}
	}

	if conf.CACertPath == "" || conf.CAPrivateKeyPath == "" {
		return derrors.NewInvalidArgumentError("caCertPath and caPrivateKey cannot be empty")
	}

	return nil
}

// LoadCert loads the management cluster certificate in memory.
func (conf *Config) loadCert() derrors.Error {
	content, err := ioutil.ReadFile(conf.ManagementClusterCertPath)
	if err != nil {
		return derrors.AsError(err, "cannot load management cluster certificate")
	}
	conf.ManagementClusterCert = string(content)
	return nil
}

// Print information about the configuration of the application.
func (conf *Config) Print() {
	log.Info().Str("app", version.AppVersion).Str("commit", version.Commit).Msg("Version")
	log.Info().Int("port", conf.Port).Msg("gRPC port")
	log.Info().Str("secret", strings.Repeat("*", len(conf.Secret))).Msg("JWT Token secret")
	if conf.ManagementClusterCert != "" {
		log.Info().Str("md5", fmt.Sprintf("%x", md5.Sum([]byte(conf.ManagementClusterCert)))).Msg("Management cluster server certificate")
	} else {
		log.Warn().Msg("Management cluster server certificate is not set")
	}
	log.Info().Str("duration", conf.ExpirationTime.String()).Msg("JWT Expiration time")
	log.Info().Str("duration", conf.DeviceExpirationTime.String()).Msg("Device expiration time")
	log.Info().Str("duration", conf.EdgeControllerExpTime.String()).Msg("Edge controller join token expiration time")

	if conf.UseInMemoryProviders {
		log.Info().Bool("UseInMemoryProviders", conf.UseInMemoryProviders).Msg("Using in-memory providers")
	}
	if conf.UseDBScyllaProviders {
		log.Info().Bool("UseDBScyllaProviders", conf.UseDBScyllaProviders).Msg("using dbScylla providers")
		log.Info().Str("URL", conf.ScyllaDBAddress).Str("KeySpace", conf.KeySpace).Int("Port", conf.ScyllaDBPort).Msg("ScyllaDB")
	}
	log.Info().Str("CA Path", conf.CACertPath).Str("CA PK Path", conf.CAPrivateKeyPath).Msg("CA files")
}
