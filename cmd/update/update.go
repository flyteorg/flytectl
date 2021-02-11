package update

import (
	cmdcore "github.com/lyft/flytectl/cmd/core"

	"github.com/spf13/cobra"
)

// Long descriptions are whitespace sensitive when generating docs using sphinx.
const (
	updateUse   = "update"
	updateShort = `
Used for updating flyte resources eg: project.
`
	updatecmdLong  = `
Currently this command only provides subcommands to update project.
Takes input project which need to be archived or unarchived. Name of the project to be updated is mandatory field.
Example update project to activate it.
::

 bin/flytectl update project -p flytesnacks --activateProject
`
	projectShort = "Updates project resources"
	projectLong  = `
Updates the project according the flags passed.Allows you to archive or activate a project.
Activates project named flytesnacks.
::

 bin/flytectl update project -p flytesnacks --activateProject

Archives project named flytesnacks.

::

 bin/flytectl get project flytesnacks --archiveProject

Activates project named flytesnacks using short option -t.
::

 bin/flytectl update project -p flytesnacks -t

Archives project named flytesnacks using short option -a.

::

 bin/flytectl update project flytesnacks -a

Incorrect usage when passing both archive and activate.

::

 bin/flytectl update project flytesnacks -a -t

Incorrect usage when passing unknown-project.

::

 bin/flytectl update project unknown-project -a

Incorrect usage when passing valid project using -p option.

::

 bin/flytectl update project unknown-project -a -p known-project

Usage
`
)

// CreateUpdateCommand will return update command
func CreateUpdateCommand() *cobra.Command {
	updateCmd := &cobra.Command{
		Use:   updateUse,
		Short: updateShort,
		Long: updatecmdLong,
	}

	updateResourcesFuncs := map[string]cmdcore.CommandEntry{
		"project": {CmdFunc: updateProjectsFunc, Aliases: []string{"projects"}, ProjectDomainNotRequired: true, PFlagProvider: projectConfig,
			Short: projectShort,
			Long:  projectLong},
	}

	cmdcore.AddCommands(updateCmd, updateResourcesFuncs)
	return updateCmd
}
