package get

import (
	cmdcore "github.com/lyft/flytectl/cmd/core"

	"github.com/spf13/cobra"
)

// CreateGetCommand will return get command
func CreateGetCommand() *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Retrieve various resource.",
	}

	getResourcesFuncs := map[string]cmdcore.CommandEntry{
		"project": {CmdFunc: getProjectsFunc, Aliases: []string{"projects"}, ProjectDomainNotRequired: true,
			Short: "Gets project resources",
			Long:  "Retrieves all the projects"},
		"task":       {CmdFunc: getTaskFunc, Aliases: []string{"tasks"},Short: "Gets task resources",
			Long:  "Retrieves all the tasks"},
		"workflow":   {CmdFunc: getWorkflowFunc, Aliases: []string{"workflows"}, Short: "Gets workflow resources",
			Long:  "Retrieves all the workflows"},
		"launchplan": {CmdFunc: getLaunchPlanFunc, Aliases: []string{"launchplans"}, Short: "Gets launchplan resources",
			Long:  "Retrieves all the launchplans"},
		"execution":  {CmdFunc: getExecutionFunc, Aliases: []string{"executions"}, Short: "Gets execution resources",
			Long:  "Retrieves all the executions"},
	}

	cmdcore.AddCommands(getCmd, getResourcesFuncs)

	return getCmd
}
