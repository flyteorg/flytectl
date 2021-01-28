package register

import "github.com/lyft/flytestdlib/config"

//go:generate pflags Config

var (
	defaultConfig = &Config{
		version : "v1",
	}
	section = config.MustRegisterSection("register", defaultConfig)
)

type Config struct {
	version string `json:"version" pflag:",version of the entity to be registered with flyte."`
}

func GetConfig() *Config {
	return section.GetConfig().(*Config)
}

