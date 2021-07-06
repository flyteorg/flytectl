package init

import (
	"context"
	"io"
	"os"
	"strings"

	initConfig "github.com/flyteorg/flytectl/cmd/config/subcommand/init"
	"github.com/flyteorg/flytectl/pkg/util"

	cmdCore "github.com/flyteorg/flytectl/cmd/core"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupConfigFunc(t *testing.T) {
	var yes = strings.NewReader("Yes")
	var no = strings.NewReader("No")
	mockOutStream := new(io.Writer)
	ctx := context.Background()
	_ = os.Remove(util.FlytectlConfig)

	cmdCtx := cmdCore.NewCommandContext(nil, *mockOutStream)
	err := configInitFunc(ctx, []string{}, cmdCtx)
	initConfig.DefaultConfig.Host = ""
	assert.Nil(t, err)

	assert.Nil(t, initFlytectlConfig(yes))
	assert.Nil(t, initFlytectlConfig(yes))
	assert.Nil(t, initFlytectlConfig(no))
	initConfig.DefaultConfig.Host = "test"
	assert.NotNil(t, initFlytectlConfig(no))
}
