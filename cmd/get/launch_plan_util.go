package get

import (
	"context"
	"fmt"

	"github.com/flyteorg/flytectl/cmd/config"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"
)

func getAllVerOfLP(ctx context.Context, lpName string, cmdCtx cmdCore.CommandContext) ([]*admin.LaunchPlan, error) {
	tList, err := cmdCtx.AdminClient().ListLaunchPlans(ctx, &admin.ResourceListRequest{
		Limit: 1,
		Id: &admin.NamedEntityIdentifier{
			Project: config.GetConfig().Project,
			Domain:  config.GetConfig().Domain,
			Name:    lpName,
		},
		SortBy: &admin.Sort{
			Key:       "created_at",
			Direction: admin.Sort_DESCENDING,
		},
	})
	if err != nil {
		return nil, err
	}
	if len(tList.LaunchPlans) == 0 {
		return nil, fmt.Errorf("no launchplans retrieved for %v", lpName)
	}
	return tList.LaunchPlans, nil
}

func FetchLPVersionOrLatest(ctx context.Context, name string, version string, cmdCtx cmdCore.CommandContext) (*admin.LaunchPlan, error) {
	var lp *admin.LaunchPlan
	var err error
	if version == "" {
		// Fetch the latest version of the task.
		var lpVersions []*admin.LaunchPlan
		lpVersions, err = getAllVerOfLP(ctx, name, cmdCtx)
		if err != nil {
			return nil, err
		}
		lp = lpVersions[0]
	} else {
		lp, err = cmdCtx.AdminClient().GetLaunchPlan(ctx, &admin.ObjectGetRequest{
			Id: &core.Identifier{
				ResourceType: core.ResourceType_LAUNCH_PLAN,
				Project:      config.GetConfig().Project,
				Domain:       config.GetConfig().Domain,
				Name:         name,
				Version:      version,
			},
		})
		if err != nil {
			return nil, err
		}
	}
	return lp, nil
}
