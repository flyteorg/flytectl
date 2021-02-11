package get

import (
	cmdcore "github.com/lyft/flytectl/cmd/core"

	"github.com/spf13/cobra"
)

const (
	projectShort    = "Gets project resources"
	projectLong     = `Retrieves all the projects\n bin/flytectl get project`
	taskShort       = "Gets task resources"
	taskLong        = "Retrieves all the tasks\n bin/flytectl get tasks -p flytesnacks -d development"
	workflowShort   = "Gets task resources"
	workflowLong    = "Retrieves all the workflows\n bin/flytectl get workflows -p flytesnacks -d development"
	launchPlanShort = "Gets launch plan resources"
	launchPlanLong  = "Retrieves all the launch plans\n bin/flytectl get launchplans -p flytesnacks -d development"
	executionShort  = "Gets execution resources"
	executionLong   = "Retrieves all the executions bin/flytectl get executions -p flytesnacks -d development"
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
