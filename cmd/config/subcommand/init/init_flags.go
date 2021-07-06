package init

//go:generate pflags Config --default-var DefaultConfig --bind-to-default
var (
	DefaultConfig = &Config{
		Insecure: false,
	}
)

//Configs
type Config struct {
	Host     string `json:"host" pflag:",Endpoint of flyte admin"`
	Insecure bool   `json:"insecure" pflag:",Enable insecure mode"`
}
