package init

import (
	"context"
	"io"
	"strings"
	"testing"

	initConfig "github.com/flyteorg/flytectl/cmd/config/subcommand/init"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/stretchr/testify/assert"
)

func TestSetupConfig(t *testing.T) {
	assert.Nil(t, initFlytectlConfig(strings.NewReader("Yes")))
	assert.Nil(t, initFlytectlConfig(strings.NewReader("Yes")))
	assert.Nil(t, initFlytectlConfig(strings.NewReader("No")))
	initConfig.DefaultConfig.Host = "test"
	assert.Nil(t, initFlytectlConfig(strings.NewReader("Yes")))
}

func TestSetupConfigFunc(t *testing.T) {
	mockOutStream := new(io.Writer)
	ctx := context.Background()
	cmdCtx := cmdCore.NewCommandContext(nil, *mockOutStream)
	err := configInitFunc(ctx, []string{}, cmdCtx)
	assert.Nil(t, err)
}
