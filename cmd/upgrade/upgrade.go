package upgrade

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	cmdUtil "github.com/flyteorg/flytectl/pkg/commandutils"
	f "github.com/flyteorg/flytectl/pkg/filesystemutils"
	"github.com/flyteorg/flytectl/pkg/util"
	stdlibversion "github.com/flyteorg/flytestdlib/version"

	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/spf13/cobra"
)

// Long descriptions are whitespace sensitive when generating docs using sphinx.
const (
	upgradeCmdShort = `Used for fetching flyte version`
	upgradeCmdLong  = `
Example version.
::

 bin/flytectl version
`
)

// SelfUpgrade will return self upgrade command
func SelfUpgrade(rootCmd *cobra.Command) map[string]cmdCore.CommandEntry {
	getResourcesFuncs := map[string]cmdCore.CommandEntry{
		"upgrade": {CmdFunc: selfUpgrade, Aliases: []string{"upgrade"}, ProjectDomainNotRequired: true,
			Short: upgradeCmdShort,
			Long:  upgradeCmdLong},
	}
	return getResourcesFuncs
}

func selfUpgrade(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	latest, err := util.GetLatestVersion(util.FlytectlReleasePath)
	if err != nil {
		return err
	}
	if CheckWindows(latest, "window") {
		return nil
	}

	isGreater, err := util.IsVersionGreaterThan(latest, stdlibversion.Version)
	if err != nil {
		return err
	}
	ext, err := os.Executable()
	if err != nil {
		return err
	}

	if isGreater {
		return upgrade(os.Stdin, latest, ext)
	}
	return nil
}

func upgrade(reader io.Reader, latest, ext string) error {
	if cmdUtil.AskForConfirmation(fmt.Sprintf("This action will upgrade flytectl version [%s]. Do you want to continue?", latest), reader) {
		arch := runtime.GOARCH
		if arch == "amd64" {
			arch = "x86_64"
		} else if arch == "386" {
			arch = "i386"
		}
		assetURL := fmt.Sprintf("/flyteorg/flytectl/releases/download/%s/flytectl_%s_%s.tar.gz", latest, strings.Title(runtime.GOOS), arch)

		response, err := util.GetRequest("https://github.com", assetURL)
		if err != nil {
			return err
		}
		gzr, err := gzip.NewReader(response.Body)
		if err != nil {
			return err
		}
		defer gzr.Close()

		tr := tar.NewReader(gzr)
		for {
			header, err := tr.Next()

			switch {
			case err == io.EOF:
				return nil

			case err != nil:
				return err

			case header == nil:
				continue
			}
			if header.Name == "flytectl" {

				target := f.FilePathJoin(f.UserHomeDir(), ".flyte", "flytectl")
				f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
				if err != nil {
					return err
				}
				for {
					_, err := io.CopyN(f, tr, 1024)
					if err != nil {
						if err == io.EOF {
							break
						}
						return err
					}
				}

				err = os.Rename(target, ext)
				if err != nil {
					rerr := os.Rename(target, ext)
					if rerr != nil {
						return nil
					}
					return err
				}
				fmt.Printf("Successfully updated to version %s", latest)
			}
		}
	}
	return nil
}

func CheckWindows(latest, goos string) bool {
	if runtime.GOOS == goos {
		fmt.Printf("A new release of flytectl is available: %s → %s \n", stdlibversion.Version, latest)
		fmt.Println("Flytectl auto upgrade is not available on windows")
		fmt.Printf("https://github.com/flyteorg/flytectl/releases/tag/%s \n", latest)
		return true
	}
	return false
}
