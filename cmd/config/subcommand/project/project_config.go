package project

import (
	"github.com/flyteorg/flytectl/pkg/filters"
)

//go:generate pflags Config --default-var DefaultConfig --bind-default-var
var (
	DefaultConfig = &Config{
		Filter: filters.DefaultFilter,
	}
)

// Config holds the flag for get project
type Config struct {
	Filter filters.Filters `json:"filter" pflag:","`
}

//go:generate pflags ConfigProject --default-var DefaultProjectConfig --bind-default-var

// ConfigProject hold configuration for project update flags.
type ConfigProject struct {
	ID              string            `json:"id" pflag:",id for the project specified as argument."`
	ActivateProject bool              `json:"activateProject" pflag:",(Deprecated) Activates the project specified as argument. Only used in update"`
	ArchiveProject  bool              `json:"archiveProject" pflag:",(Deprecated) Archives the project specified as argument. Only used in update"`
	Activate        bool              `json:"activate" pflag:",Activates the project specified as argument. Only used in update"`
	Archive         bool              `json:"archive" pflag:",Archives the project specified as argument. Only used in update"`
	Name            string            `json:"name" pflag:",name for the project specified as argument."`
	DryRun          bool              `json:"dryRun" pflag:",execute command without making any modifications."`
	Description     string            `json:"description" pflag:",description for the project specified as argument."`
	Labels          map[string]string `json:"labels" pflag:",labels for the project specified as argument."`
	File            string            `json:"file" pflag:",file for the project definition."`
}

var DefaultProjectConfig = &ConfigProject{
	Description: "",
	Labels:      map[string]string{},
}
