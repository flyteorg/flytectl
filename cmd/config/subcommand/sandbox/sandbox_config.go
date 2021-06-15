package sandbox


//go:generate pflags Config --default-var DefaultConfig
var (
	DefaultConfig = &Config{
		Debug: false,
	}
)

// Config
type Config struct {
	Debug bool `json:"debug" pflag:", Enable debugging"`
	Source string `json:"source" pflag:", Path of your repository"`
}
