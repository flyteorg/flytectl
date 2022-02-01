package configuration

import (
	"context"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/flyteorg/flytectl/pkg/configutil"

	initConfig "github.com/flyteorg/flytectl/cmd/config/subcommand/config"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flytectl/pkg/util"

	"github.com/stretchr/testify/assert"
)

func TestCreateInitCommand(t *testing.T) {
	configCmd := CreateConfigCommand()
	assert.Equal(t, configCmd.Use, "config")
	assert.Equal(t, configCmd.Short, "Run various config commands, look at the help of this command to get a list of available commands.")
	fmt.Println(configCmd.Commands())
	assert.Equal(t, len(configCmd.Commands()), 3)
	cmdNouns := configCmd.Commands()
	// Sort by Use value.
	sort.Slice(cmdNouns, func(i, j int) bool {
		return cmdNouns[i].Use < cmdNouns[j].Use
	})

	assert.Equal(t, cmdNouns[0].Use, "discover")
	assert.Equal(t, cmdNouns[0].Short, "Search for a config in one of the default search paths.")
	assert.Equal(t, cmdNouns[1].Use, "init")
	assert.Equal(t, cmdNouns[1].Short, initCmdShort)
	assert.Equal(t, cmdNouns[2].Use, "validate")
	assert.Equal(t, cmdNouns[2].Short, "Validate the loaded config.")

}

func TestSetupConfigFunc(t *testing.T) {
	var yes = strings.NewReader("Yes")
	var no = strings.NewReader("No")
	mockOutStream := new(io.Writer)
	ctx := context.Background()
	_ = os.Remove(configutil.FlytectlConfig)

	_ = util.SetupFlyteDir()

	cmdCtx := cmdCore.NewCommandContext(nil, *mockOutStream)
	err := configInitFunc(ctx, []string{}, cmdCtx)
	initConfig.DefaultConfig.Host = ""
	assert.Nil(t, err)

	assert.Nil(t, initFlytectlConfig(ctx, yes))
	assert.Nil(t, initFlytectlConfig(ctx, yes))
	assert.Nil(t, initFlytectlConfig(ctx, no))
	initConfig.DefaultConfig.Host = "flyte.org"
	assert.Nil(t, initFlytectlConfig(ctx, no))
	initConfig.DefaultConfig.Host = "localhost:30081"
	assert.Nil(t, initFlytectlConfig(ctx, no))
	initConfig.DefaultConfig.Storage = true
	assert.NotNil(t, initFlytectlConfig(ctx, yes))
}

func TestTrimFunc(t *testing.T) {
	assert.Equal(t, trimEndpoint("dns://localhost"), "localhost")
	assert.Equal(t, trimEndpoint("http://localhost"), "localhost")
	assert.Equal(t, trimEndpoint("https://localhost"), "localhost")
}

func TestValidateEndpointName(t *testing.T) {
	assert.Equal(t, true, validateEndpointName("8093405779.ap-northeast-2.elb.amazonaws.com:81"))
	assert.Equal(t, true, validateEndpointName("8093405779.ap-northeast-2.elb.amazonaws.com"))
	assert.Equal(t, false, validateEndpointName("8093405779.ap-northeast-2.elb.amazonaws.com:81/console"))
	assert.Equal(t, true, validateEndpointName("localhost"))
	assert.Equal(t, true, validateEndpointName("127.0.0.1"))
	assert.Equal(t, true, validateEndpointName("127.0.0.1:30086"))
	assert.Equal(t, true, validateEndpointName("112.11.1.1"))
	assert.Equal(t, true, validateEndpointName("112.11.1.1:8080"))
	assert.Equal(t, false, validateEndpointName("112.11.1.1:8080/console"))
	assert.Equal(t, false, validateEndpointName("flyte"))
}
