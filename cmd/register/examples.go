package register

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	rconfig "github.com/flyteorg/flytectl/cmd/config/subcommand/register"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	f "github.com/flyteorg/flytectl/pkg/filesystemutils"
	"github.com/google/go-github/github"
)

const (
	registerExampleShort = "Registers flytesnack example"
	registerExampleLong  = `
Registers all latest flytesnacks example
::

 bin/flytectl register examples  -d development  -p flytesnacks


Usage
`
)

func registerExampleFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	c := github.NewClient(nil)
	opt := &github.ListOptions{Page: 1, PerPage: 1}
	releases, _, err := c.Repositories.ListReleases(context.Background(), "flyteorg", "flytesnacks", opt)
	if err != nil {
		return err
	}
	response, err := http.Get(fmt.Sprintf("https://github.com/flyteorg/flytesnacks/releases/download/%s/flyte_tests_manifest.json", *releases[0].TagName))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(f.FilePathJoin(f.UserHomeDir(), ".flyte", "flytesnacks.yaml"), data, 0600)
	if err != nil {
		return err
	}

	jsonFile, err := os.Open(f.FilePathJoin(f.UserHomeDir(), ".flyte", "flytesnacks.yaml"))
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteValue, &FlyteSnacksRelease)
	if err != nil {
		return err
	}
	for _, v := range FlyteSnacksRelease {
		rconfig.DefaultFilesConfig.Archive = true
		args := []string{
			fmt.Sprintf("https://github.com/flyteorg/flytesnacks/releases/download/%s/flytesnacks-%s.tgz", *releases[0].TagName, v.Name),
		}
		if err := Register(ctx, args, cmdCtx); err != nil {
			return err
		}
	}
	return nil
}
