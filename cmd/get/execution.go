package get

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/lyft/flytectl/cmd/config"
	cmdCore "github.com/lyft/flytectl/cmd/core"
	"github.com/lyft/flytectl/pkg/printer"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
)

var executionColumns = []printer.Column{
	{"Version", "$.id.version"},
	{"Name", "$.id.name"},
	{"LaunchPlan", "$.spec.launchplan.name"},
	{"Phase", "$.spec.phase"},
	{"Duration", "$.closure.duration"},
	{"StartedAt", "$.closure.started_at"},
	{"Workflow", "$.closure.workflow_id.name"},
	{"Metadata", "$.spec.metadata"},
}

var executionSingleColumns = []printer.Column{
	{"TaskID", "$.id.taskID"},
	{"NodeExecutionID", "$.id.nodeExecutionID"},
	{"Phase", "$.spec.phase"},
	{"Duration", "$.closure.duration"},
	{"StartedAt", "$.closure.started_at"},
}

func ExecutionToProtoMessages(l []*admin.Execution) []proto.Message {
	messages := make([]proto.Message, 0, len(l))
	for _, m := range l {
		messages = append(messages, m)
	}
	return messages
}

func getExecutionFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	executionPrinter := printer.Printer{}
	excution, err := cmdCtx.AdminClient().ListExecutions(ctx, &admin.ResourceListRequest{
		Limit: 10,
		Id: &admin.NamedEntityIdentifier{
			Project: config.GetConfig().Project,
			Domain:  config.GetConfig().Domain,
		},
	})
	if err != nil {
		return err
	}
	err = executionPrinter.Print(config.GetConfig().MustOutputFormat(), executionColumns, ExecutionToProtoMessages(excution.Executions)...)
	if err != nil {
		return err
	}
	return nil

	return nil
}
