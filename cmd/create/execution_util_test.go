package create

import (
	"errors"
	"fmt"
	"testing"

	"github.com/flyteorg/flytectl/cmd/config"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	executionCreateResponse *admin.ExecutionCreateResponse
	relaunchRequest         *admin.ExecutionRelaunchRequest
	recoverRequest          *admin.ExecutionRecoverRequest
)

// This function needs to be called after testutils.Steup()
func createExecutionUtilSetup() {
	executionCreateResponse = &admin.ExecutionCreateResponse{
		Id: &core.WorkflowExecutionIdentifier{
			Project: "flytesnacks",
			Domain:  "development",
			Name:    "f652ea3596e7f4d80a0e",
		},
	}
	relaunchRequest = &admin.ExecutionRelaunchRequest{
		Id: &core.WorkflowExecutionIdentifier{
			Name:    "execName",
			Project: config.GetConfig().Project,
			Domain:  config.GetConfig().Domain,
		},
	}
	recoverRequest = &admin.ExecutionRecoverRequest{
		Id: &core.WorkflowExecutionIdentifier{
			Name:    "execName",
			Project: config.GetConfig().Project,
			Domain:  config.GetConfig().Domain,
		},
	}
	executionConfig = &ExecutionConfig{}
}

func TestCreateExecutionForRelaunch(t *testing.T) {
	s := setup()
	createExecutionUtilSetup()
	s.MockAdminClient.OnRelaunchExecutionMatch(s.Ctx, relaunchRequest).Return(executionCreateResponse, nil)
	err := relaunchExecution(s.Ctx, "execName", config.GetConfig().Project, config.GetConfig().Domain, s.CmdCtx, executionConfig)
	assert.Nil(t, err)
}

func TestCreateExecutionForRelaunchNotFound(t *testing.T) {
	s := setup()
	createExecutionUtilSetup()
	s.MockAdminClient.OnRelaunchExecutionMatch(s.Ctx, relaunchRequest).Return(nil, errors.New("unknown execution"))
	err := relaunchExecution(s.Ctx, "execName", config.GetConfig().Project, config.GetConfig().Domain, s.CmdCtx, executionConfig)

	assert.NotNil(t, err)
	assert.Equal(t, err, errors.New("unknown execution"))
}

func TestCreateExecutionForRecovery(t *testing.T) {
	s := setup()
	createExecutionUtilSetup()
	s.MockAdminClient.OnRecoverExecutionMatch(s.Ctx, recoverRequest).Return(executionCreateResponse, nil)
	err := recoverExecution(s.Ctx, "execName", config.GetConfig().Project, config.GetConfig().Domain, s.CmdCtx, executionConfig)
	assert.Nil(t, err)
}

func TestCreateExecutionForRecoveryNotFound(t *testing.T) {
	s := setup()
	createExecutionUtilSetup()
	s.MockAdminClient.OnRecoverExecutionMatch(s.Ctx, recoverRequest).Return(nil, errors.New("unknown execution"))
	err := recoverExecution(s.Ctx, "execName", config.GetConfig().Project, config.GetConfig().Domain, s.CmdCtx, executionConfig)
	assert.NotNil(t, err)
	assert.Equal(t, err, errors.New("unknown execution"))
}

func TestCreateExecutionRequestForWorkflow(t *testing.T) {
	t.Run("successful", func(t *testing.T) {
		setup()
		createExecutionUtilSetup()
		launchPlan := &admin.LaunchPlan{}
		fetcherClient.OnFetchLPVersionMatch(ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(launchPlan, nil)
		execCreateRequest, err := createExecutionRequestForWorkflow(ctx, "wfName", config.GetConfig().Project, config.GetConfig().Domain, cmdCtx, executionConfig)
		assert.Nil(t, err)
		assert.NotNil(t, execCreateRequest)
	})
	t.Run("failed literal conversion", func(t *testing.T) {
		setup()
		createExecutionUtilSetup()
		launchPlan := &admin.LaunchPlan{
			Spec: &admin.LaunchPlanSpec{
				DefaultInputs: &core.ParameterMap{
					Parameters: map[string]*core.Parameter{"nilparam": nil},
				},
			},
		}
		fetcherClient.OnFetchLPVersionMatch(ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(launchPlan, nil)
		execCreateRequest, err := createExecutionRequestForWorkflow(ctx, "wfName", config.GetConfig().Project, config.GetConfig().Domain, cmdCtx, executionConfig)
		assert.NotNil(t, err)
		assert.Nil(t, execCreateRequest)
		assert.Equal(t, fmt.Errorf("parameter [nilparam] has nil Variable"), err)
	})
	t.Run("failed fetch", func(t *testing.T) {
		setup()
		createExecutionUtilSetup()
		fetcherClient.OnFetchLPVersionMatch(ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("failed"))
		execCreateRequest, err := createExecutionRequestForWorkflow(ctx, "wfName", config.GetConfig().Project, config.GetConfig().Domain, cmdCtx, executionConfig)
		assert.NotNil(t, err)
		assert.Nil(t, execCreateRequest)
		assert.Equal(t, err, errors.New("failed"))
	})
	t.Run("with security context", func(t *testing.T) {
		setup()
		createExecutionUtilSetup()
		executionConfig.KubeServiceAcct = "default"
		launchPlan := &admin.LaunchPlan{}
		fetcherClient.OnFetchLPVersionMatch(ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(launchPlan, nil)
		mockClient.OnGetLaunchPlanMatch(ctx, mock.Anything).Return(launchPlan, nil)
		execCreateRequest, err := createExecutionRequestForWorkflow(ctx, "wfName", config.GetConfig().Project, config.GetConfig().Domain, cmdCtx, executionConfig)
		assert.Nil(t, err)
		assert.NotNil(t, execCreateRequest)
		executionConfig.KubeServiceAcct = ""
	})
}

func TestCreateExecutionRequestForTask(t *testing.T) {
	t.Run("successful", func(t *testing.T) {
		setup()
		createExecutionUtilSetup()
		task := &admin.Task{
			Id: &core.Identifier{
				Name: "taskName",
			},
		}
		fetcherClient.OnFetchTaskVersionMatch(ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(task, nil)
		execCreateRequest, err := createExecutionRequestForTask(ctx, "taskName", config.GetConfig().Project, config.GetConfig().Domain, cmdCtx, executionConfig)
		assert.Nil(t, err)
		assert.NotNil(t, execCreateRequest)
	})
	t.Run("failed literal conversion", func(t *testing.T) {
		setup()
		createExecutionUtilSetup()
		task := &admin.Task{
			Closure: &admin.TaskClosure{
				CompiledTask: &core.CompiledTask{
					Template: &core.TaskTemplate{
						Interface: &core.TypedInterface{
							Inputs: &core.VariableMap{
								Variables: map[string]*core.Variable{
									"nilvar": nil,
								},
							},
						},
					},
				},
			},
		}
		fetcherClient.OnFetchTaskVersionMatch(ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(task, nil)
		execCreateRequest, err := createExecutionRequestForTask(ctx, "taskName", config.GetConfig().Project, config.GetConfig().Domain, cmdCtx, executionConfig)
		assert.NotNil(t, err)
		assert.Nil(t, execCreateRequest)
		assert.Equal(t, fmt.Errorf("variable [nilvar] has nil type"), err)
	})
	t.Run("failed fetch", func(t *testing.T) {
		setup()
		createExecutionUtilSetup()
		fetcherClient.OnFetchTaskVersionMatch(ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("failed"))
		execCreateRequest, err := createExecutionRequestForTask(ctx, "taskName", config.GetConfig().Project, config.GetConfig().Domain, cmdCtx, executionConfig)
		assert.NotNil(t, err)
		assert.Nil(t, execCreateRequest)
		assert.Equal(t, err, errors.New("failed"))
	})
	t.Run("with security context", func(t *testing.T) {
		setup()
		createExecutionUtilSetup()
		executionConfig.KubeServiceAcct = "default"
		task := &admin.Task{
			Id: &core.Identifier{
				Name: "taskName",
			},
		}
		fetcherClient.OnFetchTaskVersionMatch(ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(task, nil)
		execCreateRequest, err := createExecutionRequestForTask(ctx, "taskName", config.GetConfig().Project, config.GetConfig().Domain, cmdCtx, executionConfig)
		assert.Nil(t, err)
		assert.NotNil(t, execCreateRequest)
		executionConfig.KubeServiceAcct = ""
	})
}
