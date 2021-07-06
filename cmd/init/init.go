package init

import (
	initConfig "github.com/flyteorg/flytectl/cmd/config/subcommand/init"
	cmdcore "github.com/flyteorg/flytectl/cmd/core"

	"github.com/spf13/cobra"
)

// Long descriptions are whitespace sensitive when generating docs using sphinx.
const (
	initCmdShort = `Used for generating config template.`
	initCmdLong  = `

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
