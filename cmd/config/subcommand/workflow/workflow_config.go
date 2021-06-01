package workflow

import (
	"fmt"
	"github.com/flyteorg/flytectl/pkg/filters"
	"github.com/goccy/go-graphviz"
	"strings"
)

//go:generate pflags Config --default-var DefaultConfig

var (
	DefaultConfig = &Config{
		Filter: filters.DefaultFilter,
	}
)

// Config commandline configuration
type Config struct {
	Version string          `json:"version" pflag:",version of the workflow to be fetched."`
	Latest  bool            `json:"latest" pflag:", flag to indicate to fetch the latest version, version flag will be ignored in this case"`
	Filter  filters.Filters `json:"filter" pflag:","`
	Visualize  string `json:"visualize" pflag:",optional flag to visualize a workflow as one of [png, dot, svg, jpg]"`
	OutputFile string `json:"output_file" pflag:",path and a filename of where the output image should be dumped. This can only be used in concert with visualize"`
}

// GraphvizFormat returns the graphviz format based on the user-specified value or an error
func (cfg *Config) GraphvizFormat() (graphviz.Format, error) {
	switch strings.ToLower(cfg.Visualize) {
	case "png":
		return graphviz.PNG, nil
	case "svg":
		return graphviz.SVG, nil
	case "dot":
		return graphviz.XDOT, nil
	case "jpg":
		return graphviz.JPG, nil
	}
	return graphviz.SVG, fmt.Errorf("unsupported visualization format [%s]. Supported [png, dot, svg, jpg]", cfg.Visualize)
}