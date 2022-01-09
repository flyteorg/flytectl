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



//go:generate pflags CreateConfig --default-var DefaultCreateConfig --bind-default-var

// CreateConfig Config hold configuration for project create flags.
type CreateConfig struct {
	ID          string            `json:"id" pflag:",id for the project specified as argument."`
	Name        string            `json:"name" pflag:",name for the project specified as argument."`
	File        string            `json:"file" pflag:",file for the project definition."`
	Description string            `json:"description" pflag:",description for the project specified as argument."`
	Labels      map[string]string `json:"labels" pflag:",labels for the project specified as argument."`
	DryRun      bool              `json:"dryRun" pflag:",execute command without making any modifications."`
}

var (
	DefaultCreateConfig = &CreateConfig{
		Description: "",
		Labels:      map[string]string{},
	}
)

type ProjectDefinition struct {
	ID          string            `yaml:"id"`
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Labels      map[string]string `yaml:"labels"`
}

//go:generate pflags UpdateConfig --default-var DefaultUpdateConfig --bind-default-var

// UpdateConfig hold configuration for project update flags.
type UpdateConfig struct {
	ActivateProject bool `json:"activateProject" pflag:",Activates the project specified as argument."`
	ArchiveProject  bool `json:"archiveProject" pflag:",Archives the project specified as argument."`
	DryRun          bool `json:"dryRun" pflag:",execute command without making any modifications."`
	File        string            `json:"file" pflag:",file for the project definition."`
	Description string            `json:"description" pflag:",description for the project specified as argument."`
	Labels      map[string]string `json:"labels" pflag:",labels for the project specified as argument."`
}

var DefaultUpdateConfig = &UpdateConfig{}
