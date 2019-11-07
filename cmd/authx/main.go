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
