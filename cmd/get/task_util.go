package get

import (
	"context"
	"fmt"

	"github.com/flyteorg/flytectl/cmd/config"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"
)

func getAllVerOfTask(ctx context.Context, name string, cmdCtx cmdCore.CommandContext) ([]*admin.Task, error) {
	tList, err := cmdCtx.AdminClient().ListTasks(ctx, &admin.ResourceListRequest{
		Limit: 1,
		Id: &admin.NamedEntityIdentifier{
			Project: config.GetConfig().Project,
			Domain:  config.GetConfig().Domain,
			Name:    name,
		},
		SortBy: &admin.Sort{
			Key:       "created_at",
			Direction: admin.Sort_DESCENDING,
		},
	})
	if err != nil {
		return nil, err
	}
	if len(tList.Tasks) == 0 {
		return nil, fmt.Errorf("no tasks retrieved for %v", name)
	}
	return tList.Tasks, nil
}

func FetchTaskVersionOrLatest(ctx context.Context, name string, version string, cmdCtx cmdCore.CommandContext) (*admin.Task, error) {
	var t *admin.Task
	var err error
	if version == "" {
		// Fetch the latest version of the task.
		var taskVersions []*admin.Task
		taskVersions, err = getAllVerOfTask(ctx, name, cmdCtx)
		if err != nil {
			return nil, err
		}
		t = taskVersions[0]
	} else {
		t, err = cmdCtx.AdminClient().GetTask(ctx, &admin.ObjectGetRequest{
			Id: &core.Identifier{
				ResourceType: core.ResourceType_TASK,
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
	return t, nil
}
