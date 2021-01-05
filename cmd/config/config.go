package config

import (
	"fmt"
	"strings"

	"github.com/lyft/flytestdlib/config"

	"github.com/lyft/flytectl/pkg/printer"
)

//go:generate pflags Config

var (
	defaultConfig = &Config{}
	createConfig = &CreateConfig{}
	section       = config.MustRegisterSection("root", defaultConfig)
	createSection       = config.MustRegisterSection("create", createConfig)
)

// Config hold configration for flytectl flag
type Config struct {
	Project string `json:"project" pflag:",Specifies the project to work on."`
	Domain  string `json:"domain" pflag:",Specified the domain to work on."`
	Output  string `json:"output" pflag:",Specified the output type."`
}

// CreateConfig hold configration for flytectl flag
type CreateConfig struct {
	Filename  string `json:"filename" pflag:",Specified the filename."`
	Name  string `json:"name" pflag:",Specified the name."`
	ID  string `json:"name" pflag:",Specified the id."`
	Description  string `json:"name" pflag:",Specified the description."`
	Labels  map[string]string `json:"labels" pflag:",Specified the labels."`

}

// OutputFormat will return output formate
func (cfg Config) OutputFormat() (printer.OutputFormat, error) {
	return printer.OutputFormatString(strings.ToUpper(cfg.Output))
}

// MustOutputFormat will validate the supported output formate and return output formate
func (cfg Config) MustOutputFormat() printer.OutputFormat {
	f, err := cfg.OutputFormat()
	if err != nil {
		panic(fmt.Sprintf("unsupported output format [%s], supported types %s", cfg.Output, printer.OutputFormats()))
	}
	return f
}

// GetConfig will return the config
func GetConfig() *Config {
	return section.GetConfig().(*Config)
}

// GetCreateConfig will return the config
func GetCreateConfig() *CreateConfig {
	return createSection.GetConfig().(*CreateConfig)
}
