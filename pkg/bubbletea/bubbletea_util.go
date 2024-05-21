package bubbletea

import (
	"os"

	"github.com/spf13/cobra"
)

// Check if -f bubbletea is in args
func ifRunBubbleTea(_rootCmd cobra.Command) (*cobra.Command, bool, error) {
	cmd, flags, err := _rootCmd.Find(os.Args[1:])
	if err != nil {
		return cmd, false, err
	}

	tempCmd := cmd
	for tempCmd.HasParent() {
		newArgs = append([]string{tempCmd.Use}, newArgs...)
		tempCmd = tempCmd.Parent()
	}

	for _, flag := range flags {
		if flag == "-i" || flag == "--interactive" {
			return cmd, true, nil
		}
	}

	return cmd, false, nil
	// err = _rootCmd.ParseFlags(flags)
	// if err != nil {
	// 	return nil, false, err
	// }

	// format, err := _rootCmd.Flags().GetString("format")
	// if format != "bubbletea" || err != nil {
	// 	return nil, false, err
	// } else {
	// 	return cmd, true, err
	// }
}
