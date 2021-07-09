package config

//go:generate pflags Config --default-var DefaultConfig --bind-default-var
var (
	DefaultConfig = &Config{
		Insecure: true,
	}
)

//Configs
type Config struct {
	Host        string `json:"host" pflag:",Endpoint of flyte admin"`
	Insecure    bool   `json:"insecure" pflag:",Enable insecure mode"`
	StorageType string `json:"storage" pflag:",Storage provider name S3/GCS"`
}
