package get

import (
	"context"

	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/core"

	"github.com/golang/protobuf/proto"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/lyft/flytestdlib/logger"

	"github.com/lyft/flytectl/cmd/config"
	cmdCore "github.com/lyft/flytectl/cmd/core"
	"github.com/lyft/flytectl/pkg/printer"
)

var executionColumns = []printer.Column{
	{"Name", "$.id.name"},
	{"Workflow Name", "$.closure.workflowId.name"},
	{"Type", "$.closure.workflowId.resourceType"},
	{"Phase", "$.closure.phase"},
	{"Started", "$.closure.startedAt"},
	{"Elapsed Time", "$.closure.duration"},
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
		execution, err := cmdCtx.AdminClient().GetExecution(ctx, &admin.WorkflowExecutionGetRequest{
			Id: &core.WorkflowExecutionIdentifier{
				Project: config.GetConfig().Project,
				Domain:  config.GetConfig().Domain,
				Name:    name,
			},
		})
		if err != nil {
			return err
		}
		executions = append(executions, execution)
	} else {
		executionList, err := cmdCtx.AdminClient().ListExecutions(ctx, &admin.ResourceListRequest{
			Limit: 100,
			Id: &admin.NamedEntityIdentifier{
				Project: config.GetConfig().Project,
				Domain:  config.GetConfig().Domain,
			},
		})
		if err != nil {
			return err
		}
		executions = executionList.Executions
	}
	logger.Infof(ctx, "Retrieved %v executions", len(executions))
	err := adminPrinter.Print(config.GetConfig().MustOutputFormat(), executionColumns, ExecutionToProtoMessages(executions)...)
	if err != nil {
		return err
	}
	return nil
}
