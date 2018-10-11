/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package commands

import (
	"github.com/nalej/authx/internal/app/authx"
	"github.com/spf13/cobra"
	"time"
)

const DefaultExpirationDuration = "3h"

var config = authx.Config{}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the AUTHX server",
	Long:  `Launch an instance of the AUTHX server.`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		srv := authx.NewService(config)
		srv.Run()
	},
}

func init() {

	d, _ := time.ParseDuration(DefaultExpirationDuration)

	rootCmd.AddCommand(runCmd)
	runCmd.Flags().IntVar(&config.Port, "port", 8810, "Port to launch Authx server")
	runCmd.Flags().StringVar(&config.Secret, "secret", "MyLittleSecret", "Internal secret to generate Tokens")
	runCmd.Flags().DurationVar(&config.ExpirationTime, "expiration", d, "Expiration time of Tokens")
}
