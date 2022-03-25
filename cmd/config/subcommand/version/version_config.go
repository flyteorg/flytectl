package execution

//go:generate pflags Config --default-var DefaultConfig --bind-default-var
var (
	DefaultConfig = &Config{}
)

// Config stores the flags required by version command
type Config struct {
	ControlPlane bool `json:"ctrlPlane" pflag:",gets control plane version."`
}
