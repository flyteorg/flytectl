package githubutil

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	sandboxConfig "github.com/flyteorg/flytectl/cmd/config/subcommand/sandbox"
	"github.com/flyteorg/flytestdlib/logger"
	hversion "github.com/hashicorp/go-version"
	"golang.org/x/oauth2"

	"github.com/flyteorg/flytectl/pkg/util/platformutil"
	stdlibversion "github.com/flyteorg/flytestdlib/version"
	"github.com/mouuff/go-rocket-update/pkg/provider"
	"github.com/mouuff/go-rocket-update/pkg/updater"

	"fmt"

	"github.com/google/go-github/v42/github"
)

const (
	owner                   = "flyteorg"
	flyte                   = "flyte"
	sandboxManifest         = "flyte_sandbox_manifest.yaml"
	flytectl                = "flytectl"
	sandboxSupportedVersion = "v0.10.0"
	flytectlRepository      = "github.com/flyteorg/flytectl"
	commonMessage           = "\n A new release of flytectl is available: %s → %s \n"
	brewMessage             = "To upgrade, run: brew update && brew upgrade flytectl \n"
	linuxMessage            = "To upgrade, run: flytectl upgrade \n"
	darwinMessage           = "To upgrade, run: flytectl upgrade \n"
	releaseURL              = "https://github.com/flyteorg/flytectl/releases/tag/%s \n"
	brewInstallDirectory    = "/Cellar/flytectl"
)

// FlytectlReleaseConfig represent the updater config for flytectl binary
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

//GetGHClient will return github client
func GetGHClient() *github.Client {
	if len(os.Getenv("GITHUB_TOKEN")) > 0 {
		return github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
		)))
	}
	return github.NewClient(&http.Client{})
}

// GetLatestVersion returns the latest version of provided repository
func GetLatestVersion(repository string) (*github.RepositoryRelease, error) {
	client := GetGHClient()
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), owner, repository)
	if err != nil {
		return nil, err
	}
	return release, err
}

// GetListRelease returns the list of release of provided repository
func GetListRelease(repository string) ([]*github.RepositoryRelease, error) {
	client := GetGHClient()
	releases, _, err := client.Repositories.ListReleases(context.Background(), owner, repository, &github.ListOptions{
		PerPage: 100,
	})
	if err != nil {
		return nil, err
	}
	return releases, err
}

// GetSandboxImageSha returns the sha as per input
func GetSandboxImageSha(version string) (string, error) {
	var release *github.RepositoryRelease
	if len(version) == 0 {
		releases, err := GetListRelease(flyte)
		if err != nil {
			return "", err
		}
		for _, v := range releases {
			if *v.Prerelease && sandboxConfig.DefaultConfig.Prerelease {
				logger.Infof(context.Background(), "sandbox started with pre release %s", *v.TagName)
				release = v
				break
			} else if !*v.Prerelease && !sandboxConfig.DefaultConfig.Prerelease {
				logger.Infof(context.Background(), "sandbox started with release %s", *v.TagName)
				release = v
				break
			}
		}
	} else if len(version) > 0 {
		r, err := CheckVersionExist(version, flyte)
		if err != nil {
			return "", err
		}
		release = r
	}
	isGreater, err := IsVersionGreaterThan(*release.TagName, sandboxSupportedVersion)
	if err != nil {
		return "", err
	}
	if !isGreater {
		return "", fmt.Errorf("version flag only supported with flyte %s+ release", sandboxSupportedVersion)
	}
	sha, err := GetSHAFromVersion(*release.TagName, flyte)
	if err != nil {
		return "", err
	}
	return sha, nil
}

func getFlytectlAssetName() string {
	if arch == platformutil.ArchAmd64 {
		arch = platformutil.ArchX86
	} else if arch == platformutil.ArchX86 {
		arch = platformutil.Archi386
	}
	return fmt.Sprintf("flytectl_%s_%s.tar.gz", strings.Title(runtime.GOOS), arch.String())
}

// CheckVersionExist returns the provided version release if version exist in repository
func CheckVersionExist(version, repository string) (*github.RepositoryRelease, error) {
	client := GetGHClient()
	release, _, err := client.Repositories.GetReleaseByTag(context.Background(), owner, repository, version)
	if err != nil {
		return nil, err
	}
	return release, err
}

// GetSHAFromVersion returns sha commit hash against a release
func GetSHAFromVersion(version, repository string) (string, error) {
	client := GetGHClient()
	sha, _, err := client.Repositories.GetCommitSHA1(context.Background(), owner, repository, version, "")
	if err != nil {
		return "", err
	}
	return sha, err
}

// GetAssetsFromRelease returns the asset from github release
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

// GetUpgradeMessage return the upgrade message
func GetUpgradeMessage(latest string, goos platformutil.Platform) (string, error) {
	isGreater, err := IsVersionGreaterThan(latest, stdlibversion.Version)
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

// CheckBrewInstall returns the path of symlink if flytectl is installed from brew
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

// IsVersionGreaterThan check version if it's greater then other
func IsVersionGreaterThan(version1, version2 string) (bool, error) {
	semanticVersion1, err := hversion.NewVersion(version1)
	if err != nil {
		return false, err
	}
	semanticVersion2, err := hversion.NewVersion(version2)
	if err != nil {
		return false, err
	}
	return semanticVersion1.GreaterThan(semanticVersion2), nil
}
