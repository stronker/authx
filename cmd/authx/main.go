/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package main

import (
	"github.com/nalej/authx/cmd/authx/commands"
	"github.com/nalej/authx/version"
)

// MainVersion is the variable to store the version of the project.
var MainVersion string

// MainCommit is the variable to store the current commit.
var MainCommit string

func main() {
	version.AppVersion = MainVersion
	version.Commit = MainCommit
	commands.Execute()
}
