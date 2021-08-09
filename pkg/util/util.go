package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/enescakir/emoji"
	"github.com/flyteorg/flytectl/pkg/configutil"
	"github.com/flyteorg/flytectl/pkg/docker"
	f "github.com/flyteorg/flytectl/pkg/filesystemutils"
	hversion "github.com/hashicorp/go-version"
)

const (
	progressSuccessMessage = "Flyte is ready! Flyte UI is available at http://localhost:30081/console"
)

var Ext string

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

func GetRequest(url string) (*http.Response, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("someting goes wrong while sending request to %s. Got status code %v", url, response.StatusCode)
	}
	return response, nil
}
