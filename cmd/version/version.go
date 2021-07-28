package version

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/flyteorg/flytectl/pkg/util"

	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flytestdlib/logger"
	stdlibversion "github.com/flyteorg/flytestdlib/version"
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
	flytectlAppName    = "flytectl"
	controlPlanAppName = "controlPlane"
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

// GetVersionCommand will return version command
func GetVersionCommand(rootCmd *cobra.Command) map[string]cmdCore.CommandEntry {
	getResourcesFuncs := map[string]cmdCore.CommandEntry{
		"version": {CmdFunc: getVersion, Aliases: []string{"versions"}, ProjectDomainNotRequired: true,
			Short: versionCmdShort,
			Long:  versionCmdLong},
	}
	return getResourcesFuncs
}

func getVersion(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {

	message, err := util.GetUpgradeMessage(runtime.GOOS)
	if err != nil {
		logger.Error(ctx, "Not able to detect new version because %v", err)
	}
	if len(message) > 0 {
		fmt.Println(message)
	}

	// Print Flytectl
	if err := printVersion(versionOutput{
		Build:     stdlibversion.Build,
		BuildTime: stdlibversion.BuildTime,
		Version:   stdlibversion.Version,
		App:       flytectlAppName,
	}); err != nil {
		return err
	}
	// Print Flyteadmin version if available
	if err := getControlPlaneVersion(ctx, cmdCtx); err != nil {
		logger.Debug(ctx, err)
	}
	return nil
}

func printVersion(response versionOutput) error {
	b, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return err
	}
	fmt.Print(string(b))
	return nil
}

func getControlPlaneVersion(ctx context.Context, cmdCtx cmdCore.CommandContext) error {
	v, err := cmdCtx.AdminClient().GetVersion(ctx, &admin.GetVersionRequest{})
	if err != nil || v == nil {
		logger.Debugf(ctx, "Failed to get version of control plane %v: \n", err)
		return err
	}
	// Print Flyteadmin
	if err := printVersion(versionOutput{
		Build:     v.ControlPlaneVersion.Build,
		BuildTime: v.ControlPlaneVersion.BuildTime,
		Version:   v.ControlPlaneVersion.Version,
		App:       controlPlanAppName,
	}); err != nil {
		return fmt.Errorf("not able to get control plane version..Please try again: %v", err)
	}
	return nil
}
