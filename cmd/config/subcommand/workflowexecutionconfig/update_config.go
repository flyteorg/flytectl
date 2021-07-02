package workflowexecutionconfig

//go:generate pflags AttrUpdateConfig --default-var DefaultUpdateConfig --bind-default-var

// AttrUpdateConfig Matchable resource attributes configuration passed from command line
type AttrUpdateConfig struct {
	AttrFile string `json:"attrFile" pflag:",attribute file name to be used for updating attribute for the resource type."`
}

var DefaultUpdateConfig = &AttrUpdateConfig{}
