package version

import (
	"context"
	"encoding/json"
	"fmt"

	adminclient "github.com/flyteorg/flyteidl/clients/go/admin"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	adminversion "github.com/flyteorg/flytestdlib/version"
	"github.com/spf13/cobra"
)

// Long descriptions are whitespace sensitive when generating docs using sphinx.
const (
	versionCmdShort = `Used for fetching flyte version`
	versionCmdLong  = `
Example version.
::

 bin/flytectl version
`
)

type versionOutput struct {
	// Specifies the Name of app
	App string `json:"App,omitempty"`
	// Specifies the GIT sha of the build
	Build string `json:"Build,omitempty"`
	// Version for the build, should follow a semver
	Version string `json:"Version,omitempty"`
	// Build timestamp
	BuildTime string `json:"BuildTime,omitempty"`
}

// VersionCommand will return version of flyte
func GetVersionCommand() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:     "version",
		Short:   versionCmdShort,
		Aliases: []string{"versions"},
		Long:    versionCmdLong,
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx := context.Background()
			adminClient, err := adminclient.InitializeAdminClientFromConfig(ctx)
			if err != nil {
				return fmt.Errorf("err %v: ", err)
			}

			v, err := adminClient.GetVersion(ctx, &admin.GetVersionRequest{})
			if err != nil {
				return fmt.Errorf("err %v: ", err)
			}

			// Print Flytectl
			if err := PrintVersion(versionOutput{
				Build:     adminversion.Build,
				BuildTime: adminversion.BuildTime,
				Version:   adminversion.Version,
				App:       "flytectl",
			}); err != nil {
				return err
			}

			// Print Flyteadmin
			if err := PrintVersion(versionOutput{
				Build:     v.ControlPlaneVersion.Build,
				BuildTime: v.ControlPlaneVersion.BuildTime,
				Version:   v.ControlPlaneVersion.Version,
				App:       "controlPlane",
			}); err != nil {
				return err
			}
			return nil
		},
	}

	return versionCmd
}

func PrintVersion(response versionOutput) error {
	b, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return fmt.Errorf("err %v:", err)
	}
	fmt.Print(string(b))
	return nil
}
