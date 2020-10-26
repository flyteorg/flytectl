package create

import (
	cmdcore "github.com/lyft/flytectl/cmd/core"

	"github.com/spf13/cobra"
)

// CreateCommand will return get command
func CreateCommand() *cobra.Command {
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create various resource.",
	}
	createResourcesFuncs := map[string]cmdcore.CommandEntry{
		"project":    {CmdFunc: createProjectsFunc, Aliases: []string{"projects"}, ProjectDomainNotRequired: true, CustomFlags : []cmdcore.CustomFlags{
			{&(config.GetConfig().Name), "name","n","","Specified the name of project"},
		},},
	}

	cmdcore.AddCommands(createCmd, createResourcesFuncs)

	return createCmd
}
