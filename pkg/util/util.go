package util

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	stdlibversion "github.com/flyteorg/flytestdlib/version"

	"github.com/enescakir/emoji"
	"github.com/flyteorg/flytectl/pkg/configutil"
	"github.com/flyteorg/flytectl/pkg/docker"
	"github.com/google/go-github/v37/github"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"

	f "github.com/flyteorg/flytectl/pkg/filesystemutils"
	hversion "github.com/hashicorp/go-version"
)

const (
	progressSuccessMessage = "Flyte is ready! Flyte UI is available at http://localhost:30081/console"
	FlytectlReleasePath    = "/repos/flyteorg/flytectl/releases/latest"
	commonMessage          = "\n A new release of flytectl is available: %s → %s \n"
	darwinMessage          = "To upgrade, run: brew update && brew upgrade flytectl \n"
	linuxMessage           = "To upgrade, run: flytectl upgrade \n"
	releaseURL             = "https://github.com/flyteorg/flytectl/releases/tag/%s \n"
	owner                  = "flyteorg"
)

func WriteIntoFile(data []byte, file string) error {
	err := ioutil.WriteFile(file, data, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// SetupFlyteDir will create .flyte dir if not exist
func SetupFlyteDir() error {
	if err := os.MkdirAll(f.FilePathJoin(f.UserHomeDir(), ".flyte"), os.ModePerm); err != nil {
		return err
	}
	return nil
}

func IsVersionGreaterThan(version1, version2 string) (bool, error) {
	semanticVersion1, err := hversion.NewVersion(version1)
	if err != nil {
		return false, err
	}
	semanticVersion2, err := hversion.NewVersion(version2)
	if err != nil {
		return false, err
	}
	return semanticVersion2.LessThanOrEqual(semanticVersion1), nil
}

func PrintSandboxMessage() {
	kubeconfig := strings.Join([]string{
		"$KUBECONFIG",
		f.FilePathJoin(f.UserHomeDir(), ".kube", "config"),
		docker.Kubeconfig,
	}, ":")

	fmt.Printf("%v %v %v %v %v \n", emoji.ManTechnologist, progressSuccessMessage, emoji.Rocket, emoji.Rocket, emoji.PartyPopper)
	fmt.Printf("Add KUBECONFIG and FLYTECTL_CONFIG to your environment variable \n")
	fmt.Printf("export KUBECONFIG=%v \n", kubeconfig)
	fmt.Printf("export FLYTECTL_CONFIG=%v \n", configutil.FlytectlConfig)
}

func ProgressBarForFlyteStatus(total int64, count chan int64, message string) {
	p := mpb.New(mpb.WithWidth(64))
	bar := p.AddBar(total,
		mpb.PrependDecorators(
			decor.Name(message, decor.WC{W: len(message) + 1, C: decor.DidentRight}),
			decor.Name("", decor.WCSyncSpaceR),
		),
		mpb.AppendDecorators(
			decor.OnComplete(decor.Percentage(decor.WC{W: 5}), ""),
		),
		mpb.AppendDecorators(decor.Percentage()),
	)
	for {
		c := <-count
		bar.IncrBy(int(c))
		if bar.Current() == total {
			PrintSandboxMessage()
			return
		}
	}
}

func GetUpgradeMessage(goos string) (string, error) {
	latest, err := GetLatestVersion("flytectl")
	if err != nil {
		return "", err
	}
	isGreater, err := IsVersionGreaterThan(latest.GetTagName(), stdlibversion.Version)
	if err != nil {
		return "", err
	}
	message := ""
	if isGreater {
		if goos == "darwin" {
			message = commonMessage + darwinMessage
		} else if goos == "linux" {
			message = commonMessage + linuxMessage
		}
		message += releaseURL
		message = fmt.Sprintf(message, stdlibversion.Version, latest.GetTagName(), latest.GetTagName())
	}

	return message, nil
}

func GetLatestVersion(repository string) (*github.RepositoryRelease, error) {
	client := github.NewClient(nil)
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), owner, repository)
	if err != nil {
		return nil, err
	}
	return release, err
}

func CheckVersionExist(version, repository string) (*github.RepositoryRelease, error) {
	client := github.NewClient(nil)
	release, _, err := client.Repositories.GetReleaseByTag(context.Background(), owner, repository, version)
	if err != nil {
		return nil, err
	}
	return release, err
}
