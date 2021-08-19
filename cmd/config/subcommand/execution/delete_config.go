package execution

//go:generate pflags ExecDeleteConfig --default-var DefaultExecDeleteConfig --bind-default-var

var DefaultExecDeleteConfig = &ExecDeleteConfig{}

// ExecutionDeleteConfig stores the flags required by delete execution
type ExecDeleteConfig struct {
	DryRun bool `json:"dryRun" pflag:",execute local operations without making any modifications (skip or mock all server communication)"`
}
