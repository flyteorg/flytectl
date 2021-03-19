package get

import (
	"context"
	"fmt"

	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"
)

func getAllVerOfTask(ctx context.Context, name string, project string, domain string, cmdCtx cmdCore.CommandContext) ([]*admin.Task, error) {
	tList, err := cmdCtx.AdminClient().ListTasks(ctx, &admin.ResourceListRequest{
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

func fetchTaskLatestVersion(ctx context.Context, name string, project string, domain string, cmdCtx cmdCore.CommandContext) (*admin.Task, error) {
	var t *admin.Task
	var err error
	// Fetch the latest version of the task.
	var taskVersions []*admin.Task
	taskVersions, err = getAllVerOfTask(ctx, name, project, domain, cmdCtx)
	if err != nil {
		return nil, err
	}
	t = taskVersions[0]
	return t, nil
}

func FetchTaskVersion(ctx context.Context, name string, version string, project string, domain string, cmdCtx cmdCore.CommandContext) (*admin.Task, error) {
	t, err := cmdCtx.AdminClient().GetTask(ctx, &admin.ObjectGetRequest{
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
