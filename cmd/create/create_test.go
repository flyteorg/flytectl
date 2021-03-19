package create

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCommand(t *testing.T) {
	createCommand := RemoteCreateCommand()
	assert.Equal(t, createCommand.Use, "create")
	assert.Equal(t, createCommand.Short, "Used for creating various flyte resources including tasks/workflows/launchplans/executions/project.")
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
	assert.Equal(t, cmdNouns[1].Short, "Create project resources")
}
