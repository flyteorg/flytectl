package config

import (
	"fmt"
	"strings"

	"github.com/flyteorg/flytestdlib/config"

	"github.com/flyteorg/flytectl/pkg/printer"
)

//go:generate pflags Config

var (
	DefaultLimit  int32 = 100
	defaultConfig       = &Config{
		Limit: DefaultLimit,
	}
	section = config.MustRegisterSection("root", defaultConfig)
)

// Config hold configration for flytectl flag
type Config struct {
	Project       string `json:"project" pflag:",Specifies the project to work on."`
	Domain        string `json:"domain" pflag:",Specified the domain to work on."`
	Output        string `json:"output" pflag:",Specified the output type."`
	FieldSelector string `json:"field-selector" pflag:",Specifies the Field selector"`
	SortBy        string `json:"sort-by" pflag:",Specifies which field to sort results "`
	// TODO: Support paginated queries
	Limit         int32  `json:"limit" pflag:",Specifies the limit"`
	Asc           bool   `json:"asc"  pflag:",Specifies the sorting order"`
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
