package execution

//go:generate pflags ExecutionDeleteConfig --default-var DefaultExecutionDeleteConfig --bind-default-var
var DefaultExecutionDeleteConfig = &ExecutionDeleteConfig{}

// ExecutionDeleteConfig stores the flags required by delete execution
type ExecutionDeleteConfig struct {
	DryRun bool `json:"dryRun" pflag:",execute local operations without making any modifications (skip or mock all server communication)"`
}
