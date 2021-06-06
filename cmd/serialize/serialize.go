package serialize

import (
	sconfig "github.com/flyteorg/flytectl/cmd/config/subcommand/serialize"
	cmdcore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/spf13/cobra"
)

// Long descriptions are whitespace sensitive when generating docs using sphinx.
const (
	serializeCmdShort = "Serialize flyte workflow"
	serializecmdLong  = `Example serialize workflow.
::

 flytectl serialize workflow
`
)

// RemoteSerializeCommand will return serialize command
func RemoteSerializeCommand() *cobra.Command {
	serializeCmd := &cobra.Command{
		Use:   "serialize",
		Short: serializeCmdShort,
		Long:  serializecmdLong,
	}
	serializeResourcesFuncs := map[string]cmdcore.CommandEntry{
		"workflow": {CmdFunc: serializeWorkflowFunc, Aliases: []string{"workflow"}, PFlagProvider: sconfig.DefaultFilesConfig,
			Short: serializeWorkflowShort, Long: serializeWorkflowLong},
	}
	cmdcore.AddCommands(serializeCmd, serializeResourcesFuncs)
	return serializeCmd
}
