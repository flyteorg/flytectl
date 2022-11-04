package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/flyteorg/flytectl/pkg/docker"
	f "github.com/flyteorg/flytectl/pkg/filesystemutils"

	"github.com/enescakir/emoji"
	hversion "github.com/hashicorp/go-version"
)

const (
	ProgressSuccessMessage        = "Flyte is ready! Flyte UI is available at"
	ProgressSuccessMessagePending = "Flyte would be ready after this! Flyte UI would be available at"
	SandBoxConsolePort            = 30081
	DemoConsolePort               = 30080
)

var Ext string

// WriteIntoFile will write content in a file
func WriteIntoFile(data []byte, file string) error {
	err := ioutil.WriteFile(file, data, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// SetupFlyteDir will create .flyte dir if not exist
func SetupFlyteDir() error {
	if err := os.MkdirAll(f.FilePathJoin(f.UserHomeDir(), ".flyte", "k3s"), os.ModePerm); err != nil {
		return err
	}

	// Created a empty file with right permission
	if _, err := os.Stat(docker.Kubeconfig); err != nil {
		if os.IsNotExist(err) {
			if err := ioutil.WriteFile(docker.Kubeconfig, []byte(""), os.ModePerm); err != nil {
				return err
			}
		}
	}

	return nil
}

// PrintSandboxMessage will print sandbox success message
func PrintSandboxMessage(flyteConsolePort int, dryRun bool) {
	kubeconfig := strings.Join([]string{
		"$KUBECONFIG",
		f.FilePathJoin(f.UserHomeDir(), ".kube", "config"),
		docker.Kubeconfig,
	}, ":")

	var successMsg string
	if dryRun {
		successMsg = fmt.Sprintf("%v http://localhost:%v/console", ProgressSuccessMessagePending, flyteConsolePort)
	} else {
		successMsg = fmt.Sprintf("%v http://localhost:%v/console", ProgressSuccessMessage, flyteConsolePort)

	}
	fmt.Printf("%v %v %v %v %v \n", emoji.ManTechnologist, successMsg, emoji.Rocket, emoji.Rocket, emoji.PartyPopper)
	fmt.Printf("%v Run the following command to export sandbox environment variables for accessing flytectl\n", emoji.Sparkle)
	fmt.Printf("	Add FLYTECTL_CONFIG to your environment variable \n")
	if dryRun {
		fmt.Printf("%v Run the following command to export kubeconfig variables for accessing flyte pods locally\n", emoji.Sparkle)
		fmt.Printf("	export KUBECONFIG=%v \n", kubeconfig)
	}
}

// SendRequest will create request and return the response
func SendRequest(method, url string, option io.Reader) (*http.Response, error) {
	client := &http.Client{}
	req, _ := http.NewRequest(method, url, option)
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("someting goes wrong while sending request to %s. Got status code %v", url, response.StatusCode)
	}
	return response, nil
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
