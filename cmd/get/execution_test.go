package get

import (
	"context"
	cmdCore "github.com/lyft/flytectl/cmd/core"
	"github.com/lyft/flyteidl/clients/go/admin/mocks"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/core"
	"github.com/stretchr/testify/assert"
	"testing"
)


const projectValue = "dummyProject"
const domainValue = "dummyDomain"
const executionNameValue = "e124"
const launchPlanNameValue = "lp_name"
const launchPlanVersionValue = "lp_version"
const workflowNameValue = "wf_name"
const workflowVersionValue = "wf_version"


func TestGetExecutionFunc(t *testing.T) {
	var ctx context.Context
	var args []string
	cmdCtx  := cmdCore.CommandContext{}
	mockClient := new(mocks.AdminServiceClient)
	execGetRequest := &admin.WorkflowExecutionGetRequest{
		Id: &core.WorkflowExecutionIdentifier{
			//Project: projectValue,
			//Domain:  domainValue,
			Name:    executionNameValue,
		},
	}
	execListRequest := &admin.ResourceListRequest{
		Limit: 100,
		Id: &admin.NamedEntityIdentifier{
			//Project: projectValue,
			//Domain:  domainValue,
		},
	}
	executionResponse := &admin.Execution{
		Id: &core.WorkflowExecutionIdentifier{
			Project: projectValue,
			Domain:  domainValue,
			Name:    executionNameValue,
		},
		Spec: &admin.ExecutionSpec{
			LaunchPlan: &core.Identifier{
				Project: projectValue,
				Domain:  domainValue,
				Name:    launchPlanNameValue,
				Version: launchPlanVersionValue,
			},
		},
		Closure: &admin.ExecutionClosure{
			WorkflowId: &core.Identifier{
				Project: projectValue,
				Domain:  domainValue,
				Name:    workflowNameValue,
				Version: workflowVersionValue,
			},
			Phase: core.WorkflowExecution_SUCCEEDED,
		},
	}
	var executions []* admin.Execution
	executions = append(executions, executionResponse)
	executionList := &admin.ExecutionList{
		Executions: executions,
	}
	//var callOptions []grpc.CallOption
	cmdCtx.SetAdminClient(mockClient)
	mockClient.OnGetExecution(ctx, execGetRequest).Return(executionResponse, nil)
	mockClient.OnListExecutions(ctx, execListRequest).Return(executionList, nil)
	//mockClient.OnListExecutions(ctx, execListRequest, nil).Return(executionList, nil)
	//mockClient.OnListExecutions(ctx, execListRequest, callOptions...).Return(executionList, nil)
	err := getExecutionFunc(ctx, args, cmdCtx)
	assert.Nil(t,err)
}

