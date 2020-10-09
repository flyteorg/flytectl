package get

import (
	"context"
	"encoding/json"
	"github.com/lyft/flytectl/cmd/config"
	cmdCore "github.com/lyft/flytectl/cmd/core"
	"github.com/lyft/flytectl/pkg/printer"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/core"
)

var executionSingleStructure = map[string]string{
	"Version":    "$.id.version",
	"Name":       "$.id.name",
	"LaunchPlan": "$.spec.launchplan.name",
	"Phase":      "$.spec.phase",
	"Duration":   "$.closure.duration",
	"StartedAt":  "$.closure.started_at",
	"Workflow":   "$.closure.workflow_id.name",
	"Metadata":   "$.spec.metadata",
}

func transformSingleExecution(jsonbody []byte) (interface{}, error) {
	results := PrintableSingleExecution{}
	if err := json.Unmarshal(jsonbody, &results); err != nil {
		return results, err
	}
	return results, nil
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
		err = executionPrinter.Print(config.GetConfig().MustOutputFormat(), excution, executionSingleStructure, transformSingleExecution)
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
	err = executionPrinter.Print(config.GetConfig().MustOutputFormat(), excution, executionSingleStructure, transformSingleExecution)
	if err != nil {
		return err
	}
	return nil
}
