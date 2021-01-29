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
		"activate-project":    {CmdFunc: activateProjectFunc, Aliases: []string{"activate"}},
		"archive-project":    {CmdFunc: archiveProjectFunc, Aliases: []string{"archive"}},
	}

	cmdcore.AddCommands(updateCmd, updateResourcesFuncs)
	return updateCmd
}
