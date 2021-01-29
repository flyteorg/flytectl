package update

import (
	"context"
	"errors"
	"github.com/lyft/flytectl/cmd/config"
	cmdCore "github.com/lyft/flytectl/cmd/core"
	"github.com/lyft/flyteidl/clients/go/admin/mocks"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

const projectValue = "dummyProject"

func modifyProjectFlags(archiveProject *bool, newArchiveVal bool, activateProject *bool, newActivateVal bool) {
	*archiveProject = newArchiveVal
	*activateProject = newActivateVal
}

func TestActivateProjectFunc(t *testing.T) {
	ctx := context.Background()
	config.GetConfig().Project = projectValue
	modifyProjectFlags(&(GetConfig().ArchiveProject), false, &(GetConfig().ActivateProject), true)
	var args []string
	mockClient := new(mocks.AdminServiceClient)
	mockOutStream := new(io.Writer)
	cmdCtx := cmdCore.NewCommandContext(mockClient, *mockOutStream)
	projectUpdateRequest := &admin.Project{
		Id : projectValue,
		State: admin.Project_ACTIVE,
	}
	mockClient.OnUpdateProjectMatch(ctx, projectUpdateRequest).Return(nil, nil)
	err := updateProjectsFunc(ctx, args, cmdCtx)
	assert.Nil(t, err)
	mockClient.AssertCalled(t, "UpdateProject", ctx, projectUpdateRequest)
}

func TestActivateProjectFuncWithError(t *testing.T) {
	ctx := context.Background()
	config.GetConfig().Project = projectValue
	modifyProjectFlags(&(GetConfig().ArchiveProject), false, &(GetConfig().ActivateProject), true)
	var args []string
	mockClient := new(mocks.AdminServiceClient)
	mockOutStream := new(io.Writer)
	cmdCtx := cmdCore.NewCommandContext(mockClient, *mockOutStream)
	projectUpdateRequest := &admin.Project{
		Id : projectValue,
		State: admin.Project_ACTIVE,
	}
	mockClient.OnUpdateProjectMatch(ctx, projectUpdateRequest).Return(nil, errors.New("Error Updating Project"))
	err := updateProjectsFunc(ctx, args, cmdCtx)
	assert.Equal(t, err, errors.New("Error Updating Project"))
	mockClient.AssertCalled(t, "UpdateProject", ctx, projectUpdateRequest)
}

func TestArchiveProjectFunc(t *testing.T) {
	ctx := context.Background()
	config.GetConfig().Project = projectValue
	modifyProjectFlags(&(GetConfig().ArchiveProject), true, &(GetConfig().ActivateProject), false)
	var args []string
	mockClient := new(mocks.AdminServiceClient)
	mockOutStream := new(io.Writer)
	cmdCtx := cmdCore.NewCommandContext(mockClient, *mockOutStream)
	projectUpdateRequest := &admin.Project{
		Id : projectValue,
		State: admin.Project_ARCHIVED,
	}
	mockClient.OnUpdateProjectMatch(ctx, projectUpdateRequest).Return(nil, nil)
	err := updateProjectsFunc(ctx, args, cmdCtx)
	assert.Nil(t, err)
	mockClient.AssertCalled(t, "UpdateProject", ctx, projectUpdateRequest)
}

func TestArchiveProjectFuncWithError(t *testing.T) {
	ctx := context.Background()
	config.GetConfig().Project = projectValue
	modifyProjectFlags(&(GetConfig().ArchiveProject), true, &(GetConfig().ActivateProject), false)
	var args []string
	mockClient := new(mocks.AdminServiceClient)
	mockOutStream := new(io.Writer)
	cmdCtx := cmdCore.NewCommandContext(mockClient, *mockOutStream)
	projectUpdateRequest := &admin.Project{
		Id : projectValue,
		State: admin.Project_ARCHIVED,
	}
	mockClient.OnUpdateProjectMatch(ctx, projectUpdateRequest).Return(nil, errors.New("Error Updating Project"))
	err := updateProjectsFunc(ctx, args, cmdCtx)
	assert.Equal(t, err, errors.New("Error Updating Project"))
	mockClient.AssertCalled(t, "UpdateProject", ctx, projectUpdateRequest)
}

func TestEmptyProjectInput(t *testing.T) {
	ctx := context.Background()
	config.GetConfig().Project = ""
	modifyProjectFlags(&(GetConfig().ArchiveProject), false, &(GetConfig().ActivateProject), true)
	var args []string
	mockClient := new(mocks.AdminServiceClient)
	mockOutStream := new(io.Writer)
	cmdCtx := cmdCore.NewCommandContext(mockClient, *mockOutStream)
	projectUpdateRequest := &admin.Project{
		Id : projectValue,
		State: admin.Project_ACTIVE,
	}
	mockClient.OnUpdateProjectMatch(ctx, projectUpdateRequest).Return(nil, nil)
	err := updateProjectsFunc(ctx, args, cmdCtx)
	assert.NotNil(t, err)
	assert.Equal(t, err, errors.New("Specify id of the project to be updated"))
	mockClient.AssertNotCalled(t, "UpdateProject", ctx, projectUpdateRequest)
}

func TestInvalidInput(t *testing.T) {
	ctx := context.Background()
	config.GetConfig().Project = "flytesnacks"
	modifyProjectFlags(&(GetConfig().ArchiveProject), false, &(GetConfig().ActivateProject), false)
	var args []string
	mockClient := new(mocks.AdminServiceClient)
	mockOutStream := new(io.Writer)
	cmdCtx := cmdCore.NewCommandContext(mockClient, *mockOutStream)
	projectUpdateRequest := &admin.Project{
		Id : projectValue,
		State: admin.Project_ACTIVE,
	}
	mockClient.OnUpdateProjectMatch(ctx, projectUpdateRequest).Return(nil, nil)
	err := updateProjectsFunc(ctx, args, cmdCtx)
	assert.NotNil(t, err)
	assert.Equal(t, err, errors.New("Specify either activate or archive"))
	mockClient.AssertNotCalled(t, "UpdateProject", ctx, projectUpdateRequest)
}