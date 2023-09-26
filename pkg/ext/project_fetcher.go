package ext

import (
	"context"
	"fmt"

	"github.com/flyteorg/flytectl/pkg/filters"

	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
)

func (a *AdminFetcherExtClient) ListProjects(ctx context.Context, filter filters.Filters) (*admin.Projects, error) {
	transformFilters, err := filters.BuildProjectListRequest(filter)
	if err != nil {
		return nil, err
	}
	e, err := a.AdminServiceClient().ListProjects(ctx, transformFilters)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (a *AdminFetcherExtClient) GetProjectById(ctx context.Context, projectId string) (*admin.Project, error) {
	if projectId == "" {
		return nil, fmt.Errorf("GetProjectById: projectId is empty")
	}

	response, err := a.AdminServiceClient().ListProjects(ctx, &admin.ProjectListRequest{
		Limit:   1,
		Filters: fmt.Sprintf("eq(identifier,%s)", filters.EscapeValue(projectId)),
	})
	if err != nil {
		return nil, err
	}

	if len(response.Projects) == 0 {
		return nil, NewNotFoundError("project %s", projectId)
	}

	return response.Projects[0], nil
}
