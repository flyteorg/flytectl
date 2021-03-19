package get

import (
	"testing"

	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"

	"github.com/stretchr/testify/assert"
)

func TestGetTaskFunc(t *testing.T) {
	setup()
	err = getTaskFunc(ctx, argsTask, cmdCtx)
	assert.Nil(t, err)
	mockClient.AssertCalled(t, "ListTasks", ctx, resourceListRequestTask)
	teardownAndVerify(t, `[
	{
		"id": {
			"name": "task1",
			"version": "v2"
		},
		"closure": {
			"createdAt": "1970-01-01T00:00:01Z"
		}
	},
	{
		"id": {
			"name": "task1",
			"version": "v1"
		},
		"closure": {
			"createdAt": "1970-01-01T00:00:00Z"
		}
	}
]`)
}

func TestGetTaskFuncLatest(t *testing.T) {
	setup()
	taskConfig.Latest = true
	err = getTaskFunc(ctx, argsTask, cmdCtx)
	assert.Nil(t, err)
	mockClient.AssertCalled(t, "ListTasks", ctx, resourceListRequestTask)
	teardownAndVerify(t, `{
	"id": {
		"name": "task1",
		"version": "v2"
	},
	"closure": {
		"createdAt": "1970-01-01T00:00:01Z"
	}
}`)
}

func TestGetTaskWithVersion(t *testing.T) {
	setup()
	taskConfig.Version = "v2"
	objectGetRequest.Id.ResourceType = core.ResourceType_TASK
	err = getTaskFunc(ctx, argsTask, cmdCtx)
	assert.Nil(t, err)
	mockClient.AssertCalled(t, "GetTask", ctx, objectGetRequestTask)
	teardownAndVerify(t, `{
	"id": {
		"name": "task1",
		"version": "v2"
	},
	"closure": {
		"createdAt": "1970-01-01T00:00:01Z"
	}
}`)
}

func TestGetTasks(t *testing.T) {
	setup()
	argsTask = []string{}
	err = getTaskFunc(ctx, argsTask, cmdCtx)
	assert.Nil(t, err)
	mockClient.AssertCalled(t, "ListTaskIds", ctx, namedIDRequest)
	teardownAndVerify(t, `[
	{
		"project": "dummyProject",
		"domain": "dummyDomain",
		"name": "task1"
	},
	{
		"project": "dummyProject",
		"domain": "dummyDomain",
		"name": "task2"
	}
]`)
}
