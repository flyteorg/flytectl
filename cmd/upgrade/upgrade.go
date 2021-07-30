package upgrade

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/flyteorg/flytectl/pkg/util/githubutil"

	"github.com/flyteorg/flytestdlib/logger"
	"github.com/mouuff/go-rocket-update/pkg/updater"

	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flytectl/pkg/util"
	"github.com/flyteorg/flytectl/pkg/util/platformutil"
	stdlibversion "github.com/flyteorg/flytestdlib/version"
	"github.com/spf13/cobra"
)

type Goos string

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

Note: Upgrade is not available on windows
`
	subCommand = "rollback"
)

var (
	goos = platformutil.Platform(runtime.GOOS)
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
	var err error
	// Check if it's a rollback
	if len(args) == 1 {
		if args[0] == subCommand && !isRollBackSupported(goos) {
			return nil
		}
		ext, err := githubutil.FlytectlReleaseConfig.GetExecutable()
		if err != nil {
			return err
		}
		backupBinary := fmt.Sprintf("%s.old", ext)
		if _, err := os.Stat(backupBinary); err != nil {
			return fmt.Errorf("flytectl backup doesn't exist at %s. Please first run flytectl upgrade", backupBinary)
		}
		return githubutil.FlytectlReleaseConfig.Rollback()
	}

	latest, err := githubutil.FlytectlReleaseConfig.GetLatestVersion()
	if err != nil {
		return err
	}

	if isGreater, err := util.IsVersionGreaterThan(latest, stdlibversion.Version); err != nil {
		return err
	} else if !isGreater {
		return nil
	}

	if isSupported, err := isUpgradeSupported(latest, goos); err != nil {
		return err
	} else if !isSupported {
		return nil
	}

	if message, err := upgrade(githubutil.FlytectlReleaseConfig); err != nil {
		return err
	} else if len(message) > 0 {
		logger.Info(ctx, message)
	}
	return nil
}

func upgrade(u *updater.Updater) (string, error) {
	updateStatus, err := u.Update()
	if err != nil {
		return "", err
	}

	if updateStatus == updater.Updated {
		latestVersion, err := u.GetLatestVersion()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Successfully updated to version %s", latestVersion), nil
	}
	return "", u.Rollback()
}

func isUpgradeSupported(latest string, goos platformutil.Platform) (bool, error) {
	message, err := githubutil.GetUpgradeMessage(latest, goos)
	if err != nil {
		return false, err
	}
	if goos.String() == platformutil.Windows.String() || strings.Contains(message, "brew") {
		if len(message) > 0 {
			fmt.Println(message)
		}
		return false, nil
	}
	return true, nil
}

func isRollBackSupported(goos platformutil.Platform) bool {
	if goos.String() == platformutil.Windows.String() {
		fmt.Printf("Flytectl rollback is not available on %s \n", goos.String())
		return false
	}
	return true
}
