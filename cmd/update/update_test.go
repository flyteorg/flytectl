package update

import (
	"context"
	"sort"
	"testing"

	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flytectl/cmd/testutils"
	"github.com/flyteorg/flyteidl/clients/go/admin/mocks"

	"github.com/stretchr/testify/assert"
)

var (
	err        error
	ctx        context.Context
	mockClient *mocks.AdminServiceClient
	cmdCtx     cmdCore.CommandContext
)

const (
	testDataNonExistentFile = "testdata/non-existent-file"
	testDataInvalidAttrFile = "testdata/invalid_attribute.yaml"
)

var setup = testutils.Setup
var tearDownAndVerify = testutils.TearDownAndVerify

func TestUpdateCommand(t *testing.T) {
	updateCommand := CreateUpdateCommand()
	assert.Equal(t, updateCommand.Use, updateUse)
	assert.Equal(t, updateCommand.Short, updateShort)
	assert.Equal(t, updateCommand.Long, updatecmdLong)
	assert.Equal(t, len(updateCommand.Commands()), 10)
	cmdNouns := updateCommand.Commands()
	// Sort by Use value.
	sort.Slice(cmdNouns, func(i, j int) bool {
		return cmdNouns[i].Use < cmdNouns[j].Use
	})
	useArray := []string{"cluster-resource-attribute", "execution-cluster-label", "execution-queue-attribute", "launchplan",
		"plugin-override", "project", "task", "task-resource-attribute", "workflow", "workflow-execution-config"}
	aliases := [][]string{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}}
	shortArray := []string{clusterResourceAttributesShort, executionClusterLabelShort, executionQueueAttributesShort, updateLPShort,
		pluginOverrideShort, projectShort, updateTaskShort, taskResourceAttributesShort, updateWorkflowShort, workflowExecutionConfigShort}
	longArray := []string{clusterResourceAttributesLong, executionClusterLabelLong, executionQueueAttributesLong, updateLPLong,
		pluginOverrideLong, projectLong, updateTaskLong, taskResourceAttributesLong, updateWorkflowLong, workflowExecutionConfigLong}
	for i := range cmdNouns {
		assert.Equal(t, cmdNouns[i].Use, useArray[i])
		assert.Equal(t, cmdNouns[i].Aliases, aliases[i])
		assert.Equal(t, cmdNouns[i].Short, shortArray[i])
		assert.Equal(t, cmdNouns[i].Long, longArray[i])
	}
}
