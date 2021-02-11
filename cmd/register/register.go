package register

import (
	cmdcore "github.com/lyft/flytectl/cmd/core"
	"github.com/spf13/cobra"
)

// Long descriptions are whitespace sensitive when generating docs using sphinx.
const (
	registerCmdShort = "Registers tasks/workflows/launchplans from list of generated serialized files."
	registercmdLong  = `
Takes input files as serialized versions of the tasks/workflows/launchplans and registers them with flyteadmin.
Currently these input files are protobuf files generated as output from flytekit serialize.
Project & Domain are mandatory fields to be passed for registration and an optional version which defaults to v1
If the entities are already registered with flyte for the same version then registration would fail.
`
	registerFilesShort = "Registers file resources"
	registerFilesLong  = `
Registers all the serialized protobuf files including tasks, workflows and launchplans with default v1 version.
If there are already registered entities with v1 version then the command will fail immediately on the first such encounter.
::

 bin/flytectl register file  _pb_output/* -d development  -p flytesnacks

If you want to continue executing registration on other files ignoring the errors including version conflicts then pass in
the skipOnError flag.

::

 bin/flytectl register file  _pb_output/* -d development  -p flytesnacks --skipOnError

Using short format of skipOnError flag
::

 bin/flytectl register file  _pb_output/* -d development  -p flytesnacks -s

Overriding the default version v1 using version string.
::

 bin/flytectl register file  _pb_output/* -d development  -p flytesnacks -v v2

Change the o/p format has not effect on registration. The O/p is currently available only in table format.

::

 bin/flytectl register file  _pb_output/* -d development  -p flytesnacks -s -o yaml

Usage
`
)

// RegisterCommand will return register command
func RegisterCommand() *cobra.Command {
	registerCmd := &cobra.Command{
		Use:   "register",
		Short: registerCmdShort,
		Long:  registercmdLong,
	}
	registerResourcesFuncs := map[string]cmdcore.CommandEntry{
		"files": {CmdFunc: registerFromFilesFunc, Aliases: []string{"file"}, PFlagProvider: filesConfig, Short: registerFilesShort,
			Long: registerFilesLong},
	}
	cmdcore.AddCommands(registerCmd, registerResourcesFuncs)
	return registerCmd
}
