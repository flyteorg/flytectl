package create

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"
	"testing"

	"github.com/lyft/flytectl/cmd/config"
	cmdCore "github.com/lyft/flytectl/cmd/core"
	"github.com/lyft/flyteidl/clients/go/admin/mocks"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/stretchr/testify/assert"
)

const projectValue = "dummyProject"

var (
	reader               *os.File
	writer               *os.File
	err                  error
	ctx                  context.Context
	mockClient           *mocks.AdminServiceClient
	mockOutStream        io.Writer
	args                 []string
	cmdCtx               cmdCore.CommandContext
	projectCreateRequest *admin.Project
	stdOut               *os.File
	stderr               *os.File
)

func setup() {
	reader, writer, err = os.Pipe()
	if err != nil {
		panic(err)
	}
	stdOut = os.Stdout
	stderr = os.Stderr
	os.Stdout = writer
	os.Stderr = writer
	log.SetOutput(writer)
	config.GetConfig().Project = projectValue
	mockClient = new(mocks.AdminServiceClient)
	mockOutStream = writer
	cmdCtx = cmdCore.NewCommandContext(mockClient, mockOutStream)
	projectCreateRequest = &admin.Project{
		Id:          projectValue,
		Description: "Testing",
		Name:        projectValue,
	}
}

func teardownAndVerify(t *testing.T, expectedLog string) {
	writer.Close()
	os.Stdout = stdOut
	os.Stderr = stderr
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, reader)
	assert.Equal(t, expectedLog, buf.String())
}

func TestCreateProjectIdValidation(t *testing.T) {
	setup()
	defer teardownAndVerify(t, "Project ID is required flag")
	createProjectsCommand(ctx, args, cmdCtx)
	mockClient.AssertNotCalled(t, "CreateProject", ctx, projectCreateRequest)
}

func TestCreateProjectNameValidation(t *testing.T) {
	setup()
	defer teardownAndVerify(t, "Project name is required flag")
	projectConfig.name = ""
	projectConfig.id = "test"
	createProjectsCommand(ctx, args, cmdCtx)
	mockClient.AssertNotCalled(t, "CreateProject", ctx, projectCreateRequest)
}

