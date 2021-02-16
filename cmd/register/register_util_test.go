package register

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return GetDoFunc(req)
}

var (
	reader               *os.File
	writer               *os.File
	ctx                  context.Context
	args                 []string
	stdOut               *os.File
	stderr               *os.File
	GetDoFunc func(req *http.Request) (*http.Response, error)
)

func setup() {
	ctx = context.Background()
	Client = &MockClient{}
	validTar, err := os.Open("testdata/valid-register.tar")
	if err != nil {
		fmt.Printf("unexpected error: %v", err)
		os.Exit(-1)
	}
	response := &http.Response{
		Body : validTar,
	}
	GetDoFunc = func(*http.Request) (*http.Response, error) {
		return response, nil
	}
	reader, writer, err = os.Pipe()
	if err != nil {
		panic(err)
	}
	stdOut = os.Stdout
	stderr = os.Stderr
	os.Stdout = writer
	os.Stderr = writer
	log.SetOutput(writer)
}

func teardownAndVerify(t *testing.T, expectedLog string) {
	writer.Close()
	os.Stdout = stdOut
	os.Stderr = stderr
	var buf bytes.Buffer
	io.Copy(&buf, reader)
	assert.Equal(t, expectedLog, buf.String())
}

func TestGetSortedFileList(t *testing.T) {
	setup()
	filesConfig.archive = false
	args = [] string{"file2", "file1"}
	fileList, tmpDir, err := getSortedFileList(ctx, args)
	assert.Equal(t, "file1",fileList[0])
	assert.Equal(t, "file2", fileList[1])
	assert.Equal(t, tmpDir, "")
	assert.Nil(t, err)
	teardownAndVerify(t, "")
}

func TestGetSortedArchivedFileWithParentFolderList(t *testing.T) {
	setup()
	filesConfig.archive = true
	args = [] string{"testdata/valid-parent-folder-register.tar"}
	fileList, tmpDir, err := getSortedFileList(ctx, args)
	assert.Equal(t, len(fileList),4)
	assert.Equal(t, filepath.Join(tmpDir, "parentfolder", "014_recipes.core.basic.basic_workflow.t1_1.pb"),fileList[0])
	assert.Equal(t, filepath.Join(tmpDir, "parentfolder", "015_recipes.core.basic.basic_workflow.t2_1.pb"),fileList[1])
	assert.Equal(t, filepath.Join(tmpDir, "parentfolder", "016_recipes.core.basic.basic_workflow.my_wf_2.pb"),fileList[2])
	assert.Equal(t, filepath.Join(tmpDir, "parentfolder", "017_recipes.core.basic.basic_workflow.my_wf_3.pb"),fileList[3])
	assert.True(t, strings.HasPrefix(tmpDir, "/tmp/register"))
	assert.Nil(t, err)
	// Clean up the temp directory.
	assert.Nil(t, os.RemoveAll(tmpDir), "unable to delete temp dir %v", tmpDir)
	teardownAndVerify(t, "")
}

func TestGetSortedArchivedFileList(t *testing.T) {
	setup()
	filesConfig.archive = true
	args = [] string{"testdata/valid-register.tar"}
	fileList, tmpDir, err := getSortedFileList(ctx, args)
	assert.Equal(t, len(fileList),4)
	assert.Equal(t, filepath.Join(tmpDir, "014_recipes.core.basic.basic_workflow.t1_1.pb"),fileList[0])
	assert.Equal(t, filepath.Join(tmpDir, "015_recipes.core.basic.basic_workflow.t2_1.pb"),fileList[1])
	assert.Equal(t, filepath.Join(tmpDir, "016_recipes.core.basic.basic_workflow.my_wf_2.pb"),fileList[2])
	assert.Equal(t, filepath.Join(tmpDir, "017_recipes.core.basic.basic_workflow.my_wf_3.pb"),fileList[3])
	assert.True(t, strings.HasPrefix(tmpDir, "/tmp/register"))
	assert.Nil(t, err)
	// Clean up the temp directory.
	assert.Nil(t, os.RemoveAll(tmpDir), "unable to delete temp dir %v", tmpDir)
	teardownAndVerify(t, "")
}

func TestGetSortedArchivedFileUnorderedList(t *testing.T) {
	setup()
	filesConfig.archive = true
	args = [] string{"testdata/valid-unordered-register.tar"}
	fileList, tmpDir, err := getSortedFileList(ctx, args)
	assert.Equal(t, len(fileList),4)
	assert.Equal(t, filepath.Join(tmpDir, "014_recipes.core.basic.basic_workflow.t1_1.pb"),fileList[0])
	assert.Equal(t, filepath.Join(tmpDir, "015_recipes.core.basic.basic_workflow.t2_1.pb"),fileList[1])
	assert.Equal(t, filepath.Join(tmpDir, "016_recipes.core.basic.basic_workflow.my_wf_2.pb"),fileList[2])
	assert.Equal(t, filepath.Join(tmpDir, "017_recipes.core.basic.basic_workflow.my_wf_3.pb"),fileList[3])
	assert.True(t, strings.HasPrefix(tmpDir, "/tmp/register"))
	assert.Nil(t, err)
	// Clean up the temp directory.
	assert.Nil(t, os.RemoveAll(tmpDir), "unable to delete temp dir %v", tmpDir)
	teardownAndVerify(t, "")
}

func TestGetSortedArchivedCorruptedFileList(t *testing.T) {
	setup()
	filesConfig.archive = true
	args = [] string{"testdata/invalid.tar"}
	fileList, tmpDir, err := getSortedFileList(ctx, args)
	assert.Equal(t, len(fileList),0)
	assert.True(t, strings.HasPrefix(tmpDir, "/tmp/register"))
	assert.NotNil(t, err)
	// Clean up the temp directory.
	assert.Nil(t, os.RemoveAll(tmpDir), "unable to delete temp dir %v", tmpDir)
	teardownAndVerify(t, "")
}

func TestGetSortedArchivedTgzList(t *testing.T) {
	setup()
	filesConfig.archive = true
	args = [] string{"testdata/valid-register.tgz"}
	fileList, tmpDir, err := getSortedFileList(ctx, args)
	assert.Equal(t, len(fileList),4)
	assert.Equal(t, filepath.Join(tmpDir, "014_recipes.core.basic.basic_workflow.t1_1.pb"),fileList[0])
	assert.Equal(t, filepath.Join(tmpDir, "015_recipes.core.basic.basic_workflow.t2_1.pb"),fileList[1])
	assert.Equal(t, filepath.Join(tmpDir, "016_recipes.core.basic.basic_workflow.my_wf_2.pb"),fileList[2])
	assert.Equal(t, filepath.Join(tmpDir, "017_recipes.core.basic.basic_workflow.my_wf_3.pb"),fileList[3])
	assert.True(t, strings.HasPrefix(tmpDir, "/tmp/register"))
	assert.Nil(t, err)
	// Clean up the temp directory.
	assert.Nil(t, os.RemoveAll(tmpDir), "unable to delete temp dir %v", tmpDir)
	teardownAndVerify(t, "")
}

func TestGetSortedArchivedCorruptedTgzFileList(t *testing.T) {
	setup()
	filesConfig.archive = true
	args = [] string{"testdata/invalid.tgz"}
	fileList, tmpDir, err := getSortedFileList(ctx, args)
	assert.Equal(t, 0, len(fileList))
	assert.True(t, strings.HasPrefix(tmpDir, "/tmp/register"))
	assert.NotNil(t, err)
	// Clean up the temp directory.
	assert.Nil(t, os.RemoveAll(tmpDir), "unable to delete temp dir %v", tmpDir)
	teardownAndVerify(t, "")
}

func TestGetSortedArchivedInvalidArchiveFileList(t *testing.T) {
	setup()
	filesConfig.archive = true
	args = [] string{"testdata/invalid-extension-register.zip"}
	fileList, tmpDir, err := getSortedFileList(ctx, args)
	assert.Equal(t, 0, len(fileList))
	assert.True(t, strings.HasPrefix(tmpDir, "/tmp/register"))
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("only .tar and .tgz extension archives are supported"), err)
	// Clean up the temp directory.
	assert.Nil(t, os.RemoveAll(tmpDir), "unable to delete temp dir %v", tmpDir)
	teardownAndVerify(t, "")
}

func TestGetSortedArchivedFileThroughInvalidHttpList(t *testing.T) {
	setup()
	filesConfig.archive = true
	args = [] string{"http://invalidhost:invalidport/testdata/valid-register.tar"}
	fileList, tmpDir, err := getSortedFileList(ctx, args)
	assert.Equal(t, 0, len(fileList))
	assert.True(t, strings.HasPrefix(tmpDir, "/tmp/register"))
	assert.NotNil(t, err)
	// Clean up the temp directory.
	assert.Nil(t, os.RemoveAll(tmpDir), "unable to delete temp dir %v", tmpDir)
	teardownAndVerify(t, "")
}

func TestGetSortedArchivedFileThroughValidHttpList(t *testing.T) {
	setup()
	filesConfig.archive = true
	args = [] string{"http://dummyhost:80/testdata/valid-register.tar"}
	fileList, tmpDir, err := getSortedFileList(ctx, args)
	assert.Equal(t, len(fileList),4)
	assert.Equal(t, filepath.Join(tmpDir, "014_recipes.core.basic.basic_workflow.t1_1.pb"),fileList[0])
	assert.Equal(t, filepath.Join(tmpDir, "015_recipes.core.basic.basic_workflow.t2_1.pb"),fileList[1])
	assert.Equal(t, filepath.Join(tmpDir, "016_recipes.core.basic.basic_workflow.my_wf_2.pb"),fileList[2])
	assert.Equal(t, filepath.Join(tmpDir, "017_recipes.core.basic.basic_workflow.my_wf_3.pb"),fileList[3])
	assert.True(t, strings.HasPrefix(tmpDir, "/tmp/register"))
	assert.Nil(t, err)
	// Clean up the temp directory.
	assert.Nil(t, os.RemoveAll(tmpDir), "unable to delete temp dir %v", tmpDir)
	teardownAndVerify(t, "")
}

func TestGetSortedArchivedFileThroughValidHttpWithNullContextList(t *testing.T) {
	setup()
	filesConfig.archive = true
	args = [] string{"http://dummyhost:80/testdata/valid-register.tar"}
	ctx = nil
	fileList, tmpDir, err := getSortedFileList(ctx, args)
	assert.Equal(t, len(fileList),0)
	assert.True(t, strings.HasPrefix(tmpDir, "/tmp/register"))
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("net/http: nil Context"),err)
	// Clean up the temp directory.
	assert.Nil(t, os.RemoveAll(tmpDir), "unable to delete temp dir %v", tmpDir)
	teardownAndVerify(t, "")
}