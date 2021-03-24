package cmd

import (
	"context"
	"fmt"

	"github.com/flyteorg/flytectl/version"
	"github.com/flyteorg/flyteidl/clients/go/admin"
	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Displays version information for the client and server.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			adminClient, err := admin.InitializeAdminClientFromConfig(ctx)
			if err != nil {
				fmt.Sprintf("err %v:", err)
				return
			}
			version.LogBuildInformation("flytectl")
			// TODO: Log Admin version
			v, err := adminClient.GetVersion(ctx, &admin.GetVersionRequest{})
			if err != nil {
				fmt.Sprintf("err %v:", err)
				return
			}
			version.PrintVersion("flyteadmin", v)

		},
	}
)
