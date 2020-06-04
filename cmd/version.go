package cmd

import (
	"github.com/lyft/flytectl/cmd/config"
	"github.com/lyft/flytestdlib/logger"
	"github.com/lyft/flytestdlib/version"
	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Displays version information for the client and server.",
		Run: func(cmd *cobra.Command, args []string) {
			version.LogBuildInformation("flytectl")
			logger.InfofNoCtx(config.GetConfig().Project)
			logger.InfofNoCtx(config.GetConfig().Domain)
		},
	}
)
