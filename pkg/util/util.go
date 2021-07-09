package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	f "github.com/flyteorg/flytectl/pkg/filesystemutils"
)

type githubversion struct {
	TagName string `json:"tag_name"`
}

func GetRequest(baseURL, url string) ([]byte, error) {
	response, err := http.Get(fmt.Sprintf("%v%v", baseURL, url))
	if err != nil {
		return []byte(""), err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte(""), err
	}
	return data, nil
}

func ParseGithubTag(data []byte) (string, error) {
	var result = githubversion{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return "", err
	}
	return result.TagName, nil
}

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
