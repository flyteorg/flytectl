package bubbletea

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/spf13/cobra"
)

type Command struct {
	Cmd   *cobra.Command
	Name  string
	Short string
}

var (
	nameToCommand = map[string]Command{}
	rootCmd       *cobra.Command
	targetArgs    []string
)

func generateSubCmdItems(cmd *cobra.Command) []list.Item {
	items := []list.Item{}

	for _, subcmd := range cmd.Commands() {
		nameToCommand[subcmd.Use] = Command{
			Cmd:   subcmd,
			Name:  subcmd.Use,
			Short: subcmd.Short,
		}
		items = append(items, item(subcmd.Use))
	}

	return items
}

// func isValidCommand(curArg string, cmd *cobra.Command) (*cobra.Command, bool) {
// 	for _, subCmd := range cmd.Commands() {
// 		if subCmd.Use == curArg {
// 			return subCmd, true
// 		}
// 	}
// 	return nil, false
// }

// func findSubCmdItems(cmd *cobra.Command, inputArgs []string) ([]list.Item, error) {
// 	if len(inputArgs) == 0 {
// 		return generateSubCmdItems(cmd), nil
// 	}

// 	curArg := inputArgs[0]
// 	subCmd, isValid := isValidCommand(curArg, cmd)
// 	if !isValid {
// 		return nil, fmt.Errorf("not a valid argument: %v", curArg)
// 	}

// 	return findSubCmdItems(subCmd, inputArgs[1:])
// }

func newList(i string) (list.Model, bool) {

	items := []list.Item{}

	if len(nameToCommand[i].Cmd.Commands()) == 0 {
		return list.New(items, itemDelegate{}, defaultWidth, listHeight), true
	}

	items = generateSubCmdItems(nameToCommand[i].Cmd)
	l := genList(items)

	return l, false
}
