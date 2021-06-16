package sandbox

//go:generate pflags Config --default-var DefaultConfig
var (
	DefaultConfig = &Config{
		Debug: true,
	}
)

// Config
type Config struct {
	Debug      bool   `json:"debug" pflag:", Enable debugging"`
	SnacksRepo string `json:"flytesnacks-path" pflag:", Path of your flytesnacks repository"`
}
