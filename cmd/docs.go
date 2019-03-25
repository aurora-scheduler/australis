package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func init() {
	rootCmd.AddCommand(docsCmd)
	docsCmd.Hidden = true
}

var docsCmd = &cobra.Command{
	Use:               "docs",
	Short:             "Generate documents in markdown format for Australis.",
	PersistentPreRun:  func(cmd *cobra.Command, args []string) {}, // We don't need a realis client for this cmd
	PersistentPostRun: func(cmd *cobra.Command, args []string) {}, // We don't need a realis client for this cmd
	Run: func(cmd *cobra.Command, args []string) {
		err := doc.GenMarkdownTree(rootCmd, "./docs")
		if err != nil {
			log.Fatal(err)
		}
	},
}
