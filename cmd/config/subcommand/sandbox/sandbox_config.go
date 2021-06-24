package sandbox

//go:generate pflags SandboxConfig --default-var DefaultConfig
var (
	DefaultConfig = &SandboxConfig{
		Version: "latest",
	}
)

// SandboxConfig
type SandboxConfig struct {
	Source string `json:"source" pflag:", Path of your source code"`
	Name string `json:"name" pflag:", Name of your cluster"`
	Version string `json:"version" pflag:", Flyte version"`
}
