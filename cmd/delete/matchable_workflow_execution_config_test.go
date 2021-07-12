package delete

import (
	"fmt"
	"testing"

	"github.com/flyteorg/flytectl/cmd/config/subcommand/workflowexecutionconfig"

	"github.com/flyteorg/flytectl/cmd/config"
	u "github.com/flyteorg/flytectl/cmd/testutils"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func deleteWorkflowExecutionConfigSetup() {
	ctx = u.Ctx
	cmdCtx = u.CmdCtx
	mockClient = u.MockClient
	workflowexecutionconfig.DefaultDelConfig = &workflowexecutionconfig.AttrDeleteConfig{}
	args = []string{}
}

func TestDeleteWorkflowExecutionConfig(t *testing.T) {
	t.Run("successful project domain attribute deletion commandline", func(t *testing.T) {
		setup()
		deleteWorkflowExecutionConfigSetup()
		// Empty attribute file
		workflowexecutionconfig.DefaultDelConfig.AttrFile = ""
		// No args implying project domain attribute deletion
		u.DeleterExt.OnDeleteProjectDomainAttributesMatch(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything).Return(nil)
		err = deleteWorkflowExecutionConfig(ctx, args, cmdCtx)
		assert.Nil(t, err)
		u.DeleterExt.AssertCalled(t, "DeleteProjectDomainAttributes",
			ctx, config.GetConfig().Project, config.GetConfig().Domain, admin.MatchableResource_WORKFLOW_EXECUTION_CONFIG)
	})
	t.Run("failed project domain attribute deletion", func(t *testing.T) {
		setup()
		deleteWorkflowExecutionConfigSetup()
		// No args implying project domain attribute deletion
		u.DeleterExt.OnDeleteProjectDomainAttributesMatch(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything).Return(fmt.Errorf("failed to delte project domain attributes"))
		err = deleteWorkflowExecutionConfig(ctx, args, cmdCtx)
		assert.NotNil(t, err)
		assert.Equal(t, fmt.Errorf("failed to delte project domain attributes"), err)
		u.DeleterExt.AssertCalled(t, "DeleteProjectDomainAttributes",
			ctx, config.GetConfig().Project, config.GetConfig().Domain, admin.MatchableResource_WORKFLOW_EXECUTION_CONFIG)
	})
	t.Run("successful project domain attribute deletion file", func(t *testing.T) {
		setup()
		deleteWorkflowExecutionConfigSetup()
		// Empty attribute file
		workflowexecutionconfig.DefaultDelConfig.AttrFile = "testdata/valid_project_domain_workflow_execution_config.yaml"
		// No args implying project domain attribute deletion
		u.DeleterExt.OnDeleteProjectDomainAttributesMatch(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything).Return(nil)
		err = deleteWorkflowExecutionConfig(ctx, args, cmdCtx)
		assert.Nil(t, err)
		u.DeleterExt.AssertCalled(t, "DeleteProjectDomainAttributes",
			ctx, "flytectldemo", "development", admin.MatchableResource_WORKFLOW_EXECUTION_CONFIG)
	})
	t.Run("successful workflow attribute deletion", func(t *testing.T) {
		setup()
		deleteWorkflowExecutionConfigSetup()
		// Empty attribute file
		workflowexecutionconfig.DefaultDelConfig.AttrFile = ""
		args := []string{"workflow1"}
		u.DeleterExt.OnDeleteWorkflowAttributesMatch(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Return(nil)
		err = deleteWorkflowExecutionConfig(ctx, args, cmdCtx)
		assert.Nil(t, err)
		u.DeleterExt.AssertCalled(t, "DeleteWorkflowAttributes",
			ctx, config.GetConfig().Project, config.GetConfig().Domain, "workflow1",
			admin.MatchableResource_WORKFLOW_EXECUTION_CONFIG)
	})
	t.Run("failed workflow attribute deletion", func(t *testing.T) {
		setup()
		deleteWorkflowExecutionConfigSetup()
		// Empty attribute file
		workflowexecutionconfig.DefaultDelConfig.AttrFile = ""
		args := []string{"workflow1"}
		u.DeleterExt.OnDeleteWorkflowAttributesMatch(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Return(fmt.Errorf("failed to delete workflow attribute"))
		err = deleteWorkflowExecutionConfig(ctx, args, cmdCtx)
		assert.NotNil(t, err)
		assert.Equal(t, fmt.Errorf("failed to delete workflow attribute"), err)
		u.DeleterExt.AssertCalled(t, "DeleteWorkflowAttributes",
			ctx, config.GetConfig().Project, config.GetConfig().Domain, "workflow1",
			admin.MatchableResource_WORKFLOW_EXECUTION_CONFIG)
	})
	t.Run("successful workflow attribute deletion file", func(t *testing.T) {
		setup()
		deleteWorkflowExecutionConfigSetup()
		// Empty attribute file
		workflowexecutionconfig.DefaultDelConfig.AttrFile = "testdata/valid_workflow_workflow_execution_config.yaml"
		// No args implying project domain attribute deletion
		u.DeleterExt.OnDeleteWorkflowAttributesMatch(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Return(nil)
		err = deleteWorkflowExecutionConfig(ctx, args, cmdCtx)
		assert.Nil(t, err)
		u.DeleterExt.AssertCalled(t, "DeleteWorkflowAttributes",
			ctx, "flytectldemo", "development", "core.control_flow.run_merge_sort.merge_sort",
			admin.MatchableResource_WORKFLOW_EXECUTION_CONFIG)
	})
	t.Run("workflow attribute deletion non existent file", func(t *testing.T) {
		setup()
		deleteWorkflowExecutionConfigSetup()
		// Empty attribute file
		workflowexecutionconfig.DefaultDelConfig.AttrFile = testDataNonExistentFile
		// No args implying project domain attribute deletion
		u.DeleterExt.OnDeleteWorkflowAttributesMatch(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Return(nil)
		err = deleteWorkflowExecutionConfig(ctx, args, cmdCtx)
		assert.NotNil(t, err)
		u.DeleterExt.AssertNotCalled(t, "DeleteWorkflowAttributes",
			ctx, "flytectldemo", "development", "core.control_flow.run_merge_sort.merge_sort",
			admin.MatchableResource_WORKFLOW_EXECUTION_CONFIG)
	})
	t.Run("attribute deletion invalid file", func(t *testing.T) {
		setup()
		deleteWorkflowExecutionConfigSetup()
		// Empty attribute file
		workflowexecutionconfig.DefaultDelConfig.AttrFile = testDataInvalidAttrFile
		// No args implying project domain attribute deletion
		err = deleteWorkflowExecutionConfig(ctx, args, cmdCtx)
		assert.NotNil(t, err)
		assert.Equal(t,
			fmt.Errorf("error unmarshaling JSON: while decoding JSON: json: unknown field \"InvalidDomain\""),
			err)
		u.DeleterExt.AssertNotCalled(t, "DeleteProjectDomainAttributes",
			ctx, "flytectldemo", "development", admin.MatchableResource_WORKFLOW_EXECUTION_CONFIG)
	})
}
