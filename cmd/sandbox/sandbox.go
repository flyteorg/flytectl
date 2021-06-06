package sandbox

import (
	cmdcore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/spf13/cobra"
)

// Long descriptions are whitespace sensitive when generating docs using sphinx.
const (
	sandboxShort = `Used for playing with sandbox.`
	sandboxLong  = `
Example Create sandbox cluster.
::

 flytectl sandbox cluster 
`
)

// CreateSandboxCommand will return sandbox command
func CreateSandboxCommand() *cobra.Command {
	sandbox := &cobra.Command{
		Use:   "sandbox",
		Short: sandboxShort,
		Long:  sandboxLong,
	}

	sandboxResourcesFuncs := map[string]cmdcore.CommandEntry{
		"start": {CmdFunc: startSandboxCluster, Aliases: []string{"create"}, ProjectDomainNotRequired: true,
			Short: startShort,
			Long:  sandboxLong},
		"register": {CmdFunc: registerSandboxCluster, Aliases: []string{}, ProjectDomainNotRequired: false,
			Short: registerShort,
			Long:  registerLong},
		"teardown": {CmdFunc: teardownSandboxCluster, Aliases: []string{}, ProjectDomainNotRequired: true,
			Short: teardownShort,
			Long:  teardownLong},
	}

	cmdcore.AddCommands(sandbox, sandboxResourcesFuncs)

	return sandbox
}
