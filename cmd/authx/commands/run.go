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

const DefaultExpirationDuration = "3h"

const DefaultPort = 8810

const DefaultSecret = "myLittleSecret"

var config = authx.Config{}

var secretFile = ""

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the AUTHX server",
	Long:  `Launch an instance of the AUTHX server.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if secretFile == "" {
			config.Secret = DefaultSecret
		} else {
			dat, err := ioutil.ReadFile(secretFile)
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

	rootCmd.AddCommand(runCmd)
	runCmd.Flags().IntVar(&config.Port, "port", DefaultPort, "Port to launch Authx server")
	runCmd.Flags().StringVar(&secretFile, "secret", "", "Path to internal secret to generate Tokens")
	runCmd.Flags().DurationVar(&config.ExpirationTime, "expiration", d, "Expiration time of Tokens")
}
