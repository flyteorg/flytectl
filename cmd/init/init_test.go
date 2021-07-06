package init

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateInitCommand(t *testing.T) {
	initCmd := CreateInitCommand()
	assert.Equal(t, initCmd.Use, "init")
	assert.Equal(t, initCmd.Short, initCmdShort)
	fmt.Println(initCmd.Commands())
	assert.Equal(t, len(initCmd.Commands()), 1)
	cmdNouns := initCmd.Commands()
	// Sort by Use value.
	sort.Slice(cmdNouns, func(i, j int) bool {
		return cmdNouns[i].Use < cmdNouns[j].Use
	})

	assert.Equal(t, cmdNouns[0].Use, "config")
	assert.Equal(t, cmdNouns[0].Short, initCmdShort)
	assert.Equal(t, cmdNouns[0].Long, initCmdLong)
}
