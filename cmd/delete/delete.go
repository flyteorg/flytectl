package delete

import (
	cmdcore "github.com/lyft/flytectl/cmd/core"

	"github.com/spf13/cobra"
)

// CreateDeleteCommand will return delete command
func CreateDeleteCommand() *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete various resource.",
	}

	terminateResourcesFuncs := map[string]cmdcore.CommandEntry{
		"execution": {CmdFunc: terminateExecutionFunc, Aliases: []string{"executions"}},
	}

	cmdcore.AddCommands(deleteCmd, terminateResourcesFuncs)

	return deleteCmd
}
