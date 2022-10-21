package demo

import (
	sandboxCmdConfig "github.com/flyteorg/flytectl/cmd/config/subcommand/sandbox"
	cmdcore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/spf13/cobra"
)

// Long descriptions are whitespace sensitive when generating docs using sphinx.
const (
	demoShort = `Helps with demo interactions like start, teardown, status, and exec.`
	demoLong  = `
Flyte Demo is a fully standalone minimal environment for running Flyte.
It provides a simplified way of running Flyte demo as a single Docker container locally.
	
To create a demo cluster, run:
::

 flytectl demo start 

To remove a demo cluster, run:
::

 flytectl demo teardown

To check the status of the demo container, run:
::

 flytectl demo status

To execute commands inside the demo container, use exec:
::

 flytectl demo exec -- pwd 	
`
)

// Long descriptions are whitespace sensitive when generating docs using sphinx.
const (
	initShort = `Download the Flyte sandbox image, create local state folder and place a default config file in it`
	initLong  = `
Flyte Demo	
When you run::
::

 flytectl demo init  

flytectl will ensure you have the latest run time image, create a local state directory for you if not present,
and place a default configuration file for the Flyte binary in it.

You may update the flyte binary configuration file after the demo cluster has been started, but this command is useful
in cases where you know you will want to modify the config before creating the cluster.
`
)

// CreateDemoCommand will return demo command
func CreateDemoCommand() *cobra.Command {
	demo := &cobra.Command{
		Use:   "demo",
		Short: demoShort,
		Long:  demoLong,
	}

	demoResourcesFuncs := map[string]cmdcore.CommandEntry{
		"init": {CmdFunc: initDemoCluster, Aliases: []string{}, ProjectDomainNotRequired: true,
			Short: initShort,
			Long:  initLong, PFlagProvider: sandboxCmdConfig.DefaultConfig, DisableFlyteClient: true},
		"start": {CmdFunc: startDemoCluster, Aliases: []string{}, ProjectDomainNotRequired: true,
			Short: startShort,
			Long:  startLong, PFlagProvider: sandboxCmdConfig.DefaultConfig, DisableFlyteClient: true},
		"teardown": {CmdFunc: teardownDemoCluster, Aliases: []string{}, ProjectDomainNotRequired: true,
			Short: teardownShort,
			Long:  teardownLong, DisableFlyteClient: true},
		"status": {CmdFunc: demoClusterStatus, Aliases: []string{}, ProjectDomainNotRequired: true,
			Short: statusShort,
			Long:  statusLong},
		"exec": {CmdFunc: demoClusterExec, Aliases: []string{}, ProjectDomainNotRequired: true,
			Short: execShort,
			Long:  execLong, DisableFlyteClient: true},
	}

	cmdcore.AddCommands(demo, demoResourcesFuncs)

	return demo
}
