package get

import (
	"context"

	"github.com/flyteorg/flytectl/cmd/config"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flytectl/pkg/adminutils"
	"github.com/flyteorg/flytectl/pkg/printer"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flytestdlib/logger"
	"github.com/golang/protobuf/proto"
)

const (
	launchPlanShort = "Gets launch plan resources"
	launchPlanLong  = `
Retrieves all the launch plans within project and domain.(launchplan,launchplans can be used interchangeably in these commands)
::

 bin/flytectl get launchplan -p flytesnacks -d development

Retrieves launch plan by name within project and domain.

::

 bin/flytectl get launchplan -p flytesnacks -d development core.basic.lp.go_greet

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
)

//go:generate pflags LaunchPlanConfig --default-var launchPlanConfig
var (
	launchPlanConfig = &LaunchPlanConfig{}
)

// LaunchPlanConfig
type LaunchPlanConfig struct {
	ExecFile string `json:"execFile" pflag:",execution file name to be used for generating execution spec of a single launchplan."`
	Version  string `json:"version" pflag:",version of the launchplan to be fetched."`
}

var launchplanColumns = []printer.Column{
	{Header: "Version", JSONPath: "$.id.version"},
	{Header: "Name", JSONPath: "$.id.name"},
	{Header: "Type", JSONPath: "$.closure.compiledTask.template.type"},
	{Header: "State", JSONPath: "$.spec.state"},
	{Header: "Schedule", JSONPath: "$.spec.entityMetadata.schedule"},
}

func LaunchplanToProtoMessages(l []*admin.LaunchPlan) []proto.Message {
	messages := make([]proto.Message, 0, len(l))
	for _, m := range l {
		messages = append(messages, m)
	}
	return messages
}

func getLaunchPlanFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	launchPlanPrinter := printer.Printer{}
	if len(args) == 1 {
		name := args[0]
		var launchPlans []*admin.LaunchPlan
		var err error
		// Right only support writing execution file for single task version.
		if launchPlanConfig.ExecFile != "" {
			var lp *admin.LaunchPlan
			if lp, err = FetchLPVersionOrLatest(ctx, name, taskConfig.Version, cmdCtx); err != nil {
				return err
			}
			launchPlans = append(launchPlans, lp)
			if err = createAndWriteExecConfigForWorkflow(lp, launchPlanConfig.ExecFile); err != nil {
				return err
			}
		} else {
			launchPlans, err = getAllVerOfLP(ctx, name, cmdCtx)
			if err != nil {
				return err
			}
		}
		logger.Debugf(ctx, "Retrieved %v launch plans", len(launchPlans))
		err = launchPlanPrinter.Print(config.GetConfig().MustOutputFormat(), launchplanColumns, LaunchplanToProtoMessages(launchPlans)...)
		if err != nil {
			return err
		}
		return nil
	}

	launchPlans, err := adminutils.GetAllNamedEntities(ctx, cmdCtx.AdminClient().ListLaunchPlanIds, adminutils.ListRequest{Project: config.GetConfig().Project, Domain: config.GetConfig().Domain})
	if err != nil {
		return err
	}
	logger.Debugf(ctx, "Retrieved %v launch plans", len(launchPlans))
	return launchPlanPrinter.Print(config.GetConfig().MustOutputFormat(), entityColumns, adminutils.NamedEntityToProtoMessage(launchPlans)...)
}
