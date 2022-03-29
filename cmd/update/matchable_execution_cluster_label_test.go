package update

import (
	"fmt"
	"testing"

	"github.com/flyteorg/flytectl/cmd/config/subcommand/executionclusterlabel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func updateExecutionClusterLabelSetup() {
	executionclusterlabel.DefaultUpdateConfig = &executionclusterlabel.AttrUpdateConfig{}
}

func TestExecutionClusterLabel(t *testing.T) {
	t.Run("no input file for update", func(t *testing.T) {
		s := setup()
		updateExecutionClusterLabelSetup()
		err := updateExecutionClusterLabelFunc(s.Ctx, []string{}, s.CmdCtx)
		assert.NotNil(t, err)
		assert.Equal(t, fmt.Errorf("attrFile is mandatory while calling update for execution cluster label"), err)
		tearDownAndVerify(t, s.Reader, ``)
	})
	t.Run("successful update project domain attribute", func(t *testing.T) {
		s := setup()
		updateExecutionClusterLabelSetup()
		executionclusterlabel.DefaultUpdateConfig.AttrFile = "testdata/valid_project_domain_execution_cluster_label.yaml"
		// No args implying project domain attribute deletion
		s.UpdaterExt.OnUpdateProjectDomainAttributesMatch(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything).Return(nil)
		err := updateExecutionClusterLabelFunc(s.Ctx, []string{}, s.CmdCtx)
		assert.Nil(t, err)
		tearDownAndVerify(t, s.Reader, `Updated attributes from flytectldemo project and domain development`)
	})
	t.Run("failed update project domain attribute", func(t *testing.T) {
		s := setup()
		updateExecutionClusterLabelSetup()
		executionclusterlabel.DefaultUpdateConfig.AttrFile = "testdata/valid_project_domain_execution_cluster_label.yaml"
		// No args implying project domain attribute deletion
		s.UpdaterExt.OnUpdateProjectDomainAttributesMatch(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything).Return(fmt.Errorf("failed to update attributes"))
		err := updateExecutionClusterLabelFunc(s.Ctx, []string{}, s.CmdCtx)
		assert.NotNil(t, err)
		assert.Equal(t, fmt.Errorf("failed to update attributes"), err)
		tearDownAndVerify(t, s.Reader, ``)
	})
	t.Run("successful update workflow attribute", func(t *testing.T) {
		s := setup()
		updateExecutionClusterLabelSetup()
		executionclusterlabel.DefaultUpdateConfig.AttrFile = "testdata/valid_workflow_execution_cluster_label.yaml"
		// No args implying project domain attribute deletion
		s.UpdaterExt.OnUpdateWorkflowAttributesMatch(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Return(nil)
		err := updateExecutionClusterLabelFunc(s.Ctx, nil, s.CmdCtx)
		assert.Nil(t, err)
		tearDownAndVerify(t, s.Reader, `Updated attributes from flytectldemo project and domain development and workflow core.control_flow.run_merge_sort.merge_sort`)
	})
	t.Run("failed update workflow attribute", func(t *testing.T) {
		s := setup()
		updateExecutionClusterLabelSetup()
		executionclusterlabel.DefaultUpdateConfig.AttrFile = "testdata/valid_workflow_execution_cluster_label.yaml"
		// No args implying project domain attribute deletion
		s.UpdaterExt.OnUpdateWorkflowAttributesMatch(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Return(fmt.Errorf("failed to update attributes"))
		err := updateExecutionClusterLabelFunc(s.Ctx, nil, s.CmdCtx)
		assert.NotNil(t, err)
		assert.Equal(t, fmt.Errorf("failed to update attributes"), err)
		tearDownAndVerify(t, s.Reader, ``)
	})
	t.Run("non existent file", func(t *testing.T) {
		s := setup()
		updateExecutionClusterLabelSetup()
		executionclusterlabel.DefaultUpdateConfig.AttrFile = testDataNonExistentFile
		err := updateExecutionClusterLabelFunc(s.Ctx, nil, s.CmdCtx)
		assert.NotNil(t, err)
		assert.Equal(t, fmt.Errorf("unable to read from testdata/non-existent-file yaml file"), err)
		tearDownAndVerify(t, s.Reader, ``)
	})
	t.Run("invalid update file", func(t *testing.T) {
		s := setup()
		updateExecutionClusterLabelSetup()
		executionclusterlabel.DefaultUpdateConfig.AttrFile = testDataInvalidAttrFile
		err := updateExecutionClusterLabelFunc(s.Ctx, nil, s.CmdCtx)
		assert.NotNil(t, err)
		assert.Equal(t, fmt.Errorf("error unmarshaling JSON: while decoding JSON: json: unknown field \"InvalidDomain\""), err)
		tearDownAndVerify(t, s.Reader, ``)
	})
}
