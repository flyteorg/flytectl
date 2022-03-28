package create

import (
	"context"
	"github.com/flyteorg/flyteidl/clients/go/admin"
	"sort"
	"testing"

	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flytectl/cmd/testutils"
	"github.com/stretchr/testify/assert"
)

const testDataFolder = "../testdata/"

var (
	err        error
	ctx        context.Context
	mockClient *admin.Clientset
	args       []string
	cmdCtx     cmdCore.CommandContext
)
var setup = testutils.Setup
var tearDownAndVerify = testutils.TearDownAndVerify

func TestCreateCommand(t *testing.T) {
	createCommand := RemoteCreateCommand()
	assert.Equal(t, createCommand.Use, "create")
	assert.Equal(t, createCommand.Short, "Creates various Flyte resources such as tasks, workflows, launch plans, executions, and projects.")
	assert.Equal(t, len(createCommand.Commands()), 2)
	cmdNouns := createCommand.Commands()
	// Sort by Use value.
	sort.Slice(cmdNouns, func(i, j int) bool {
		return cmdNouns[i].Use < cmdNouns[j].Use
	})
	assert.Equal(t, cmdNouns[0].Use, "execution")
	assert.Equal(t, cmdNouns[0].Aliases, []string{"executions"})
	assert.Equal(t, cmdNouns[0].Short, executionShort)
	assert.Equal(t, cmdNouns[1].Use, "project")
	assert.Equal(t, cmdNouns[1].Aliases, []string{"projects"})
	assert.Equal(t, cmdNouns[1].Short, "Creates project resources.")
}
