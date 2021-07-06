package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"

	f "github.com/flyteorg/flytectl/pkg/filesystemutils"
)

type githubversion struct {
	TagName string `json:"tag_name"`
}

const (
	AdminConfigTemplate = `admin:
  # For GRPC endpoints you might want to use dns:///flyte.myexample.com
  endpoint: {{.Host}}
  authType: Pkce
  insecure: {{.Insecure}}
logger:
  show-source: true
  level: 0`
	StorageConfigTemplate = `
storage:
  connection:
    access-key: minio
    auth-type: accesskey
    disable-ssl: true
    endpoint: http://localhost:30084
    region: us-east-1
    secret-key: miniostorage
  type: minio
  container: "my-s3-bucket"
  enable-multicontainer: true`
	StorageS3ConfigTemplate = `
#storage:
#  kind: s3
#  config:
#    auth_type: iam
#    region: <replace> # Example: us-east-2
#  container: <replace> # Example my-bucket. Flyte k8s cluster / service account for execution should have read access to this bucket`
	StorageGCSConfigTemplate = `
#storage:
#  kind: google
#  config:
#    json: ""
#    project_id: <replace-me> # TODO: replace <project-id> with the GCP project ID
#    scopes: https://www.googleapis.com/auth/devstorage.read_write
#  container: <replace> # Example my-bucket. Flyte k8s cluster / service account for execution should have read access to this bucket`
)

type ConfigTemplateValuesSpec struct {
	Host     string
	Insecure bool
	Template string
}

var (
	FlytectlConfig = f.FilePathJoin(f.UserHomeDir(), ".flyte", "config-sandbox.yaml")
	ConfigFile     = f.FilePathJoin(f.UserHomeDir(), ".flyte", "config.yaml")
	Kubeconfig     = f.FilePathJoin(f.UserHomeDir(), ".flyte", "k3s", "k3s.yaml")
)

// GetSandboxTemplate return sandbox cluster config
func GetSandboxTemplate() string {
	return AdminConfigTemplate + StorageConfigTemplate
}

// GetAWSCloudTemplate return aws flyte config with storage config
func GetAWSCloudTemplate() string {
	return AdminConfigTemplate + StorageS3ConfigTemplate
}

// GetGoogleCloudTemplate return google flyte config with storage config
func GetGoogleCloudTemplate() string {
	return AdminConfigTemplate + StorageGCSConfigTemplate
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
	err := ioutil.WriteFile(file, data, 0600)
	if err != nil {
		return err
	}
	return nil
}

// SetupFlyteDir will create .flyte dir if not exist
func SetupFlyteDir() error {
	if err := os.MkdirAll(f.FilePathJoin(f.UserHomeDir(), ".flyte"), 0755); err != nil {
		return err
	}
	return nil
}

// SetupConfig download the flyte sandbox config
func SetupConfig(filename string, templateValue ConfigTemplateValuesSpec) error {
	tmpl := template.New("config")
	tmpl, err := tmpl.Parse(templateValue.Template)
	if err != nil {
		return err
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return tmpl.Execute(file, templateValue)
}

// ConfigCleanup will remove the sandbox config from flyte dir
func ConfigCleanup() error {
	err := os.Remove(FlytectlConfig)
	if err != nil {
		return err
	}
	err = os.RemoveAll(f.FilePathJoin(f.UserHomeDir(), ".flyte", "k3s"))
	if err != nil {
		return err
	}
	return nil
}
