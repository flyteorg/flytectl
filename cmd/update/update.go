package update

import (
	cmdcore "github.com/lyft/flytectl/cmd/core"

	"github.com/spf13/cobra"
)

// CreateUpdateCommand will return update command
func CreateUpdateCommand() *cobra.Command {
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update various resources.",
	}

	updateResourcesFuncs := map[string]cmdcore.CommandEntry{
		"project":    {CmdFunc: updateProjectsFunc, Aliases: []string{"projects"}, ProjectDomainNotRequired: true},
	}

	cmdcore.AddCommands(updateCmd, updateResourcesFuncs)
	updateCmd.PersistentFlags().BoolVarP(&(GetConfig().ActivateProject), "activate", "t", false, "Activates the project specified as argument.")
	updateCmd.PersistentFlags().BoolVarP(&(GetConfig().ArchiveProject), "archive", "a", false, "Activates the project specified as argument.")
	return updateCmd
}
