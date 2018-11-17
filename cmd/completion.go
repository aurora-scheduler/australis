package cmd

import (
	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Create auto completion for bash.",
	Long: `Create auto completion bash file for australis. Auto completion file must be placed in the correct
directory in order for bash to pick up the definitions.

Copy australis.completion.sh into the correct folder and rename to australis

In Linux, this directory is usually /etc/bash_completion.d/
In MacOS this directory is $(brew --prefix)/etc/bash_completion.d if auto completion was install through brew.
`,
	PersistentPreRun:  func(cmd *cobra.Command, args []string) {}, //We don't need a realis client for this cmd
	PersistentPostRun: func(cmd *cobra.Command, args []string) {}, //We don't need a realis client for this cmd
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenBashCompletionFile("australis.completion.sh")
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
