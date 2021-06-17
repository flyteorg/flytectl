package workflow

import (
	"github.com/flyteorg/flytectl/pkg/filters"
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