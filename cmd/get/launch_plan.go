package get

import (
	"context"
	"github.com/lyft/flytectl/cmd/config"
	cmdCore "github.com/lyft/flytectl/cmd/core"
	"github.com/lyft/flytectl/pkg/printer"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/lyft/flytestdlib/logger"
)

func getLaunchPlanFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	launchPlanPrinter := printer.Printer{}

	launchPlans, err := cmdCtx.AdminClient().ListLaunchPlans(ctx, &admin.ResourceListRequest{
		Limit: 10,
		Id: &admin.NamedEntityIdentifier{
			Project: config.GetConfig().Project,
			Domain:  config.GetConfig().Domain,
		},
	})
	if err != nil {
		return err
	}
	if len(args) == 1 {
		name := args[0]
		logger.Debugf(ctx, "Retrieved %v excutions", len(launchPlans.LaunchPlans))
		for _, v := range launchPlans.LaunchPlans {
			if v.Id.Name == name {
				err := launchPlanPrinter.Print(config.GetConfig().MustOutputFormat(), v, launchPlanStructure, transformLaunchPlan)
				if err != nil {
					return err
				}
				return nil
			}
		}
		return nil
	}
	logger.Debugf(ctx, "Retrieved %v launch plan", len(launchPlans.LaunchPlans))
	return launchPlanPrinter.Print(config.GetConfig().MustOutputFormat(), launchPlans.LaunchPlans, launchPlanStructure, transformLaunchPlan)
	return nil
}
