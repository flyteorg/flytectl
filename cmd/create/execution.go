package create

import (
	"context"
	"fmt"

	"github.com/flyteorg/flytectl/cmd/config"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flytestdlib/logger"
)

const (
	executionShort = "Create execution resources"
	executionLong  = `
Create the executions for given workflow/task in a project and domain.

There are three steps in generating an execution.

- Generate the execution spec file using the get command.
- Update the inputs for the execution if needed.
- Run the execution by passing in the generated yaml file.

The spec file should be generated first and then run the execution using the spec file.
You can reference the flytectl get task for more details

::

 flytectl get tasks -d development -p flytectldemo core.advanced.run_merge_sort.merge  --version v2 --execFile execution_spec.yaml

The generated file would look similar to this

.. code-block:: yaml

	 iamRoleARN: ""
	 inputs:
	   sorted_list1:
	   - 0
	   sorted_list2:
	   - 0
	 kubeServiceAcct: ""
	 targetDomain: ""
	 targetProject: ""
	 task: core.advanced.run_merge_sort.merge
	 version: "v2"


The generated file can be modified to change the input values.

.. code-block:: yaml

	 iamRoleARN: 'arn:aws:iam::12345678:role/defaultrole'
	 inputs:
	   sorted_list1:
	   - 2
	   - 4
	   - 6
	   sorted_list2:
	   - 1
	   - 3
	   - 5
	 kubeServiceAcct: ""
	 targetDomain: ""
	 targetProject: ""
	 task: core.advanced.run_merge_sort.merge
	 version: "v2"

And then can be passed through the command line.
Notice the source and target domain/projects can be different.
The root project and domain flags of -p and -d should point to task/launch plans project/domain.

::

 flytectl create execution --execFile execution_spec.yaml -p flytectldemo -d development --targetProject flytesnacks

Also an execution can be relaunched by passing in current execution id.

::

 flytectl create execution --relaunch ffb31066a0f8b4d52b77 -p flytectldemo -d development

An execution can be recovered, that is recreated from the last known failure point for a previously-run workflow execution.
See :ref:` + "`ref_flyteidl.admin.ExecutionRecoverRequest`" + ` for more details.

::

 flytectl create execution --recover ffb31066a0f8b4d52b77 -p flytectldemo -d development

Generic data types are also supported for execution in similar way.Following is sample of how the inputs need to be specified while creating the execution.
As usual the spec file should be generated first and then run the execution using the spec file.

::

 flytectl get task -d development -p flytectldemo  core.type_system.custom_objects.add --execFile adddatanum.yaml

The generated file would look similar to this. Here you can see empty values dumped for generic data type x and y. 

::

    iamRoleARN: ""
    inputs:
      "x": {}
      "y": {}
    kubeServiceAcct: ""
    targetDomain: ""
    targetProject: ""
    task: core.type_system.custom_objects.add
    version: v3

Modified file with struct data populated for x and y parameters for the task core.type_system.custom_objects.add

::

  iamRoleARN: "arn:aws:iam::123456789:role/dummy"
  inputs:
    "x":
      "x": 2
      "y": ydatafory
      "z":
        1 : "foo"
        2 : "bar"
    "y":
      "x": 3
      "y": ydataforx
      "z":
        3 : "buzz"
        4 : "lightyear"
  kubeServiceAcct: ""
  targetDomain: ""
  targetProject: ""
  task: core.type_system.custom_objects.add
  version: v3

Usage
`
)

//go:generate pflags ExecutionConfig --default-var executionConfig --bind-default-var

// ExecutionConfig hold configuration for create execution flags and configuration of the actual task or workflow  to be launched.
type ExecutionConfig struct {
	// pflag section
	ExecFile        string `json:"execFile,omitempty" pflag:",file for the execution params.If not specified defaults to <<workflow/task>_name>.execution_spec.yaml"`
	TargetDomain    string `json:"targetDomain" pflag:",project where execution needs to be created.If not specified configured domain would be used."`
	TargetProject   string `json:"targetProject" pflag:",project where execution needs to be created.If not specified configured project would be used."`
	KubeServiceAcct string `json:"kubeServiceAcct" pflag:",kubernetes service account AuthRole for launching execution."`
	IamRoleARN      string `json:"iamRoleARN" pflag:",iam role ARN AuthRole for launching execution."`
	Relaunch        string `json:"relaunch" pflag:",execution id to be relaunched."`
	Recover         string `json:"recover" pflag:",execution id to be recreated from the last known failure point."`
	DryRun          bool   `json:"dryRun" pflag:",execute local operations without making any modifications (skip or mock all server communication)"`
	// Non plfag section is read from the execution config generated by get task/launchplan
	Workflow string                 `json:"workflow,omitempty"`
	Task     string                 `json:"task,omitempty"`
	Version  string                 `json:"version"`
	Inputs   map[string]interface{} `json:"inputs"`
}

type ExecutionType int

const (
	Task ExecutionType = iota
	Workflow
	Relaunch
	Recover
)

type ExecutionParams struct {
	name     string
	execType ExecutionType
}

var (
	executionConfig = &ExecutionConfig{}
)

func createExecutionCommand(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	var execParams ExecutionParams
	var err error
	sourceProject := config.GetConfig().Project
	sourceDomain := config.GetConfig().Domain
	if execParams, err = readConfigAndValidate(config.GetConfig().Project, config.GetConfig().Domain); err != nil {
		return err
	}
	var executionRequest *admin.ExecutionCreateRequest
	switch execParams.execType {
	case Relaunch:
		return relaunchExecution(ctx, execParams.name, sourceProject, sourceDomain, cmdCtx)
	case Recover:
		return recoverExecution(ctx, execParams.name, sourceProject, sourceDomain, cmdCtx)
	case Task:
		if executionRequest, err = createExecutionRequestForTask(ctx, execParams.name, sourceProject, sourceDomain, cmdCtx); err != nil {
			return err
		}
	case Workflow:
		if executionRequest, err = createExecutionRequestForWorkflow(ctx, execParams.name, sourceProject, sourceDomain, cmdCtx); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid execution type %v", execParams.execType)
	}
	if executionConfig.DryRun {
		logger.Debugf(ctx, "skipping CreateExecution request (DryRun)")
	} else {
		exec, _err := cmdCtx.AdminClient().CreateExecution(ctx, executionRequest)
		if _err != nil {
			return _err
		}
		fmt.Printf("execution identifier %v\n", exec.Id)
	}
	return nil
}
