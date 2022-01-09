package create

import (
	"fmt"
	"testing"

	"github.com/flyteorg/flytectl/cmd/config/subcommand/project"

	"github.com/flyteorg/flytectl/cmd/testutils"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const projectValue = "dummyProject"

var (
	projectRegisterRequest *admin.ProjectRegisterRequest
)

func createProjectSetup() {
	ctx = testutils.Ctx
	cmdCtx = testutils.CmdCtx
	mockClient = testutils.MockClient
	projectRegisterRequest = &admin.ProjectRegisterRequest{
		Project: &admin.Project{
			Id:          projectValue,
			Name:        projectValue,
			Description: "",
			Labels: &admin.Labels{
				Values: map[string]string{},
			},
		},
	}
	project.DefaultCreateConfig.ID = ""
	project.DefaultCreateConfig.Name = ""
	project.DefaultCreateConfig.Labels = map[string]string{}
	project.DefaultCreateConfig.Description = ""
}
func TestCreateProjectFunc(t *testing.T) {
	setup()
	createProjectSetup()
	defer tearDownAndVerify(t, "project Created successfully")
	project.DefaultCreateConfig.ID = projectValue
	project.DefaultCreateConfig.Name = projectValue
	project.DefaultCreateConfig.Labels = map[string]string{}
	project.DefaultCreateConfig.Description = ""
	mockClient.OnRegisterProjectMatch(ctx, projectRegisterRequest).Return(nil, nil)
	err := createProjectsCommand(ctx, args, cmdCtx)
	assert.Nil(t, err)
	mockClient.AssertCalled(t, "RegisterProject", ctx, projectRegisterRequest)
}

func TestEmptyProjectID(t *testing.T) {
	setup()
	createProjectSetup()
	defer tearDownAndVerify(t, "")
	project.DefaultCreateConfig.Name = projectValue
	project.DefaultCreateConfig.Labels = map[string]string{}
	mockClient.OnRegisterProjectMatch(ctx, projectRegisterRequest).Return(nil, nil)
	err := createProjectsCommand(ctx, args, cmdCtx)
	assert.Equal(t, fmt.Errorf("project ID is required flag"), err)
	mockClient.AssertNotCalled(t, "RegisterProject", ctx, mock.Anything)
}

func TestEmptyProjectName(t *testing.T) {
	setup()
	createProjectSetup()
	defer tearDownAndVerify(t, "")
	project.DefaultCreateConfig.ID = projectValue
	project.DefaultCreateConfig.Labels = map[string]string{}
	project.DefaultCreateConfig.Description = ""
	mockClient.OnRegisterProjectMatch(ctx, projectRegisterRequest).Return(nil, nil)
	err := createProjectsCommand(ctx, args, cmdCtx)
	assert.Equal(t, fmt.Errorf("project name is required flag"), err)
	mockClient.AssertNotCalled(t, "RegisterProject", ctx, mock.Anything)
}
