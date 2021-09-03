package sandbox

//go:generate pflags Config --default-var DefaultConfig --bind-default-var
var (
	DefaultConfig = &Config{
		Image: "cr.flyte.org/flyteorg/flyte-sandbox",
		Version: "dind",
	}
)

//Config
type Config struct {
	Source  string `json:"source" pflag:",Path of your source code"`
	Image  string `json:"image" pflag:",flyte sandbox custom image"`
	Version string `json:"version" pflag:",Version of flyte"`
}
