package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "ionbeam",
	Short: "ionbeam is a client for Apache Aurora",
	Long:  `A light-weight, intuitive command line client for use with Apache Aurora.`,
}

func Execute() {
	rootCmd.Execute()
}
