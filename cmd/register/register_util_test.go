package register

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"

	storageMocks "github.com/flyteorg/flytestdlib/storage/mocks"

	rconfig "github.com/flyteorg/flytectl/cmd/config/subcommand/register"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	u "github.com/flyteorg/flytectl/cmd/testutils"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return GetDoFunc(req)
}

func registerFilesSetup() {
	httpClient = &MockHTTPClient{}
	validTar, err := os.Open("testdata/valid-register.tar")
	if err != nil {
		fmt.Printf("unexpected error: %v", err)
		os.Exit(-1)
	}
	response := &http.Response{
		Body: validTar,
	}
	GetDoFunc = func(*http.Request) (*http.Response, error) {
		return response, nil
	}
	ctx = u.Ctx
	mockAdminClient = u.MockClient
	cmdCtx = cmdCore.NewCommandContext(mockAdminClient, u.MockOutStream)

	rconfig.DefaultFilesConfig.AssumableIamRole = ""
	rconfig.DefaultFilesConfig.K8ServiceAccount = ""
	rconfig.DefaultFilesConfig.OutputLocationPrefix = ""
}

func TestGetSortedArchivedFileWithParentFolderList(t *testing.T) {
	setup()
	registerFilesSetup()
	rconfig.DefaultFilesConfig.Archive = true
	args = []string{"testdata/valid-parent-folder-register.tar"}
	fileList, tmpDir, err := getSerializeOutputFiles(ctx, args)
	assert.Equal(t, len(fileList), 4)
	assert.Equal(t, filepath.Join(tmpDir, "parentfolder", "014_recipes.core.basic.basic_workflow.t1_1.pb"), fileList[0])
	assert.Equal(t, filepath.Join(tmpDir, "parentfolder", "015_recipes.core.basic.basic_workflow.t2_1.pb"), fileList[1])
	assert.Equal(t, filepath.Join(tmpDir, "parentfolder", "016_recipes.core.basic.basic_workflow.my_wf_2.pb"), fileList[2])
	assert.Equal(t, filepath.Join(tmpDir, "parentfolder", "017_recipes.core.basic.basic_workflow.my_wf_3.pb"), fileList[3])
	assert.True(t, strings.HasPrefix(tmpDir, "/tmp/register"))
	assert.Nil(t, err)
	// Clean up the temp directory.
	assert.Nil(t, os.RemoveAll(tmpDir), "unable to delete temp dir %v", tmpDir)
}

func TestGetSortedArchivedFileList(t *testing.T) {
	setup()
	registerFilesSetup()
	rconfig.DefaultFilesConfig.Archive = true
	args = []string{"testdata/valid-register.tar"}
	fileList, tmpDir, err := getSerializeOutputFiles(ctx, args)
	assert.Equal(t, len(fileList), 4)
	assert.Equal(t, filepath.Join(tmpDir, "014_recipes.core.basic.basic_workflow.t1_1.pb"), fileList[0])
	assert.Equal(t, filepath.Join(tmpDir, "015_recipes.core.basic.basic_workflow.t2_1.pb"), fileList[1])
	assert.Equal(t, filepath.Join(tmpDir, "016_recipes.core.basic.basic_workflow.my_wf_2.pb"), fileList[2])
	assert.Equal(t, filepath.Join(tmpDir, "017_recipes.core.basic.basic_workflow.my_wf_3.pb"), fileList[3])
	assert.True(t, strings.HasPrefix(tmpDir, "/tmp/register"))
	assert.Nil(t, err)
	// Clean up the temp directory.
	assert.Nil(t, os.RemoveAll(tmpDir), "unable to delete temp dir %v", tmpDir)
}

func TestGetSortedArchivedFileUnorderedList(t *testing.T) {
	setup()
	registerFilesSetup()
	rconfig.DefaultFilesConfig.Archive = true
	args = []string{"testdata/valid-unordered-register.tar"}
	fileList, tmpDir, err := getSerializeOutputFiles(ctx, args)
	assert.Equal(t, len(fileList), 4)
	assert.Equal(t, filepath.Join(tmpDir, "014_recipes.core.basic.basic_workflow.t1_1.pb"), fileList[0])
	assert.Equal(t, filepath.Join(tmpDir, "015_recipes.core.basic.basic_workflow.t2_1.pb"), fileList[1])
	assert.Equal(t, filepath.Join(tmpDir, "016_recipes.core.basic.basic_workflow.my_wf_2.pb"), fileList[2])
	assert.Equal(t, filepath.Join(tmpDir, "017_recipes.core.basic.basic_workflow.my_wf_3.pb"), fileList[3])
	assert.True(t, strings.HasPrefix(tmpDir, "/tmp/register"))
	assert.Nil(t, err)
	// Clean up the temp directory.
	assert.Nil(t, os.RemoveAll(tmpDir), "unable to delete temp dir %v", tmpDir)
}

func TestGetSortedArchivedCorruptedFileList(t *testing.T) {
	setup()
	registerFilesSetup()
	rconfig.DefaultFilesConfig.Archive = true
	args = []string{"testdata/invalid.tar"}
	fileList, tmpDir, err := getSerializeOutputFiles(ctx, args)
	assert.Equal(t, len(fileList), 0)
	assert.True(t, strings.HasPrefix(tmpDir, "/tmp/register"))
	assert.NotNil(t, err)
	// Clean up the temp directory.
	assert.Nil(t, os.RemoveAll(tmpDir), "unable to delete temp dir %v", tmpDir)
}

func TestGetSortedArchivedTgzList(t *testing.T) {
	setup()
	registerFilesSetup()
	rconfig.DefaultFilesConfig.Archive = true
	args = []string{"testdata/valid-register.tgz"}
	fileList, tmpDir, err := getSerializeOutputFiles(ctx, args)
	assert.Equal(t, len(fileList), 4)
	assert.Equal(t, filepath.Join(tmpDir, "014_recipes.core.basic.basic_workflow.t1_1.pb"), fileList[0])
	assert.Equal(t, filepath.Join(tmpDir, "015_recipes.core.basic.basic_workflow.t2_1.pb"), fileList[1])
	assert.Equal(t, filepath.Join(tmpDir, "016_recipes.core.basic.basic_workflow.my_wf_2.pb"), fileList[2])
	assert.Equal(t, filepath.Join(tmpDir, "017_recipes.core.basic.basic_workflow.my_wf_3.pb"), fileList[3])
	assert.True(t, strings.HasPrefix(tmpDir, "/tmp/register"))
	assert.Nil(t, err)
	// Clean up the temp directory.
	assert.Nil(t, os.RemoveAll(tmpDir), "unable to delete temp dir %v", tmpDir)
}

func TestGetSortedArchivedCorruptedTgzFileList(t *testing.T) {
	setup()
	rconfig.DefaultFilesConfig.Archive = true
	args = []string{"testdata/invalid.tgz"}
	fileList, tmpDir, err := getSerializeOutputFiles(ctx, args)
	assert.Equal(t, 0, len(fileList))
	assert.True(t, strings.HasPrefix(tmpDir, "/tmp/register"))
	assert.NotNil(t, err)
	// Clean up the temp directory.
	assert.Nil(t, os.RemoveAll(tmpDir), "unable to delete temp dir %v", tmpDir)
}

func TestGetSortedArchivedInvalidArchiveFileList(t *testing.T) {
	setup()
	registerFilesSetup()
	rconfig.DefaultFilesConfig.Archive = true
	args = []string{"testdata/invalid-extension-register.zip"}
	fileList, tmpDir, err := getSerializeOutputFiles(ctx, args)
	assert.Equal(t, 0, len(fileList))
	assert.True(t, strings.HasPrefix(tmpDir, "/tmp/register"))
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("only .tar and .tgz extension archives are supported"), err)
	// Clean up the temp directory.
	assert.Nil(t, os.RemoveAll(tmpDir), "unable to delete temp dir %v", tmpDir)
}

func TestGetSortedArchivedFileThroughInvalidHttpList(t *testing.T) {
	setup()
	rconfig.DefaultFilesConfig.Archive = true
	args = []string{"http://invalidhost:invalidport/testdata/valid-register.tar"}
	fileList, tmpDir, err := getSerializeOutputFiles(ctx, args)
	assert.Equal(t, 0, len(fileList))
	assert.True(t, strings.HasPrefix(tmpDir, "/tmp/register"))
	assert.NotNil(t, err)
	// Clean up the temp directory.
	assert.Nil(t, os.RemoveAll(tmpDir), "unable to delete temp dir %v", tmpDir)
}

func TestGetSortedArchivedFileThroughValidHttpList(t *testing.T) {
	setup()
	registerFilesSetup()
	rconfig.DefaultFilesConfig.Archive = true
	args = []string{"http://dummyhost:80/testdata/valid-register.tar"}
	fileList, tmpDir, err := getSerializeOutputFiles(ctx, args)
	assert.Equal(t, len(fileList), 4)
	assert.Equal(t, filepath.Join(tmpDir, "014_recipes.core.basic.basic_workflow.t1_1.pb"), fileList[0])
	assert.Equal(t, filepath.Join(tmpDir, "015_recipes.core.basic.basic_workflow.t2_1.pb"), fileList[1])
	assert.Equal(t, filepath.Join(tmpDir, "016_recipes.core.basic.basic_workflow.my_wf_2.pb"), fileList[2])
	assert.Equal(t, filepath.Join(tmpDir, "017_recipes.core.basic.basic_workflow.my_wf_3.pb"), fileList[3])
	assert.True(t, strings.HasPrefix(tmpDir, "/tmp/register"))
	assert.Nil(t, err)
	// Clean up the temp directory.
	assert.Nil(t, os.RemoveAll(tmpDir), "unable to delete temp dir %v", tmpDir)
}

func TestGetSortedArchivedFileThroughValidHttpWithNullContextList(t *testing.T) {
	setup()
	registerFilesSetup()
	rconfig.DefaultFilesConfig.Archive = true
	args = []string{"http://dummyhost:80/testdata/valid-register.tar"}
	ctx = nil
	fileList, tmpDir, err := getSerializeOutputFiles(ctx, args)
	assert.Equal(t, len(fileList), 0)
	assert.True(t, strings.HasPrefix(tmpDir, "/tmp/register"))
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("net/http: nil Context"), err)
	// Clean up the temp directory.
	assert.Nil(t, os.RemoveAll(tmpDir), "unable to delete temp dir %v", tmpDir)
}

func TestRegisterFile(t *testing.T) {
	t.Run("Successful run", func(t *testing.T) {
		setup()
		registerFilesSetup()
		mockAdminClient.OnCreateTaskMatch(mock.Anything, mock.Anything).Return(nil, nil)
		args = []string{"testdata/69_core.flyte_basics.lp.greet_1.pb"}
		var registerResults []Result
		results, err := registerFile(ctx, args[0], "", registerResults, cmdCtx)
		assert.Equal(t, 1, len(results))
		assert.Nil(t, err)
	})
	t.Run("Non existent file", func(t *testing.T) {
		setup()
		registerFilesSetup()
		args = []string{"testdata/non-existent.pb"}
		var registerResults []Result
		results, err := registerFile(ctx, args[0], "", registerResults, cmdCtx)
		assert.Equal(t, 1, len(results))
		assert.Equal(t, "Failed", results[0].Status)
		assert.Equal(t, "Error reading file due to open testdata/non-existent.pb: no such file or directory", results[0].Info)
		assert.NotNil(t, err)
	})
	t.Run("unmarhal failure", func(t *testing.T) {
		setup()
		registerFilesSetup()
		args = []string{"testdata/valid-register.tar"}
		var registerResults []Result
		results, err := registerFile(ctx, args[0], "", registerResults, cmdCtx)
		assert.Equal(t, 1, len(results))
		assert.Equal(t, "Failed", results[0].Status)
		assert.Equal(t, "Error unmarshalling file due to failed unmarshalling file testdata/valid-register.tar", results[0].Info)
		assert.NotNil(t, err)
	})
	t.Run("AlreadyExists", func(t *testing.T) {
		setup()
		registerFilesSetup()
		mockAdminClient.OnCreateTaskMatch(mock.Anything, mock.Anything).Return(nil,
			status.Error(codes.AlreadyExists, "AlreadyExists"))
		args = []string{"testdata/69_core.flyte_basics.lp.greet_1.pb"}
		var registerResults []Result
		results, err := registerFile(ctx, args[0], "", registerResults, cmdCtx)
		assert.Equal(t, 1, len(results))
		assert.Equal(t, "Success", results[0].Status)
		assert.Equal(t, "AlreadyExists", results[0].Info)
		assert.Nil(t, err)
	})
	t.Run("Registration Error", func(t *testing.T) {
		setup()
		registerFilesSetup()
		mockAdminClient.OnCreateTaskMatch(mock.Anything, mock.Anything).Return(nil,
			status.Error(codes.InvalidArgument, "Invalid"))
		args = []string{"testdata/69_core.flyte_basics.lp.greet_1.pb"}
		var registerResults []Result
		results, err := registerFile(ctx, args[0], "", registerResults, cmdCtx)
		assert.Equal(t, 1, len(results))
		assert.Equal(t, "Failed", results[0].Status)
		assert.Equal(t, "Error registering file due to rpc error: code = InvalidArgument desc = Invalid", results[0].Info)
		assert.NotNil(t, err)
	})
}

func TestHydrateLaunchPlanSpec(t *testing.T) {
	t.Run("IamRole override", func(t *testing.T) {
		setup()
		registerFilesSetup()
		rconfig.DefaultFilesConfig.AssumableIamRole = "iamRole"
		lpSpec := &admin.LaunchPlanSpec{}
		hydrateLaunchPlanSpec(lpSpec)
		assert.Equal(t, &admin.AuthRole{AssumableIamRole: "iamRole"}, lpSpec.AuthRole)
	})
	t.Run("k8Service account override", func(t *testing.T) {
		setup()
		registerFilesSetup()
		rconfig.DefaultFilesConfig.K8ServiceAccount = "k8Account"
		lpSpec := &admin.LaunchPlanSpec{}
		hydrateLaunchPlanSpec(lpSpec)
		assert.Equal(t, &admin.AuthRole{KubernetesServiceAccount: "k8Account"}, lpSpec.AuthRole)
	})
	t.Run("Both k8Service and IamRole", func(t *testing.T) {
		setup()
		registerFilesSetup()
		rconfig.DefaultFilesConfig.AssumableIamRole = "iamRole"
		rconfig.DefaultFilesConfig.K8ServiceAccount = "k8Account"
		lpSpec := &admin.LaunchPlanSpec{}
		hydrateLaunchPlanSpec(lpSpec)
		assert.Equal(t, &admin.AuthRole{AssumableIamRole: "iamRole",
			KubernetesServiceAccount: "k8Account"}, lpSpec.AuthRole)
	})
	t.Run("Output prefix", func(t *testing.T) {
		setup()
		registerFilesSetup()
		rconfig.DefaultFilesConfig.OutputLocationPrefix = "prefix"
		lpSpec := &admin.LaunchPlanSpec{}
		hydrateLaunchPlanSpec(lpSpec)
		assert.Equal(t, &admin.RawOutputDataConfig{OutputLocationPrefix: "prefix"}, lpSpec.RawOutputDataConfig)
	})
}

func TestFlyteManifest(t *testing.T) {
	_, tag, err := getFlyteTestManifest(githubOrg, githubRepository)
	assert.Nil(t, err)
	assert.Contains(t, tag, "v")
	assert.NotEmpty(t, tag)
}

func TestGetAdditionalDistributionLoc(t *testing.T) {
	remoteLocation := getRemoteStoragePath("s3://dummy", "test.tar.gz", "v1")
	assert.Equal(t, "s3://dummy/v1-test.tar.gz", string(remoteLocation))
}

func TestUploadFastRegisterArtifact(t *testing.T) {
	t.Run("Successful upload", func(t *testing.T) {
		rawStoreWrite := &storageMocks.ComposedProtobufStore{}
		Client = rawStoreWrite
		rawStoreWrite.OnWriteRawMatch(ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
		err := uploadFastRegisterArtifact(ctx, "testdata/flytesnacks-core.tgz", "", "", "")
		assert.Nil(t, err)
	})
	t.Run("Failed upload", func(t *testing.T) {
		rawStoreWrite := &storageMocks.ComposedProtobufStore{}
		Client = rawStoreWrite
		rawStoreWrite.OnWriteRawMatch(ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("error"))
		err := uploadFastRegisterArtifact(ctx, "testdata/flytesnacks-core.tgz", "", "", "")
		assert.NotNil(t, err)
	})
	t.Run("Failed upload", func(t *testing.T) {
		rawStoreWrite := &storageMocks.ComposedProtobufStore{}
		Client = rawStoreWrite
		rawStoreWrite.OnWriteRawMatch(ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
		err := uploadFastRegisterArtifact(ctx, "testdata/flytesnacksre.tgz", "", "", "")
		assert.NotNil(t, err)
	})
}

func TestGetStorageClient(t *testing.T) {
	t.Run("Failed to create storage client", func(t *testing.T) {
		Client = nil
		s, err := getStorageClient(ctx)
		assert.NotNil(t, err)
		assert.Nil(t, s)
	})
}

func TestGetFlyteTestManifest(t *testing.T) {
	t.Run("Failed to get manifest with wrong name", func(t *testing.T) {
		_, tag, err := getFlyteTestManifest("no////ne", "no////ne")
		assert.NotNil(t, err)
		assert.Equal(t, len(tag), 0)
	})
	t.Run("Failed to get release", func(t *testing.T) {
		_, tag, err := getFlyteTestManifest("flyteorg", "homebrew-tap")
		assert.NotNil(t, err)
		assert.Equal(t, len(tag), 0)
	})
	t.Run("Failed to get manifest", func(t *testing.T) {
		flyteManifest = ""
		_, tag, err := getFlyteTestManifest("flyteorg", "flytesnacks")
		assert.NotNil(t, err)
		assert.Equal(t, len(tag), 0)
	})
}

func TestRegister(t *testing.T) {
	t.Run("Failed to register", func(t *testing.T) {
		setup()
		registerFilesSetup()
		node := &admin.NodeExecution{}
		err := register(ctx, node, cmdCtx)
		assert.NotNil(t, err)
	})
}

func TestHydrateNode(t *testing.T) {
	t.Run("Failed hydrate node", func(t *testing.T) {
		setup()
		registerFilesSetup()
		node := &core.Node{}
		err := hydrateNode(node)
		assert.NotNil(t, err)
	})

	t.Run("hydrateSpec with wrong type", func(t *testing.T) {
		setup()
		registerFilesSetup()
		task := &admin.Task{}
		err := hydrateSpec(task, "")
		assert.NotNil(t, err)
	})
}
