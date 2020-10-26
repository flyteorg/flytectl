package create

import (
	cmdcore "github.com/lyft/flytectl/cmd/core"
	"github.com/lyft/flytectl/cmd/config"
	"github.com/spf13/cobra"
)

// CreateCommand will return create command
func CreateCommand() *cobra.Command {
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create various resource.",
	}
	createResourcesFuncs := map[string]cmdcore.CommandEntry{
		"project":    {CmdFunc: createProjectsFunc, Aliases: []string{"projects"}, ProjectDomainNotRequired: true, CustomFlags : []cmdcore.CustomFlags{
			{(config.GetCreateConfig().Name), "name","n","","Specified the name of project"},
			{(config.GetCreateConfig().ID), "id","i","","Specified the id of project"},
			{(config.GetCreateConfig().Labels), "labels","l","","Specified the labels of project"},
			{(config.GetCreateConfig().Description), "description","","","Specified the description of project"},
			{(config.GetCreateConfig().Filename), "file","f","","Specified the filename of project"},
		},},
	}
	cmdcore.AddCommands(createCmd, createResourcesFuncs)
	return createCmd
}
