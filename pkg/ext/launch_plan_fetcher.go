package ext

import (
	"context"
	"fmt"

	"github.com/flyteorg/flytectl/pkg/filters"

	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"
)

// FetchAllVerOfLP fetches all the versions for give launch plan name
func (a *AdminFetcherExtClient) FetchAllVerOfLP(ctx context.Context, lpName, project, domain string, filter filters.Filters) ([]*admin.LaunchPlan, error) {
	transformFilters, err := filters.BuildResourceListRequestWithName(filter, project, domain, lpName)
	if err != nil {
		return nil, err
	}
	tList, err := a.AdminServiceClient().ListLaunchPlans(ctx, transformFilters)
	if err != nil {
		return nil, err
	}
	if len(tList.LaunchPlans) == 0 {
		return nil, fmt.Errorf("no launchplans retrieved for %v", lpName)
	}
	return tList.LaunchPlans, nil
}

// FetchAllLPs fetches all launchplans in project domain
func (a *AdminFetcherExtClient) FetchAllLPs(ctx context.Context, project, domain string, filter filters.Filters) ([]*admin.NamedEntity, error) {
	tranformFilters, err := filters.BuildNamedEntityListRequest(filter, project, domain, core.ResourceType_LAUNCH_PLAN)
	if err != nil {
		return nil, err
	}
	wList, err := a.AdminServiceClient().ListNamedEntities(ctx, tranformFilters)
	if err != nil {
		return nil, err
	}
	if len(wList.Entities) == 0 {
		return nil, fmt.Errorf("no launchplan retrieved for %v project %v domain", project, domain)
	}
	return wList.Entities, nil
}

// FetchLPLatestVersion fetches latest version for give launch plan name
func (a *AdminFetcherExtClient) FetchLPLatestVersion(ctx context.Context, name, project, domain string, filter filters.Filters) (*admin.LaunchPlan, error) {
	// Fetch the latest version of the task.
	lpVersions, err := a.FetchAllVerOfLP(ctx, name, project, domain, filter)
	if err != nil {
		return nil, err
	}
	lp := lpVersions[0]
	return lp, nil
}

// FetchLPVersion fetches particular version of launch plan
func (a *AdminFetcherExtClient) FetchLPVersion(ctx context.Context, name, version, project, domain string) (*admin.LaunchPlan, error) {
	lp, err := a.AdminServiceClient().GetLaunchPlan(ctx, &admin.ObjectGetRequest{
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
