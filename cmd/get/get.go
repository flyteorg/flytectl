package get

import (
	cmdcore "github.com/lyft/flytectl/cmd/core"

	"github.com/spf13/cobra"
)

const (
	projectShort    = "Gets project resources"
	projectLong     = "Retrieves all the projects"
	taskShort       = "Gets task resources"
	taskLong        = "Retrieves all the tasks"
	workflowShort   = "Gets task resources"
	workflowLong    = "Retrieves all the tasks"
	launchPlanShort = "Gets launch plan resources"
	launchPlanLong  = "Retrieves all the launch plans"
	executionShort  = "Gets execution resources"
	executionLong   = "Retrieves all the executions"
)

// CreateGetCommand will return get command
func CreateGetCommand() *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Retrieve various resource.",
	}

	getResourcesFuncs := map[string]cmdcore.CommandEntry{
		"project": {CmdFunc: getProjectsFunc, Aliases: []string{"projects"}, ProjectDomainNotRequired: true,
			Short: projectShort,
			Long:  projectLong},
		"task": {CmdFunc: getTaskFunc, Aliases: []string{"tasks"}, Short: taskShort,
			Long: taskLong},
		"workflow": {CmdFunc: getWorkflowFunc, Aliases: []string{"workflows"}, Short: workflowShort,
			Long: workflowLong},
		"launchplan": {CmdFunc: getLaunchPlanFunc, Aliases: []string{"launchplans"}, Short: launchPlanShort,
			Long: launchPlanLong},
		"execution": {CmdFunc: getExecutionFunc, Aliases: []string{"executions"}, Short: executionShort,
			Long: executionLong},
	}

	cmdcore.AddCommands(getCmd, getResourcesFuncs)

	return getCmd
}
