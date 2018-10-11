/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package main

import (
	"github.com/nalej/authx/cmd/authx/commands"
	"github.com/nalej/authx/version"
)

var MainVersion string

var MainCommit string

func main() {
	version.AppVersion = MainVersion
	version.Commit = MainCommit
	commands.Execute()
}
