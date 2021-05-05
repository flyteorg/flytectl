package impl

import (
	"context"
	"fmt"

	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"
)

func (f FetcherImpl) FetchAllVerOfTask(ctx context.Context, name, project, domain string) ([]*admin.Task, error) {
	tList, err := f.ListTasks(ctx, &admin.ResourceListRequest{
		Id: &admin.NamedEntityIdentifier{
			Project: project,
			Domain:  domain,
			Name:    name,
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
	if len(tList.Tasks) == 0 {
		return nil, fmt.Errorf("no tasks retrieved for %v", name)
	}
	return tList.Tasks, nil
}

func (f FetcherImpl) FetchTaskLatestVersion(ctx context.Context, name, project, domain string) (*admin.Task, error) {
	var t *admin.Task
	var err error
	// Fetch the latest version of the task.
	var taskVersions []*admin.Task
	taskVersions, err = f.FetchAllVerOfTask(ctx, name, project, domain)
	if err != nil {
		return nil, err
	}
	t = taskVersions[0]
	return t, nil
}

func (f FetcherImpl) FetchTaskVersion(ctx context.Context, name, version, project, domain string) (*admin.Task, error) {
	t, err := f.GetTask(ctx, &admin.ObjectGetRequest{
		Id: &core.Identifier{
			ResourceType: core.ResourceType_TASK,
			Project:      project,
			Domain:       domain,
			Name:         name,
			Version:      version,
		},
	})
	if err != nil {
		return nil, err
	}
	return t, nil
}
