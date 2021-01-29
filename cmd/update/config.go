package update

import (
	"github.com/lyft/flytestdlib/config"
)

//go:generate pflags Config

var (
	defaultConfig = &Config{
	}
	section       = config.MustRegisterSection("update", defaultConfig)
)

// Config hold configuration for project update flags.
type Config struct {
	ActivateProject bool `json:"activateProject" pflag:",Activates the project specified as argument."`
	ArchiveProject  bool `json:"archiveProject" pflag:",Archives the project specified as argument."`
}

// GetConfig will return the config
func GetConfig() *Config {
	return section.GetConfig().(*Config)
}
