package version

import (
	"context"
	cmdCore "github.com/lyft/flytectl/cmd/core"
	cmdcore "github.com/lyft/flytectl/cmd/core"
	"github.com/lyft/flytectl/pkg/version"
	"github.com/spf13/cobra"
)

func getVersion(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	version.LogBuildInformation("flytectl")
	return nil
}

// CreateVersionCommand will return get command
func CreateVersionCommand() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Retrieve version of the flytectl.",
	}

	versionResourcesFuncs := map[string]cmdcore.CommandEntry{
		"projects": {CmdFunc: getVersion, ProjectDomainNotRequired: true},
	}

	cmdcore.AddCommands(versionCmd, versionResourcesFuncs)

	return versionCmd
}
