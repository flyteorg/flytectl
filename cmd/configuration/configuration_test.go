package configuration

import (
	"context"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"testing"

	initConfig "github.com/flyteorg/flytectl/cmd/config/subcommand/config"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flytectl/pkg/util"

	"github.com/stretchr/testify/assert"
)

func TestCreateInitCommand(t *testing.T) {
	initCmd := CreateInitCommand()
	assert.Equal(t, initCmd.Use, "config")
	assert.Equal(t, initCmd.Short, "Runs various config commands, look at the help of this command to get a list of available commands..")
	fmt.Println(initCmd.Commands())
	assert.Equal(t, len(initCmd.Commands()), 3)
	cmdNouns := initCmd.Commands()
	// Sort by Use value.
	sort.Slice(cmdNouns, func(i, j int) bool {
		return cmdNouns[i].Use < cmdNouns[j].Use
	})

	assert.Equal(t, cmdNouns[0].Use, "discover")
	assert.Equal(t, cmdNouns[0].Short, "Searches for a config in one of the default search paths.")
	assert.Equal(t, cmdNouns[1].Use, "init")
	assert.Equal(t, cmdNouns[1].Short, initCmdShort)
	assert.Equal(t, cmdNouns[2].Use, "validate")
	assert.Equal(t, cmdNouns[2].Short, "Validates the loaded config.")

}

func TestSetupConfigFunc(t *testing.T) {
	var yes = strings.NewReader("Yes")
	var no = strings.NewReader("No")
	mockOutStream := new(io.Writer)
	ctx := context.Background()
	_ = os.Remove(util.FlytectlConfig)

	_ = util.SetupFlyteDir()

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
