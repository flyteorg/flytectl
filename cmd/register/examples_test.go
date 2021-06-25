package register

import (
	"testing"

	storageMocks "github.com/flyteorg/flytestdlib/storage/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterExamplesFunc(t *testing.T) {
	setup()
	registerFilesSetup()
	mockStorage := &storageMocks.ComposedProtobufStore{}
	args = []string{""}
	mockStorage.OnWriteRawMatch(ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockAdminClient.OnCreateTaskMatch(mock.Anything, mock.Anything).Return(nil, nil)
	mockAdminClient.OnCreateWorkflowMatch(mock.Anything, mock.Anything).Return(nil, nil)
	mockAdminClient.OnCreateLaunchPlanMatch(mock.Anything, mock.Anything).Return(nil, nil)
	Client = mockStorage
	err := registerExamplesFunc(ctx, args, cmdCtx)
	assert.Nil(t, err)
}
func TestRegisterExamplesFuncErr(t *testing.T) {
	setup()
	registerFilesSetup()
	mockStorage := &storageMocks.ComposedProtobufStore{}
	githubRepository = "testingsnacks"
	args = []string{""}
	mockStorage.OnWriteRawMatch(ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockAdminClient.OnCreateTaskMatch(mock.Anything, mock.Anything).Return(nil, nil)
	mockAdminClient.OnCreateWorkflowMatch(mock.Anything, mock.Anything).Return(nil, nil)
	mockAdminClient.OnCreateLaunchPlanMatch(mock.Anything, mock.Anything).Return(nil, nil)
	Client = mockStorage
	err := registerExamplesFunc(ctx, args, cmdCtx)
	// TODO (Yuvraj) make test to success after fixing flytesnacks bug
	assert.NotNil(t, err)
	githubRepository = "flytesnacks"
}
