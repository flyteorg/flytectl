package sandbox

//go:generate pflags SandboxConfig --default-var DefaultConfig
var (
	DefaultConfig = &SandboxConfig{
		Version: "latest",
	}
)

// SandboxConfig
type SandboxConfig struct {
	SourcesPath string `json:"sourcesPath" pflag:",Path to your source code path where flyte workflows and tasks are."`
	Name        string `json:"name" pflag:", Name of your cluster"`
	Version     string `json:"version" pflag:", Flyte version"`
}
