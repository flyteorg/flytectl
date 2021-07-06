package init

import (
	initConfig "github.com/flyteorg/flytectl/cmd/config/subcommand/init"
	cmdcore "github.com/flyteorg/flytectl/cmd/core"

	"github.com/spf13/cobra"
)

// Long descriptions are whitespace sensitive when generating docs using sphinx.
const (
	initCmdShort = `Generates flytectl config file in the user's home directory.`
	initCmdLong  = `Creates a flytectl config file in flyte directory i.e ~/.flyte
	
Generate sandbox config. Flyte Sandbox is a fully standalone minimal environment for running Flyte. Read more about sandbox https://docs.flyte.org/en/latest/deployment/sandbox.html

::

 bin/flytectl init config 

Generate remote cluster config. Read more about the remote deployment https://docs.flyte.org/en/latest/deployment/index.html
	
::

 bin/flytectl init config --host=flyte.myexample.com
`
)

// CreateInitCommand will return init command
func CreateInitCommand() *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: initCmdShort,
		Long:  initCmdLong,
	}

	getResourcesFuncs := map[string]cmdcore.CommandEntry{
		"config": {CmdFunc: configInitFunc, Aliases: []string{""}, ProjectDomainNotRequired: true,
			Short: initCmdShort,
			Long:  initCmdLong, PFlagProvider: initConfig.DefaultConfig},
	}

	cmdcore.AddCommands(initCmd, getResourcesFuncs)

	return initCmd
}
