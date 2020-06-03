package cmd

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/protobuf/jsonpb"

	"github.com/stretchr/testify/mock"

	adminIdl "github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
	coreIdl "github.com/lyft/flyteidl/gen/pb-go/flyteidl/core"

	"github.com/lyft/flyteidl/clients/go/admin"
	"github.com/lyft/flytestdlib/config"
	"github.com/lyft/flytestdlib/config/viper"

	"github.com/lyft/flyteidl/clients/go/admin/mocks"

	"github.com/stretchr/testify/assert"
)

var update = flag.Bool("update", false, "Updates testdata")

func refStr(s string) *string {
	return &s
}

func copyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}

	destination, err := os.Create(dst)
	if err != nil {
		err2 := source.Close()
		if err2 != nil {
			return config.ErrorCollection([]error{err, err2})
		}

		return err
	}

	errs := make([]error, 0)
	_, err = io.Copy(destination, source)
	if err != nil {
		errs = append(errs, err)
	}

	err2 := source.Close()
	if err2 != nil {
		errs = append(errs, err2)
	}

	err3 := destination.Close()
	if err3 != nil {
		errs = append(errs, err3)
	}

	if len(errs) > 0 {
		return config.ErrorCollection(errs)
	}

	return nil
}

func Test_updateAdminResponse(t *testing.T) {
	if !*update {
		t.SkipNow()
	}

	accessor := viper.NewAccessor(config.Options{
		SearchPaths: []string{filepath.Join("testdata", "config.yaml")},
	})

	ctx := context.Background()
	assert.NoError(t, accessor.UpdateConfig(ctx))
	c := admin.InitializeAdminClient(ctx, *admin.GetConfig(ctx))
	resp, err := c.ListNodeExecutions(ctx, &adminIdl.NodeExecutionListRequest{
		WorkflowExecutionId: &coreIdl.WorkflowExecutionIdentifier{
			Project: "priceoptimizeroffline",
			Domain:  "production",
			Name:    "eqwdb3jwg7",
		},
		Limit: 100,
	})

	assert.NoError(t, err)

	if err != nil {
		t.FailNow()
	}

	m := &jsonpb.Marshaler{}
	var buf bytes.Buffer
	err = m.Marshal(&buf, resp)
	assert.NoError(t, err)
	assert.NoError(t, ioutil.WriteFile(filepath.Join("testdata", "listNodeExecutions.pb"), buf.Bytes(), os.ModePerm))
}

func Test_visualizeTimeline(t *testing.T) {
	ctx := context.Background()

	respBytes, err := ioutil.ReadFile(filepath.Join("testdata", "listNodeExecutions.pb"))
	assert.NoError(t, err)

	resp := &adminIdl.NodeExecutionList{}
	assert.NoError(t, jsonpb.Unmarshal(bytes.NewReader(respBytes), resp))

	m := &mocks.AdminServiceClient{}
	m.OnListNodeExecutionsMatch(mock.Anything, mock.Anything, mock.Anything).Return(resp, nil)
	assert.NoError(t, visualizeTimeline(ctx, m, timelineFlags{
		persistentFlags: persistentFlags{
			Project: refStr("priceoptimizeroffline"),
			Domain:  refStr("production"),
		},
		ExecutionName: refStr("eqwdb3jwg7"),
	}))
}

func Test_visualizeTimeline_test_output(t *testing.T) {
	ctx := context.Background()

	respBytes, err := ioutil.ReadFile(filepath.Join("testdata", "listNodeExecutions.pb"))
	assert.NoError(t, err)

	resp := &adminIdl.NodeExecutionList{}
	assert.NoError(t, jsonpb.Unmarshal(bytes.NewReader(respBytes), resp))

	tmpLoc, err := ioutil.TempFile(os.TempDir(), "visualize_time_line.png")
	assert.NoError(t, err)
	assert.NoError(t, tmpLoc.Close())

	m := &mocks.AdminServiceClient{}
	m.OnListNodeExecutionsMatch(mock.Anything, mock.Anything, mock.Anything).Return(resp, nil)
	assert.NoError(t, visualizeTimeline(ctx, m, timelineFlags{
		persistentFlags: persistentFlags{
			Project: refStr("priceoptimizeroffline"),
			Domain:  refStr("production"),
		},
		ExecutionName: refStr("eqwdb3jwg7"),
		OutputPath:    refStr(tmpLoc.Name()),
	}))

	expectedPath := filepath.Join("testdata", "expected.png")
	if *update {
		assert.NoError(t, copyFile(tmpLoc.Name(), expectedPath))
	}

	expectedBytes, err := ioutil.ReadFile(expectedPath)
	assert.NoError(t, err)

	actualBytes, err := ioutil.ReadFile(tmpLoc.Name())
	assert.NoError(t, err)

	if assert.Equal(t, expectedBytes, actualBytes) {
		assert.NoError(t, os.Remove(tmpLoc.Name()))
	} else {
		t.Logf("Files are different, expected file [%v] vs actual file [%v]", expectedPath, tmpLoc.Name())
	}
}
