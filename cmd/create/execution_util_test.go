package create

import (
	"errors"
	"testing"

	"github.com/flyteorg/flytectl/cmd/config"
	cmdGet "github.com/flyteorg/flytectl/cmd/get"
	"github.com/flyteorg/flytectl/cmd/get/interfaces/mocks"
	"github.com/flyteorg/flytectl/cmd/testutils"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	mockFetcher *mocks.Fetcher
	exec        *admin.Execution
	launchPlan  *admin.LaunchPlan
)

// This function needs to be called after testutils.Steup()
func createExecutionUtilSetup() {
	ctx = testutils.Ctx
	cmdCtx = testutils.CmdCtx
	mockClient = testutils.MockClient
	mockFetcher = &mocks.Fetcher{}
	cmdGet.DefaultFetcher = mockFetcher
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
}

func TestCreateExecutionForRelaunch(t *testing.T) {
	setup()
	createExecutionUtilSetup()
	mockFetcher.OnFetchExecutionMatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(exec, nil)
	mockFetcher.OnFetchLPVersionMatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(launchPlan, nil)
	var actualExecCreateRequest *admin.ExecutionCreateRequest
	actualExecCreateRequest, err = createExecutionRequestForRelaunch(ctx, "execName", projectValue, "domainValue", cmdCtx)
	assert.Nil(t, err)
	assert.NotNil(t, actualExecCreateRequest)
}

func TestCreateExecutionForRelaunchNotFound(t *testing.T) {
	setup()
	createExecutionUtilSetup()
	mockFetcher.OnFetchExecutionMatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("unknown execution"))
	_, err = createExecutionRequestForRelaunch(ctx, "execName", projectValue, "domainValue", cmdCtx)
	assert.NotNil(t, err)
	assert.Equal(t, err, errors.New("unknown execution"))
}

func TestCreateExecutionForRelaunchLaunchPlanError(t *testing.T) {
	setup()
	createExecutionUtilSetup()
	mockFetcher.OnFetchExecutionMatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(exec, nil)
	mockFetcher.OnFetchLPVersionMatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("unknown launchplan"))
	_, err = createExecutionRequestForRelaunch(ctx, "execName", projectValue, "domainValue", cmdCtx)
	assert.NotNil(t, err)
	assert.Equal(t, err, errors.New("unknown launchplan"))
}
