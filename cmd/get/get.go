package get

import (
	cmdcore "github.com/lyft/flytectl/cmd/core"

	"github.com/spf13/cobra"
)

// Long descriptions are whitespace sensitive when generating docs using sphinx.
const (
	getLong = `
Example get projects.
::

 bin/flytectl get project

`
	projectShort = "Gets project resources"
	projectLong  = `
Retrieves all the projects.(project,projects can be used interchangeably in these commands)
::

 bin/flytectl get project

Retrieves project by name

::

 bin/flytectl get project flytesnacks

Retrieves project by filters
::

 Not yet implemented

Retrieves all the projects in yaml format

::

 bin/flytectl get project -o yaml

Retrieves all the projects in json format

::

 bin/flytectl get project -o json

Usage
`
	taskShort = "Gets task resources"
	taskLong  = `
Retrieves all the task within project and domain.(task,tasks can be used interchangeably in these commands)
::

 bin/flytectl get task -p flytesnacks -d development

Retrieves task by name within project and domain.

::

 bin/flytectl task -p flytesnacks -d development square

Retrieves project by filters.
::

 Not yet implemented

Retrieves all the tasks within project and domain in yaml format.

::

 bin/flytectl get task -p flytesnacks -d development -o yaml

Retrieves all the tasks within project and domain in json format.

::

 bin/flytectl get task -p flytesnacks -d development -o json

Usage
`
	workflowShort = "Gets workflow resources"
	workflowLong  = `
Retrieves all the workflows within project and domain.(workflow,workflows can be used interchangeably in these commands)
::

 bin/flytectl get workflow -p flytesnacks -d development

Retrieves workflow by name within project and domain.

::

 bin/flytectl workflow -p flytesnacks -d development  recipes.plugins.k8s_spark.pyspark_pi.my_spark

Retrieves workflow by filters. 
::

 Not yet implemented

Retrieves all the workflow within project and domain in yaml format.

::

 bin/flytectl get workflow -p flytesnacks -d development -o yaml

Retrieves all the workflow within project and domain in json format.

::

 bin/flytectl get workflow -p flytesnacks -d development -o json

Usage
`
	launchPlanShort = "Gets launch plan resources"
	launchPlanLong  = `
Retrieves all the launch plans within project and domain.(launchplan,launchplans can be used interchangeably in these commands)
::

 bin/flytectl get launchplan -p flytesnacks -d development

Retrieves launch plan by name within project and domain.

::

 bin/flytectl launchplan -p flytesnacks -d development recipes.core.basic.lp.my_wf

Retrieves launchplan by filters.
::

 Not yet implemented

Retrieves all the launchplan within project and domain in yaml format.

::

 bin/flytectl get launchplan -p flytesnacks -d development -o yaml

Retrieves all the launchplan within project and domain in json format

::

 bin/flytectl get launchplan -p flytesnacks -d development -o json

Usage
`
	executionShort = "Gets execution resources"
	executionLong  = `
Retrieves all the executions within project and domain.(execution,executions can be used interchangeably in these commands)
::

 bin/flytectl get execution -p flytesnacks -d development

Retrieves execution by name within project and domain.

::

 bin/flytectl execution -p flytesnacks -d development oeh94k9r2r

Retrieves execution by filters
::

 Not yet implemented

Retrieves all the execution within project and domain in yaml format

::

 bin/flytectl get execution -p flytesnacks -d development -o yaml

Retrieves all the execution within project and domain in json format.

::

 bin/flytectl get execution -p flytesnacks -d development -o json

Usage
`
)

// CreateGetCommand will return get command
func CreateGetCommand() *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get",
		Short: getLong,
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
