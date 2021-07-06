package init

import (
	initConfig "github.com/flyteorg/flytectl/cmd/config/subcommand/init"
	cmdcore "github.com/flyteorg/flytectl/cmd/core"

	"github.com/spf13/cobra"
)

// Long descriptions are whitespace sensitive when generating docs using sphinx.
const (
	initCmdShort = `Used for generating config template.`
	initCmdLong  = `init config will create a config in flyte directory i.e ~/.flyte
Generate sandbox config.
	
::

 bin/flytectl init config 

Generate remote cluster config. 
	
::

 bin/flytectl init config --host="flyte.myexample.com"
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
			Short: initConfigCmdShort,
			Long:  initConfigCmdLong, PFlagProvider: initConfig.DefaultConfig},
	}

	cmdcore.AddCommands(initCmd, getResourcesFuncs)

	return initCmd
}
