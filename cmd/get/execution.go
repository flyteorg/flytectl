package get

import (
	"context"

	"github.com/flyteorg/flytectl/pkg/filters"

	"github.com/flyteorg/flytectl/cmd/config"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flytectl/pkg/printer"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flytestdlib/logger"
	"github.com/golang/protobuf/proto"
)

const (
	executionShort = "Gets execution resources"
	executionLong  = `
Retrieves all the executions within project and domain.(execution,executions can be used interchangeably in these commands)
::

 bin/flytectl get execution -p flytesnacks -d development

Retrieves execution by name within project and domain.

::

 bin/flytectl get execution -p flytesnacks -d development oeh94k9r2r

Retrieves all the execution with filters.
::
 
  bin/flytectl get execution -p flytesnacks -d development --filter.field-selector="execution.phase in (FAILED;SUCCEEDED),execution.duration<200" 
 
Retrieve specific execution with filters.
::
 
  bin/flytectl get execution -p flytesnacks -d development  y8n2wtuspj --filter.field-selector="execution.phase in (FAILED),execution.duration<200" 
 
Retrieves all the execution with limit and sorting.
::
  
   bin/flytectl get execution -p flytesnacks -d development --filter.sort-by=created_at --filter.limit=1 --filter.asc
   

Retrieves all the execution within project and domain in yaml format

::

 bin/flytectl get execution -p flytesnacks -d development -o yaml

Retrieves all the execution within project and domain in json format.

::

 bin/flytectl get execution -p flytesnacks -d development -o json

Usage
`
)

//go:generate pflags ExecutionsConfig --default-var executionConfig
var (
	executionConfig = &ExecutionsConfig{
		Filter: filters.DefaultFilter,
	}
)

// ExecutionsConfig
type ExecutionsConfig struct {
	Filter filters.Filters `json:"filter" pflag:","`
}

var executionColumns = []printer.Column{
	{Header: "Name", JSONPath: "$.id.name"},
	{Header: "Launch Plan Name", JSONPath: "$.spec.launchPlan.name"},
	{Header: "Type", JSONPath: "$.spec.launchPlan.resourceType"},
	{Header: "Phase", JSONPath: "$.closure.phase"},
	{Header: "Started", JSONPath: "$.closure.startedAt"},
	{Header: "Elapsed Time", JSONPath: "$.closure.duration"},
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
		execution, err := cmdCtx.AdminFetcherExt().FetchExecution(ctx, name, config.GetConfig().Project, config.GetConfig().Domain)
		if err != nil {
			return err
		}
		executions = append(executions, execution)
	} else {
		executionList, err := cmdCtx.AdminClient().ListExecutions(ctx, filters.BuildResourceListRequestWithName(executionConfig.Filter, ""))
		if err != nil {
			return err
		}
		executions = executionList.Executions
	}
	logger.Infof(ctx, "Retrieved %v executions", len(executions))
	err := adminPrinter.Print(config.GetConfig().MustOutputFormat(), executionColumns,
		ExecutionToProtoMessages(executions)...)
	if err != nil {
		return err
	}
	return nil
}
