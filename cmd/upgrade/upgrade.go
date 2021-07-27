package upgrade

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/avast/retry-go"

	"github.com/google/go-github/v37/github"

	cmdUtil "github.com/flyteorg/flytectl/pkg/commandutils"
	f "github.com/flyteorg/flytectl/pkg/filesystemutils"
	"github.com/flyteorg/flytectl/pkg/util"
	stdlibversion "github.com/flyteorg/flytestdlib/version"

	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/spf13/cobra"
)

// Long descriptions are whitespace sensitive when generating docs using sphinx.
const (
	upgradeCmdShort = `Used for upgrade/rollback flyte version`
	upgradeCmdLong  = `
Upgrade flytectl
::

 bin/flytectl upgrade
	
Rollback flytectl binary
::

 bin/flytectl upgrade rollback	
`
)

var (
	target = f.FilePathJoin(f.UserHomeDir(), ".flyte", "flytectl")
	backup = f.FilePathJoin(f.UserHomeDir(), ".flyte", "flytectl.bak")
	ext    = "/bin/flytectl"
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
	var source, message string
	// Get the binary path
	if err := getExecutable(); err != nil {
		return err
	}
	// Check if it's a rollback
	if len(args) > 0 {
		if args[0] == "rollback" {
			// Restore the binary
			source = backup
		}
	} else {
		latest, err := util.GetLatestVersion("flytectl")
		if err != nil {
			return err
		}

		isGreater, err := util.IsVersionGreaterThan(latest.GetTagName(), stdlibversion.Version)
		if err != nil {
			return err
		}
		if !isGreater {
			return nil
		}
		if err := upgrade(os.Stdin, latest, ext, runtime.GOARCH); err != nil {
			return err
		}
		source = target
		message = fmt.Sprintf("Successfully updated to version %s", latest.GetTagName())
	}

	if err := move(source, ext); err != nil {
		return err
	}

	if err := os.Chmod(ext, 0111); err != nil {
		return err
	}

	if len(message) > 0 {
		fmt.Println(message)
	}
	return nil
}

func upgrade(reader io.Reader, latest *github.RepositoryRelease, ext, goos string) error {
	// Check if GOOS is windows
	if checkWindows(latest.GetTagName(), goos) {
		return nil
	}

	if cmdUtil.AskForConfirmation(fmt.Sprintf("This action will upgrade flytectl version [%s]. Do you want to continue?", latest.GetTagName()), reader) {
		arch := runtime.GOARCH
		if arch == "amd64" {
			arch = "x86_64"
		} else if arch == "386" {
			arch = "i386"
		}
		for _, v := range latest.Assets {
			asset := fmt.Sprintf("flytectl_%s_%s.tar.gz", strings.Title(goos), arch)
			if v.GetName() == asset {
				response, err := http.Get(*v.BrowserDownloadURL)
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
						// Backup the binary
						if err := copy(ext, backup); err != nil {
							return err
						}
						break
					}
				}
			}
		}
	}
	return nil
}

func checkWindows(latest, goos string) bool {
	if goos == "windows" {
		fmt.Printf("A new release of flytectl is available: %s → %s \n", stdlibversion.Version, latest)
		fmt.Println("Flytectl auto upgrade is not available on windows")
		fmt.Printf("https://github.com/flyteorg/flytectl/releases/tag/%s \n", latest)
		return true
	}
	return false
}

func move(source, destination string) error {
	err := retry.Do(func() error {
		err := os.Rename(source, destination)
		if err != nil {
			return err
		}
		return nil
	}, retry.Attempts(3))
	if err != nil {
		return err
	}
	return nil
}

func copy(source, destination string) error {
	err := retry.Do(func() error {
		data, err := ioutil.ReadFile(source)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(destination, data, 0600)
		if err != nil {
			return err
		}
		return nil
	}, retry.Attempts(3))
	if err != nil {
		return err
	}
	return nil
}

func getExecutable() error {
	if len(ext) == 0 {
		binary, err := os.Executable()
		if err != nil {
			return err
		}
		ext = binary
	}
	return nil
}
