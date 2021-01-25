package register

import "github.com/lyft/flytestdlib/config"

//go:generate pflags Config

var (
	defaultConfig = &Config{
		Files:  [] string{},
	}
	section = config.MustRegisterSection("register", defaultConfig)
)

type Config struct {
	Files []string `json:"files" pflag:",List of serialized files used for registering to flyte."`
}

func GetConfig() *Config {
	return section.GetConfig().(*Config)
}

