/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "This action generates completion scripts.",
	Long: `To load completions in Bash, the general format is:

Bash:

  $ source <(flytectl completion bash)

  # To load completions for each session in Linux, execute the below line:
  # Linux:
  $ flytectl completion bash > /etc/bash_completion.d/flytectl
  # To load completions for each session in macOS, execute the below line:
  # macOS:
  $ flytectl completion bash > /usr/local/etc/bash_completion.d/flytectl

Zsh:

  # If shell completion is not already enabled in your environment,
  # it needs to be enabled. This can be done by executing the below command:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions in Zsh, execute the below line:
  $ flytectl completion zsh > "${fpath[1]}/_flytectl"

  # Note: A new shell has to be started for the setup to take effect.

fish:
	# To load completions in fish, the general format is:
  $ flytectl completion fish | source

  # To load completions for each session in fish, execute the below line:
  $ flytectl completion fish > ~/.config/fish/completions/flytectl.fish

PowerShell:
	# To load completions on PowerShell, the general format is:
  PS> flytectl completion powershell | Out-String | Invoke-Expression

  # To load completions for a new session, execute the below line, and source this file from your Powershell profile:
  PS> flytectl completion powershell > flytectl.ps1

`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			return cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			return cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
		return nil
	},
}
