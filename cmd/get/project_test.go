package get

import (
	"io"
	"testing"

	"github.com/flyteorg/flytectl/cmd/config"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flytectl/cmd/testutils"
	"github.com/flyteorg/flyteidl/clients/go/admin/mocks"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/stretchr/testify/assert"
)

func TestListProjectFunc(t *testing.T) {
	ctx := testutils.Ctx
	config.GetConfig().Output = output
	config.GetConfig().SortBy = ""
	var args []string
	mockClient := new(mocks.AdminServiceClient)
	mockOutStream := new(io.Writer)
	cmdCtx := cmdCore.NewCommandContext(mockClient, *mockOutStream)
	projectListRequest := &admin.ProjectListRequest{
		Limit: 100,
	}
	projectResponse := &admin.Project{
		Id:   "flytesnacks",
		Name: "flytesnacks",
		Domains: []*admin.Domain{
			{
				Id:   "production",
				Name: "production",
			},
			{
				Id:   "staging",
				Name: "staging",
			},
			{
				Id:   "development",
				Name: "development",
			},
		},
	}
	projects := []*admin.Project{projectResponse}
	projectList := &admin.Projects{
		Projects: projects,
	}
	mockClient.OnListProjectsMatch(ctx, projectListRequest).Return(projectList, nil)
	err := getProjectsFunc(ctx, args, cmdCtx)
	assert.Nil(t, err)
	mockClient.AssertCalled(t, "ListProjects", ctx, projectListRequest)
}
