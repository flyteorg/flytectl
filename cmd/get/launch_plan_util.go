package get

import (
	"context"
	"fmt"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/service"

	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"
)

// FetchLPForName fetches the launchplan give it name.
func (f FetcherImpl) FetchLPForName(ctx context.Context, adminClient service.AdminServiceClient, name, project,
	domain string) ([]*admin.LaunchPlan, error) {
	var launchPlans []*admin.LaunchPlan
	var lp *admin.LaunchPlan
	var err error
	if launchPlanConfig.Latest {
		if lp, err = f.FetchLPLatestVersion(ctx, adminClient, name, project, domain); err != nil {
			return nil, err
		}
		launchPlans = append(launchPlans, lp)
	} else if launchPlanConfig.Version != "" {
		if lp, err = f.FetchLPVersion(ctx, adminClient, name, launchPlanConfig.Version,
			project, domain); err != nil {
			return nil, err
		}
		launchPlans = append(launchPlans, lp)
	} else {
		launchPlans, err = f.FetchAllVerOfLP(ctx, adminClient, name, project, domain)
		if err != nil {
			return nil, err
		}
	}
	if launchPlanConfig.ExecFile != "" {
		// There would be atleast one launchplan object when code reaches here and hence the length
		// assertion is not required.
		lp = launchPlans[0]
		// Only write the first task from the tasks object.
		if err = CreateAndWriteExecConfigForWorkflow(lp, launchPlanConfig.ExecFile); err != nil {
			return nil, err
		}
	}
	return launchPlans, nil
}

// FetchAllVerOfLP fetches all the versions for give launchplan name
func (f FetcherImpl) FetchAllVerOfLP(ctx context.Context, adminClient service.AdminServiceClient, lpName, project,
	domain string) ([]*admin.LaunchPlan, error) {
	tList, err := adminClient.ListLaunchPlans(ctx, &admin.ResourceListRequest{
		Id: &admin.NamedEntityIdentifier{
			Project: project,
			Domain:  domain,
			Name:    lpName,
		},
		SortBy: &admin.Sort{
			Key:       "created_at",
			Direction: admin.Sort_DESCENDING,
		},
		Limit: 100,
	})
	if err != nil {
		return nil, err
	}
	if len(tList.LaunchPlans) == 0 {
		return nil, fmt.Errorf("no launchplans retrieved for %v", lpName)
	}
	return tList.LaunchPlans, nil
}


// FetchLPLatestVersion fetches latest version for give launchplan name
func (f FetcherImpl) FetchLPLatestVersion(ctx context.Context, adminClient service.AdminServiceClient, name, project,
	domain string) (*admin.LaunchPlan, error) {
	// Fetch the latest version of the task.
	lpVersions, err := f.FetchAllVerOfLP(ctx, adminClient, name, project, domain)
	if err != nil {
		return nil, err
	}
	lp := lpVersions[0]
	return lp, nil
}

func (f FetcherImpl) FetchLPVersion(ctx context.Context, adminClient service.AdminServiceClient, name, version,
	project, domain string) (*admin.LaunchPlan, error) {
	lp, err := adminClient.GetLaunchPlan(ctx, &admin.ObjectGetRequest{
		Id: &core.Identifier{
			ResourceType: core.ResourceType_LAUNCH_PLAN,
			Project:      project,
			Domain:       domain,
			Name:         name,
			Version:      version,
		},
	})
	if err != nil {
		return nil, err
	}
	return lp, nil
}
