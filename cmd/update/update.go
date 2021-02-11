package update

import (
	cmdcore "github.com/lyft/flytectl/cmd/core"

	"github.com/spf13/cobra"
)

const (
	updateUse = "update"
	updateShort = "Update various resources."
	projectShort = "Updates project resources"
	projectLong = "Updates the project according the flags passed.Allows you to archive or activate a project"

)
// CreateUpdateCommand will return update command
func CreateUpdateCommand() *cobra.Command {
	updateCmd := &cobra.Command{
		Use:   updateUse,
		Short: updateShort,
	}

	updateResourcesFuncs := map[string]cmdcore.CommandEntry{
		"project":    {CmdFunc: updateProjectsFunc, Aliases: []string{"projects"}, ProjectDomainNotRequired: true, PFlagProvider: projectConfig,
			Short: projectShort,
			Long:  projectLong},
	}

	cmdcore.AddCommands(updateCmd, updateResourcesFuncs)
	return updateCmd
}
