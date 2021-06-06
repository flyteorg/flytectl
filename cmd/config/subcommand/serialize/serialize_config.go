package serialize

//go:generate pflags Config --default-var DefaultFilesConfig --bind-default-var

var (
	DefaultFilesConfig = &Config{}
)

// Config containing flags used for serialize
type Config struct {
	Image            string `json:"image" pflag:", Docker image name"`
	Registry         string `json:"registry" pflag:", Docker image name"`
	ServiceAccount   string `json:"service-account" pflag:",Service account name."`
	Version          string `json:"version" pflag:", workflow version"`
	OutputDir        string `json:"output-dir" pflag:", custom output location "`
	OutputDirprefix  string `json:"output-dir-prefix" pflag:", Container registry "`
	FlyteAwsEndpoint string `json:"flyte-aws-endpoint" pflag:", Container registry "`
	FlyteAwsKey      string `json:"flyte-aws-key" pflag:", Container registry "`
	FlyteAwsSecret   string `json:"flyte-aws-secret" pflag:", Container registry "`
	Command          string `json:"command" pflag:", Container registry "`
}
