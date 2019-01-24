/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package commands

import (
	"github.com/nalej/authx/internal/app/authx"
	"github.com/spf13/cobra"
	"io/ioutil"
	"time"
)

// DefaultExpirationDuration is the default token expiration duration
const DefaultExpirationDuration = "3h"
const DefaultDeviceExpiration = "10m"

// DefaultPort is the default port where the service is deployed
const DefaultPort = 8810

// DefaultSecret is the default secret that Authx uses to sign the token.
const DefaultSecret = "myLittleSecret"

var config = authx.Config{}

var secretPath = ""

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the AUTHX server",
	Long:  `Launch an instance of the AUTHX server.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if secretPath == "" {
			config.Secret = DefaultSecret
		} else {
			dat, err := ioutil.ReadFile(secretPath)
			if err != nil {
				panic(err)
			}
			config.Secret = string(dat)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		srv := authx.NewService(config)
		srv.Run()
	},
}

func init() {

	d, _ := time.ParseDuration(DefaultExpirationDuration)
	e, _ := time.ParseDuration(DefaultDeviceExpiration)

	rootCmd.AddCommand(runCmd)
	runCmd.Flags().IntVar(&config.Port, "port", DefaultPort, "Port to launch Authx server")
	runCmd.Flags().StringVar(&secretPath, "secret", "", "Path to internal secret to generate Tokens")
	runCmd.Flags().DurationVar(&config.ExpirationTime, "expiration", d, "Expiration time of Tokens. No more than 3 hours allowed")
	runCmd.Flags().DurationVar(&config.DeviceExpirationTime, "deviceExpiration", e, "Expiration time of devices Tokens")

	runCmd.Flags().BoolVar(&config.UseInMemoryProviders, "userInMemoryProviders", false, "Whether in-memory providers should be used. ONLY for development")
	runCmd.Flags().BoolVar(&config.UseDBScyllaProviders, "useDBScyllaProviders", true, "Whether dbscylla providers should be used")
	runCmd.Flags().StringVar(&config.ScyllaDBAddress, "scyllaDBAddress", "", "address to connect to scylla database")
	runCmd.Flags().IntVar(&config.ScyllaDBPort, "scyllaDBPort", 9042, "port to connect to scylla database")
	runCmd.Flags().StringVar(&config.KeySpace, "scyllaDBKeyspace", "", "keyspace of scylla database")
}
