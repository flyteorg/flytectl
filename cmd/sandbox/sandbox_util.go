package sandbox

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/enescakir/emoji"
	f "github.com/flyteorg/flytectl/pkg/filesystemutils"
)

var (
	KUBECONFIG     = f.FilePathJoin(f.UserHomeDir(), ".flyte", "kube.yaml")
	FLYTECTLCONFIG = f.FilePathJoin(f.UserHomeDir(), ".flyte", "config.yaml")
)

func SetupFlytectlConfig() error {
	response, err := http.Get(FlytectlConfig)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(FLYTECTLCONFIG, data, 0600)
	if err != nil {
		fmt.Printf("Please create ~/.flyte dir %v \n", emoji.ManTechnologist)
		return err
	}
	return nil
}

func ConfigCleanup() error {
	err := os.Remove(FLYTECTLCONFIG)
	if err != nil {
		return err
	}
	err = os.Remove(KUBECONFIG)
	if err != nil {
		return err
	}
	return nil
}
