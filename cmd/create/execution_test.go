package create

import (
	"fmt"
	"testing"

	"github.com/flyteorg/flytectl/cmd/config"
	cmdGet "github.com/flyteorg/flytectl/cmd/get"
	"github.com/flyteorg/flytectl/cmd/get/interfaces/mocks"
	"github.com/flyteorg/flytectl/cmd/testutils"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// This function needs to be called after testutils.Steup()
func createExecutionSetup() {
	ctx = testutils.Ctx
	cmdCtx = testutils.CmdCtx
	mockClient = testutils.MockClient
	sortedListLiteralType := core.Variable{
		Type: &core.LiteralType{
			Type: &core.LiteralType_CollectionType{
				CollectionType: &core.LiteralType{
					Type: &core.LiteralType_Simple{
						Simple: core.SimpleType_INTEGER,
					},
				},
			},
		},
	}
	variableMap := map[string]*core.Variable{
		"sorted_list1": &sortedListLiteralType,
		"sorted_list2": &sortedListLiteralType,
	}

	task1 := &admin.Task{
		Id: &core.Identifier{
			Name:    "task1",
			Version: "v2",
		},
		Closure: &admin.TaskClosure{
			CreatedAt: &timestamppb.Timestamp{Seconds: 1, Nanos: 0},
			CompiledTask: &core.CompiledTask{
				Template: &core.TaskTemplate{
					Interface: &core.TypedInterface{
						Inputs: &core.VariableMap{
							Variables: variableMap,
						},
					},
				},
			},
		},
	}
	mockClient.OnGetTaskMatch(ctx, mock.Anything).Return(task1, nil)
	parameterMap := map[string]*core.Parameter{
		"numbers": {
			Var: &core.Variable{
				Type: &core.LiteralType{
					Type: &core.LiteralType_CollectionType{
						CollectionType: &core.LiteralType{
							Type: &core.LiteralType_Simple{
								Simple: core.SimpleType_INTEGER,
							},
						},
					},
				},
			},
		},
		"numbers_count": {
			Var: &core.Variable{
				Type: &core.LiteralType{
					Type: &core.LiteralType_Simple{
						Simple: core.SimpleType_INTEGER,
					},
				},
			},
		},
		"run_local_at_count": {
			Var: &core.Variable{
				Type: &core.LiteralType{
					Type: &core.LiteralType_Simple{
						Simple: core.SimpleType_INTEGER,
					},
				},
			},
			Behavior: &core.Parameter_Default{
				Default: &core.Literal{
					Value: &core.Literal_Scalar{
						Scalar: &core.Scalar{
							Value: &core.Scalar_Primitive{
								Primitive: &core.Primitive{
									Value: &core.Primitive_Integer{
										Integer: 10,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	launchPlan1 := &admin.LaunchPlan{
		Id: &core.Identifier{
			Name:    "core.advanced.run_merge_sort.merge_sort",
			Version: "v3",
		},
		Spec: &admin.LaunchPlanSpec{
			DefaultInputs: &core.ParameterMap{
				Parameters: parameterMap,
			},
		},
		Closure: &admin.LaunchPlanClosure{
			CreatedAt: &timestamppb.Timestamp{Seconds: 0, Nanos: 0},
			ExpectedInputs: &core.ParameterMap{
				Parameters: parameterMap,
			},
		},
	}
	objectGetRequest := &admin.ObjectGetRequest{
		Id: &core.Identifier{
			ResourceType: core.ResourceType_LAUNCH_PLAN,
			Project:      config.GetConfig().Project,
			Domain:       config.GetConfig().Domain,
			Name:         "core.advanced.run_merge_sort.merge_sort",
			Version:      "v3",
		},
	}
	mockClient.OnGetLaunchPlanMatch(ctx, objectGetRequest).Return(launchPlan1, nil)
}

func TestCreateTaskExecutionFunc(t *testing.T) {
	setup()
	createExecutionSetup()
	executionCreateResponseTask := &admin.ExecutionCreateResponse{
		Id: &core.WorkflowExecutionIdentifier{
			Project: "flytesnacks",
			Domain:  "development",
			Name:    "ff513c0e44b5b4a35aa5",
		},
	}
	mockClient.OnCreateExecutionMatch(ctx, mock.Anything).Return(executionCreateResponseTask, nil)
	executionConfig.ExecFile = testDataFolder + "task_execution_spec.yaml"
	err = createExecutionCommand(ctx, args, cmdCtx)
	assert.Nil(t, err)
	mockClient.AssertCalled(t, "CreateExecution", ctx, mock.Anything)
	tearDownAndVerify(t, `execution identifier project:"flytesnacks" domain:"development" name:"ff513c0e44b5b4a35aa5" `)
}

func TestCreateTaskExecutionFuncError(t *testing.T) {
	setup()
	createExecutionSetup()
	mockClient.OnCreateExecutionMatch(ctx, mock.Anything).Return(nil, fmt.Errorf("error launching task"))
	executionConfig.ExecFile = testDataFolder + "task_execution_spec.yaml"
	err = createExecutionCommand(ctx, args, cmdCtx)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("error launching task"), err)
	mockClient.AssertCalled(t, "CreateExecution", ctx, mock.Anything)
}

func TestCreateLaunchPlanExecutionFunc(t *testing.T) {
	setup()
	createExecutionSetup()
	executionCreateResponseLP := &admin.ExecutionCreateResponse{
		Id: &core.WorkflowExecutionIdentifier{
			Project: "flytesnacks",
			Domain:  "development",
			Name:    "f652ea3596e7f4d80a0e",
		},
	}
	mockClient.OnCreateExecutionMatch(ctx, mock.Anything).Return(executionCreateResponseLP, nil)
	executionConfig.ExecFile = testDataFolder + "launchplan_execution_spec.yaml"
	err = createExecutionCommand(ctx, args, cmdCtx)
	assert.Nil(t, err)
	mockClient.AssertCalled(t, "CreateExecution", ctx, mock.Anything)
	tearDownAndVerify(t, `execution identifier project:"flytesnacks" domain:"development" name:"f652ea3596e7f4d80a0e"`)
}

func TestCreateRelaunchExecutionFunc(t *testing.T) {
	setup()
	createExecutionSetup()
	executionCreateResponseLP := &admin.ExecutionCreateResponse{
		Id: &core.WorkflowExecutionIdentifier{
			Project: "flytesnacks",
			Domain:  "development",
			Name:    "f652ea3596e7f4d80a0e",
		},
	}
	literalMap := &core.LiteralMap{
		Literals: nil,
	}
	exec = &admin.Execution{
		Id: &core.WorkflowExecutionIdentifier{
			Project: config.GetConfig().Project,
			Domain:  config.GetConfig().Domain,
			Name:    "ffb31066a0f8b4d52b77",
		},
		Spec: &admin.ExecutionSpec{
			LaunchPlan: &core.Identifier{
				Name:    "core.advanced.run_merge_sort.merge_sort",
				Version: "v3",
			},
			Inputs: literalMap,
		},
	}
	launchPlan = &admin.LaunchPlan{
		Id: &core.Identifier{
			Name:    "core.advanced.run_merge_sort.merge_sort",
			Version: "v3",
		},
	}
	mockFetcher = &mocks.Fetcher{}
	cmdGet.DefaultFetcher = mockFetcher
	mockFetcher.OnFetchExecutionMatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(exec, nil)
	mockFetcher.OnFetchLPVersionMatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(launchPlan, nil)
	mockClient.OnCreateExecutionMatch(ctx, mock.Anything).Return(executionCreateResponseLP, nil)
	executionConfig.Relaunch = "xb5317xbty"
	err = createExecutionCommand(ctx, args, cmdCtx)
	assert.Nil(t, err)
	mockClient.AssertCalled(t, "CreateExecution", ctx, mock.Anything)
	tearDownAndVerify(t, `execution identifier project:"flytesnacks" domain:"development" name:"f652ea3596e7f4d80a0e"`)
}
