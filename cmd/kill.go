package cmd

import (
	"github.com/spf13/cobra"
	"fmt")

func init() {
	rootCmd.AddCommand(killCmd)
  }

var killCmd = &cobra.Command{
	Use:   "kill",
	Short: "kill an Aurora Job",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Kill job")
	},
  }