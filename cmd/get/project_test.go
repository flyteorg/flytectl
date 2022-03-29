package get

import (
	"fmt"
	"testing"

	"github.com/flyteorg/flytectl/cmd/config/subcommand/project"

	"github.com/flyteorg/flytectl/pkg/filters"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/stretchr/testify/assert"
)

var (
	resourceListRequestProject *admin.ProjectListRequest
	projectListResponse        *admin.Projects
	argsProject                = []string{"flyteexample"}
	project1                   *admin.Project
)

func getProjectSetup() {
	resourceListRequestProject = &admin.ProjectListRequest{}

	project1 = &admin.Project{
		Id:   "flyteexample",
		Name: "flyteexample",
		Domains: []*admin.Domain{
			{
				Id:   "development",
				Name: "development",
			},
		},
	}

	project2 := &admin.Project{
		Id:   "flytesnacks",
		Name: "flytesnacks",
		Domains: []*admin.Domain{
			{
				Id:   "development",
				Name: "development",
			},
		},
	}

	projects := []*admin.Project{project1, project2}

	projectListResponse = &admin.Projects{
		Projects: projects,
	}
}

func TestListProjectFunc(t *testing.T) {
	s := setup()
	getProjectSetup()
	project.DefaultConfig.Filter = filters.Filters{}
	s.MockAdminClient.OnListProjectsMatch(s.Ctx, resourceListRequestProject).Return(projectListResponse, nil)
	err := getProjectsFunc(s.Ctx, argsProject, s.CmdCtx)

	assert.Nil(t, err)
	s.MockAdminClient.AssertCalled(t, "ListProjects", s.Ctx, resourceListRequestProject)
}

func TestGetProjectFunc(t *testing.T) {
	s := setup()
	getProjectSetup()
	argsProject = []string{}

	project.DefaultConfig.Filter = filters.Filters{}
	s.MockAdminClient.OnListProjectsMatch(s.Ctx, resourceListRequestProject).Return(projectListResponse, nil)
	err := getProjectsFunc(s.Ctx, argsProject, s.CmdCtx)
	assert.Nil(t, err)
	s.MockAdminClient.AssertCalled(t, "ListProjects", s.Ctx, resourceListRequestProject)
}

func TestGetProjectFuncError(t *testing.T) {
	s := setup()
	getProjectSetup()
	project.DefaultConfig.Filter = filters.Filters{
		FieldSelector: "hello=",
	}
	s.MockAdminClient.OnListProjectsMatch(s.Ctx, resourceListRequestProject).Return(nil, fmt.Errorf("Please add a valid field selector"))
	err := getProjectsFunc(s.Ctx, argsProject, s.CmdCtx)
	assert.NotNil(t, err)
}
