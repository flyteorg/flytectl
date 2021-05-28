package get

import (
	"github.com/flyteorg/flytectl/pkg/filters"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	resourceListRequestProject *admin.ProjectListRequest
	projectListResponse *admin.Projects
	argsProject []string
	project1 *admin.Project
 )

func getProjectSetup(){
	argsProject =  []string{"flyteexample"}
	resourceListRequestProject = &admin.ProjectListRequest{
	}

	project1 = &admin.Project{
		Id: "flyteexample",
		Name: "flyteexample",
		Domains: []*admin.Domain{
			{
				Id :"development",
				Name: "development",
			},
		},
	}

	project2 := &admin.Project{
		Id: "flytesnacks",
		Name: "flytesnacks",
		Domains: []*admin.Domain{
			{
				Id :"development",
				Name: "development",
			},
		},
	}

	projects := []*admin.Project{project1, project2}

	projectListResponse = &admin.Projects{
		Projects: projects,
	}
}

func TestProjectFunc(t *testing.T) {
	setup()
	getProjectSetup()
	projectConfig.Filter = filters.Filters{}
	mockClient.OnListProjectsMatch(ctx, resourceListRequestProject).Return(projectListResponse, nil)
	err = getProjectsFunc(ctx, argsProject, cmdCtx)
	assert.Nil(t, err)
	mockClient.AssertCalled(t, "ListProjects", ctx, resourceListRequestProject)
}

func TestGetProjectFunc(t *testing.T) {
	setup()
	getProjectSetup()
	projectConfig.Filter = filters.Filters{}
	mockClient.OnListProjectsMatch(ctx, resourceListRequestProject).Return(projectListResponse, nil)
	err = getProjectsFunc(ctx, argsProject, cmdCtx)
	assert.Nil(t, err)
	mockClient.AssertCalled(t, "ListProjects", ctx, resourceListRequestProject)
}