package sandbox

//go:generate pflags SandboxConfig --default-var DefaultConfig
var (
	DefaultConfig = &SandboxConfig{}
)

// SandboxConfig
type SandboxConfig struct {
	Source string `json:"source" pflag:", Path of your source code"`
}
