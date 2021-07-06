package init

//go:generate pflags Config --default-var DefaultConfig
var (
	DefaultConfig = &Config{}
)

//Configs
type Config struct {
	Host string `json:"host" pflag:", Endpoint of flyte admin"`
}
