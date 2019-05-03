/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package commands

import (
	"github.com/nalej/authx/internal/app/authx"
	"github.com/nalej/authx/internal/app/authx/config"
	"github.com/spf13/cobra"
	"io/ioutil"
	"time"
)

// DefaultExpirationDuration is the default token expiration duration
const DefaultExpirationDuration = "3h"
const DefaultDeviceExpiration = "10m"
const DefaultEdgeControllerJoinExpiration = "1h"

// DefaultPort is the default port where the service is deployed
const DefaultPort = 8810

// DefaultSecret is the default secret that Authx uses to sign the token.
const DefaultSecret = "myLittleSecret"

var cfg = config.Config{}

var secretPath = ""

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the AUTHX server",
	Long:  `Launch an instance of the AUTHX server.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if secretPath == "" {
			cfg.Secret = DefaultSecret
		} else {
			dat, err := ioutil.ReadFile(secretPath)
			if err != nil {
				panic(err)
			}
			cfg.Secret = string(dat)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cfg.Debug = debugLevel
		srv := authx.NewService(cfg)
		srv.Run()
	},
}

func init() {

	d, _ := time.ParseDuration(DefaultExpirationDuration)
	e, _ := time.ParseDuration(DefaultDeviceExpiration)
	ece, _ := time.ParseDuration(DefaultEdgeControllerJoinExpiration)

	rootCmd.AddCommand(runCmd)
	runCmd.Flags().IntVar(&cfg.Port, "port", DefaultPort, "Port to launch Authx server")
	runCmd.Flags().StringVar(&secretPath, "secret", "", "Path to internal secret to generate Tokens")
	runCmd.Flags().DurationVar(&cfg.ExpirationTime, "expiration", d, "Expiration time of Tokens. No more than 3 hours allowed")
	runCmd.Flags().DurationVar(&cfg.DeviceExpirationTime, "deviceExpiration", e, "Expiration time of devices Tokens")
	runCmd.Flags().DurationVar(&cfg.EdgeControllerExpTime, "edgeControllerJoinExpiration", ece, "Expiration time of Edge Controller join tokens")

	runCmd.Flags().BoolVar(&cfg.UseInMemoryProviders, "userInMemoryProviders", false, "Whether in-memory providers should be used. ONLY for development")
	runCmd.Flags().BoolVar(&cfg.UseDBScyllaProviders, "useDBScyllaProviders", true, "Whether dbscylla providers should be used")
	runCmd.Flags().StringVar(&cfg.ScyllaDBAddress, "scyllaDBAddress", "", "address to connect to scylla database")
	runCmd.Flags().IntVar(&cfg.ScyllaDBPort, "scyllaDBPort", 9042, "port to connect to scylla database")
	runCmd.Flags().StringVar(&cfg.KeySpace, "scyllaDBKeyspace", "", "keyspace of scylla database")

	runCmd.Flags().StringVar(&cfg.ManagementClusterCertPath, "managementClusterCertPath", "", "Server certificate that joining entities can use for authentication")
}
