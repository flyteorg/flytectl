package get

import (
	"context"
	"github.com/lyft/flytectl/cmd/config"
	cmdCore "github.com/lyft/flytectl/cmd/core"
	"github.com/lyft/flytectl/pkg/printer"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/core"
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


func getExecutionFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	executionPrinter := printer.Printer{}

	if len(args) == 1 {
		name := args[0]
		excution, err := cmdCtx.AdminClient().GetTaskExecution(ctx, &admin.TaskExecutionGetRequest{
			Id: &core.TaskExecutionIdentifier{
				TaskId: &core.Identifier{
					Project : config.GetConfig().Project,
					Domain:  config.GetConfig().Domain,
					Name : name,
				},
			},
		})
		if err != nil {
			return err
		}
		err = executionPrinter.Print(config.GetConfig().MustOutputFormat(), excution, executionSingleColumns)
		if err != nil {
					return err
				}

		return nil
	}
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
	err = executionPrinter.Print(config.GetConfig().MustOutputFormat(), excution, executionColumns)
	if err != nil {
		return err
	}
	return nil
}
