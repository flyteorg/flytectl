package workflowexecutionconfig

//go:generate pflags AttrDeleteConfig --default-var DefaultDelConfig --bind-default-var

// AttrDeleteConfig Matchable resource attributes configuration passed from command line
type AttrDeleteConfig struct {
	AttrFile string `json:"attrFile" pflag:",attribute file name to be used for delete attribute for the resource type."`
	DryRun   bool   `json:"dryRun" pflag:",execute local operations without making any modifications (skip or mock all server communication)"`
}

var DefaultDelConfig = &AttrDeleteConfig{}
