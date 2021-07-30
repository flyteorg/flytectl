package githubutil

import (
	"context"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/flyteorg/flytectl/pkg/util/platformutil"
	stdlibversion "github.com/flyteorg/flytestdlib/version"
	"github.com/mouuff/go-rocket-update/pkg/provider"
	"github.com/mouuff/go-rocket-update/pkg/updater"

	"github.com/flyteorg/flytectl/pkg/util"

	"fmt"
	"io/ioutil"

	"github.com/google/go-github/v37/github"
)

const (
	owner                = "flyteorg"
	flyte                = "flyte"
	sandboxManifest      = "flyte_sandbox_manifest.yaml"
	flytectl             = "flytectl"
	flytectlRepository   = "github.com/flyteorg/flytectl"
	commonMessage        = "\n A new release of flytectl is available: %s → %s \n"
	brewMessage          = "To upgrade, run: brew update && brew upgrade flytectl \n"
	linuxMessage         = "To upgrade, run: flytectl upgrade \n"
	darwinMessage        = "To upgrade, run: flytectl upgrade \n"
	releaseURL           = "https://github.com/flyteorg/flytectl/releases/tag/%s \n"
	brewInstallDirectory = "/Cellar/flytectl"
)

var FlytectlReleaseConfig = &updater.Updater{
	Provider: &provider.Github{
		RepositoryURL: flytectlRepository,
		ArchiveName:   getFlytectlAssetName(),
	},
	ExecutableName: flytectl,
	Version:        stdlibversion.Version,
}

var (
	arch = platformutil.Arch(runtime.GOARCH)
)

func GetLatestVersion(repository string) (*github.RepositoryRelease, error) {
	client := github.NewClient(nil)
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), owner, repository)
	if err != nil {
		return nil, err
	}
	return release, err
}

func getFlytectlAssetName() string {
	if arch == platformutil.ArchAmd64 {
		arch = platformutil.ArchX86
	} else if arch == platformutil.ArchX86 {
		arch = platformutil.Archi386
	}
	return fmt.Sprintf("flytectl_%s_%s.tar.gz", strings.Title(runtime.GOOS), arch.String())
}

func CheckVersionExist(version, repository string) (*github.RepositoryRelease, error) {
	client := github.NewClient(nil)
	release, _, err := client.Repositories.GetReleaseByTag(context.Background(), owner, repository, version)
	if err != nil {
		return nil, err
	}
	return release, err
}

func GetAssetsFromRelease(version, assets, repository string) (*github.ReleaseAsset, error) {
	release, err := CheckVersionExist(version, repository)
	if err != nil {
		return nil, err
	}
	for _, v := range release.Assets {
		if v.GetName() == assets {
			return v, nil
		}
	}
	return nil, fmt.Errorf("assest is not found in %s[%s] release", repository, version)
}

func GetFlyteManifest(version string, target string) error {
	asset, err := GetAssetsFromRelease(version, sandboxManifest, flyte)
	if err != nil {
		return err
	}
	response, err := util.GetRequest(asset.GetBrowserDownloadURL())
	if err != nil {
		return err
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if err := util.WriteIntoFile(data, target); err != nil {
		return err
	}
	return nil

}

func GetUpgradeMessage(latest string, goos platformutil.Platform) (string, error) {
	isGreater, err := util.IsVersionGreaterThan(latest, stdlibversion.Version)
	if err != nil {
		return "", err
	}
	message := fmt.Sprintf(commonMessage, stdlibversion.Version, latest)
	if isGreater {
		symlink, err := CheckBrewInstall(goos)
		if err != nil {
			return "", err
		}
		if len(symlink) > 0 {
			message += brewMessage
		} else if goos == platformutil.Darwin {
			message += darwinMessage
		} else if goos == platformutil.Linux {
			message += linuxMessage
		}
		message += fmt.Sprintf(releaseURL, latest)
	}

	return message, nil
}

func CheckBrewInstall(goos platformutil.Platform) (string, error) {
	if goos.String() == platformutil.Darwin.String() {
		executable, err := FlytectlReleaseConfig.GetExecutable()
		if err != nil {
			return executable, err
		}
		if symlink, err := filepath.EvalSymlinks(executable); err != nil {
			return symlink, err
		} else if len(symlink) > 0 {
			if strings.Contains(symlink, brewInstallDirectory) {
				return symlink, nil
			}
		}
	}
	return "", nil
}
