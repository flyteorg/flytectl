package get

import (
	"context"
	"fmt"

	"github.com/flyteorg/flytectl/cmd/config"
	"github.com/flyteorg/flytectl/cmd/config/subcommand/execution"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flytectl/pkg/printer"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flytestdlib/logger"

	"github.com/golang/protobuf/proto"
)

const (
	executionShort = "Gets execution resources"
	executionLong  = `
Retrieve all executions within the project and domain (execution, executions can be used interchangeably):
::

 flytectl get execution -p flytesnacks -d development

Retrieve executions by name within the project and domain:

::

 flytectl get execution -p flytesnacks -d development oeh94k9r2r

Retrieve all the executions with filters:
::
 
 flytectl get execution -p flytesnacks -d development --filter.fieldSelector="execution.phase in (FAILED;SUCCEEDED),execution.duration<200" 

 
Retrieve executions as per the specified limit and sorting parameters:
::
  
 flytectl get execution -p flytesnacks -d development --filter.sortBy=created_at --filter.limit=1 --filter.asc
   

Retrieve executions within the project and domain in YAML format:

::

 flytectl get execution -p flytesnacks -d development -o yaml

Retrieve executions within the project and domain in JSON format:

::

 flytectl get execution -p flytesnacks -d development -o json


Get more details of the execution using the --details flag, which shows node and task executions. The default view is a tree view, and the TABLE view format is not supported on this view.

::

 flytectl get execution -p flytesnacks -d development oeh94k9r2r --details

Fetch execution details in YAML format. In this view, only node details are available. For task, send the --nodeID flag.

::

 flytectl get execution -p flytesnacks -d development oeh94k9r2r --details -o yaml

Fetch task executions on a specific node using the --nodeID flag. Use the nodeID attribute given by the node details view.

::

 flytectl get execution -p flytesnacks -d development oeh94k9r2r --nodeID n0

Task execution view is also available in YAML/JSON format. The following example showcases YAML, where the output also contains input and output data of each node.

::

 flytectl get execution -p flytesnacks -d development oeh94k9r2r --nodID n0 -o yaml

Usage
`
)

var hundredChars = 100

var executionColumns = []printer.Column{
	{Header: "Name", JSONPath: "$.id.name"},
	{Header: "Launch Plan Name", JSONPath: "$.spec.launchPlan.name"},
	{Header: "Type", JSONPath: "$.spec.launchPlan.resourceType"},
	{Header: "Phase", JSONPath: "$.closure.phase"},
	{Header: "Scheduled Time", JSONPath: "$.spec.metadata.scheduledAt"},
	{Header: "Started", JSONPath: "$.closure.startedAt"},
	{Header: "Elapsed Time", JSONPath: "$.closure.duration"},
	{Header: "Abort data (Trunc)", JSONPath: "$.closure.abortMetadata[\"cause\"]", TruncateTo: &hundredChars},
	{Header: "Error data (Trunc)", JSONPath: "$.closure.error[\"message\"]", TruncateTo: &hundredChars},
}

func ExecutionToProtoMessages(l []*admin.Execution) []proto.Message {
	messages := make([]proto.Message, 0, len(l))
	for _, m := range l {
		messages = append(messages, m)
	}
	return messages
}

func getExecutionFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	adminPrinter := printer.Printer{}
	var executions []*admin.Execution
	if len(args) > 0 {
		name := args[0]
		exec, err := cmdCtx.AdminFetcherExt().FetchExecution(ctx, name, config.GetConfig().Project, config.GetConfig().Domain)
		if err != nil {
			return err
		}
		executions = append(executions, exec)
		logger.Infof(ctx, "Retrieved %v executions", len(executions))

		if execution.DefaultConfig.Details || len(execution.DefaultConfig.NodeID) > 0 {
			// Fetching Node execution details
			nExecDetailsForView, err := getExecutionDetails(ctx, config.GetConfig().Project, config.GetConfig().Domain, name, execution.DefaultConfig.NodeID, cmdCtx)
			if err != nil {
				return err
			}
			// o/p format of table is not supported on the details. TODO: Add tree format in printer
			if config.GetConfig().MustOutputFormat() == printer.OutputFormatTABLE {
				fmt.Println("TABLE format is not supported on detailed view and defaults to tree view. Choose either json/yaml")
				nodeExecTree := createNodeDetailsTreeView(nil, nExecDetailsForView)
				fmt.Println(nodeExecTree.Print())
				return nil
			}
			return adminPrinter.PrintInterface(config.GetConfig().MustOutputFormat(), nodeExecutionColumns, nExecDetailsForView)
		}
		return adminPrinter.Print(config.GetConfig().MustOutputFormat(), executionColumns,
			ExecutionToProtoMessages(executions)...)
	}
	executionList, err := cmdCtx.AdminFetcherExt().ListExecution(ctx, config.GetConfig().Project, config.GetConfig().Domain, execution.DefaultConfig.Filter)
	if err != nil {
		return err
	}
	logger.Infof(ctx, "Retrieved %v executions", len(executionList.Executions))
	return adminPrinter.Print(config.GetConfig().MustOutputFormat(), executionColumns,
		ExecutionToProtoMessages(executionList.Executions)...)
}
