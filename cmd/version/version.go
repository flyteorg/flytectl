package version

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flytestdlib/logger"
	stdlibversion "github.com/flyteorg/flytestdlib/version"
	hversion "github.com/hashicorp/go-version"
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

var githubBaseURL = "https://api.github.com"

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
	latest, err := getLatestVersion(githubBaseURL)
	if err != nil {
		return fmt.Errorf("err %v: ", err)
	}

	fmt.Println(compareVersion(latest, stdlibversion.Version))
	// Print Flytectl
	if err := printVersion(versionOutput{
		Build:     stdlibversion.Build,
		BuildTime: stdlibversion.BuildTime,
		Version:   stdlibversion.Version,
		App:       "flytectl",
	}); err != nil {
		return err
	}

	getControlePlaneVersion(ctx, cmdCtx)
	return nil
}

func printVersion(response versionOutput) error {
	b, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return fmt.Errorf("err %v: ", err)
	}
	fmt.Print(string(b))
	return nil
}

func compareVersion(latest, current string) string {
	latestVersion, _ := hversion.NewVersion(latest)
	fversion, _ := hversion.NewVersion(current)

	if fversion.LessThan(latestVersion) {
		return fmt.Sprintf("A newer version of flytectl is available [%v] Please upgrade using - https://docs.flyte.org/projects/flytectl/en/latest/index.html", latest)
	}

	return "Installed flytectl version is the latest"
}

func getControlePlaneVersion(ctx context.Context, cmdCtx cmdCore.CommandContext) {
	v, err := cmdCtx.AdminClient().GetVersion(ctx, &admin.GetVersionRequest{})
	if err != nil || v == nil {
		logger.Debugf(ctx, "Failed to get version of control plane %v: \n", err)
		return
	}
	// Print Flyteadmin
	if err := printVersion(versionOutput{
		Build:     v.ControlPlaneVersion.Build,
		BuildTime: v.ControlPlaneVersion.BuildTime,
		Version:   v.ControlPlaneVersion.Version,
		App:       "controlPlane",
	}); err != nil {
		logger.Debugf(ctx, "Not able to get control plane version..Please try again : %v \n", err)
	}
}

func getLatestVersion(baseURL string) (string, error) {
	response, err := http.Get(fmt.Sprintf("%v/repos/flyteorg/flytectl/releases/latest", baseURL))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	var result = make(map[string]interface{})
	err = json.Unmarshal(data, &result)
	if err != nil {
		return "", err
	}
	return result["tag_name"].(string), nil
}
