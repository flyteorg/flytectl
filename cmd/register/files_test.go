package register

import (
	"fmt"
	"testing"

	storageMocks "github.com/flyteorg/flytestdlib/storage/mocks"

	rconfig "github.com/flyteorg/flytectl/cmd/config/subcommand/register"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	s3Output = "s3://dummy/prefix"
)

func TestRegisterFromFiles(t *testing.T) {
	t.Run("Valid registration", func(t *testing.T) {
		setup()
		registerFilesSetup()
		rconfig.DefaultFilesConfig.Archive = true
		args = []string{"testdata/valid-parent-folder-register.tar"}
		mockAdminClient.OnCreateTaskMatch(mock.Anything, mock.Anything).Return(nil, nil)
		mockAdminClient.OnCreateWorkflowMatch(mock.Anything, mock.Anything).Return(nil, nil)
		mockAdminClient.OnCreateLaunchPlanMatch(mock.Anything, mock.Anything).Return(nil, nil)
		err := registerFromFilesFunc(ctx, args, cmdCtx)
		assert.Nil(t, err)
	})
	t.Run("Valid fast registration", func(t *testing.T) {
		setup()
		registerFilesSetup()
		rconfig.DefaultFilesConfig.Archive = true

		rconfig.DefaultFilesConfig.DestinationDir = "/"
		rconfig.DefaultFilesConfig.OutputLocationPrefix = s3Output
		rconfig.DefaultFilesConfig.AdditionalDistributionDir = s3Output
		mockStorage := &storageMocks.ComposedProtobufStore{}
		args = []string{"testdata/flytesnacks-core.tgz"}
		mockStorage.OnWriteRawMatch(ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
		mockAdminClient.OnCreateTaskMatch(mock.Anything, mock.Anything).Return(nil, nil)
		mockAdminClient.OnCreateWorkflowMatch(mock.Anything, mock.Anything).Return(nil, nil)
		mockAdminClient.OnCreateLaunchPlanMatch(mock.Anything, mock.Anything).Return(nil, nil)
		Client = mockStorage
		err := registerFromFilesFunc(ctx, args, cmdCtx)
		assert.Nil(t, err)
	})
	t.Run("Failed fast registration while uploading the codebase", func(t *testing.T) {
		setup()
		registerFilesSetup()
		rconfig.DefaultFilesConfig.Archive = true
		rconfig.DefaultFilesConfig.DestinationDir = "/"
		rconfig.DefaultFilesConfig.OutputLocationPrefix = s3Output
		rconfig.DefaultFilesConfig.ContinueOnError = false
		rconfig.DefaultFilesConfig.AdditionalDistributionDir = s3Output
		mockStorage := &storageMocks.ComposedProtobufStore{}
		args = []string{"testdata/valid-fast-register.tgz"}
		mockStorage.OnWriteRawMatch(ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("error"))
		mockAdminClient.OnCreateTaskMatch(mock.Anything, mock.Anything).Return(nil, nil)
		mockAdminClient.OnCreateWorkflowMatch(mock.Anything, mock.Anything).Return(nil, nil)
		mockAdminClient.OnCreateLaunchPlanMatch(mock.Anything, mock.Anything).Return(nil, nil)
		Client = mockStorage
		err := Register(ctx, args, cmdCtx)
		assert.NotNil(t, err)
	})
	t.Run("Failed fast registration while flags are missing", func(t *testing.T) {
		setup()
		registerFilesSetup()
		rconfig.DefaultFilesConfig.Archive = true
		rconfig.DefaultFilesConfig.OutputLocationPrefix = s3Output
		mockStorage := &storageMocks.ComposedProtobufStore{}
		args = []string{"testdata/valid-fast-register.tgz"}
		mockStorage.OnWriteRawMatch(ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
		mockAdminClient.OnCreateTaskMatch(mock.Anything, mock.Anything).Return(nil, nil)
		mockAdminClient.OnCreateWorkflowMatch(mock.Anything, mock.Anything).Return(nil, nil)
		mockAdminClient.OnCreateLaunchPlanMatch(mock.Anything, mock.Anything).Return(nil, nil)
		Client = mockStorage
		err := registerFromFilesFunc(ctx, args, cmdCtx)
		assert.NotNil(t, err)
	})
	t.Run("Valid registration of fast serialize", func(t *testing.T) {
		setup()
		registerFilesSetup()
		rconfig.DefaultFilesConfig.Archive = true
		rconfig.DefaultFilesConfig.DestinationDir = "/"
		rconfig.DefaultFilesConfig.OutputLocationPrefix = s3Output
		rconfig.DefaultFilesConfig.AdditionalDistributionDir = s3Output
		mockStorage := &storageMocks.ComposedProtobufStore{}
		args = []string{"testdata/flytesnacks-core.tgz"}
		mockStorage.OnWriteRawMatch(ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
		mockAdminClient.OnCreateTaskMatch(mock.Anything, mock.Anything).Return(nil, nil)
		mockAdminClient.OnCreateWorkflowMatch(mock.Anything, mock.Anything).Return(nil, nil)
		mockAdminClient.OnCreateLaunchPlanMatch(mock.Anything, mock.Anything).Return(nil, nil)
		Client = mockStorage
		err := registerFromFilesFunc(ctx, args, cmdCtx)
		assert.Nil(t, err)
	})

	t.Run("Registration with proto files ", func(t *testing.T) {
		setup()
		registerFilesSetup()
		rconfig.DefaultFilesConfig.Archive = false
		rconfig.DefaultFilesConfig.DestinationDir = "/"
		rconfig.DefaultFilesConfig.OutputLocationPrefix = s3Output
		rconfig.DefaultFilesConfig.AdditionalDistributionDir = s3Output
		mockStorage := &storageMocks.ComposedProtobufStore{}
		args = []string{"testdata/69_core.flyte_basics.lp.greet_1.pb"}
		mockStorage.OnWriteRawMatch(ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
		mockAdminClient.OnCreateTaskMatch(mock.Anything, mock.Anything).Return(nil, nil)
		mockAdminClient.OnCreateWorkflowMatch(mock.Anything, mock.Anything).Return(nil, nil)
		mockAdminClient.OnCreateLaunchPlanMatch(mock.Anything, mock.Anything).Return(nil, nil)
		Client = mockStorage
		err := registerFromFilesFunc(ctx, args, cmdCtx)
		assert.Nil(t, err)
	})

}
