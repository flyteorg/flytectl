package get

import (
	"github.com/lyft/flytectl/cmd/core"

	"github.com/spf13/cobra"
)

func CreateGetCommand() *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Retrieve various resource.",
	}

	getResourcesFuncs := map[string]cmdCore.CommandFunc{
		"projects": getProjectsFunc,
		"tasks":    getTaskFunc,
		"workflows":    getWorkflowFunc,
	}

	cmdCore.AddCommands(getCmd, getResourcesFuncs)

	return getCmd
}


