package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var env, role, name string

func init() {
	rootCmd.AddCommand(killCmd)
	killCmd.Flags().StringVarP(&env, "environment", "e", "", "Aurora Environment")
	killCmd.Flags().StringVarP(&role, "role", "r", "", "Aurora Role")
	killCmd.Flags().StringVarP(&name, "name", "n", "", "Aurora Name")
	killCmd.MarkFlagRequired("environment")
	killCmd.MarkFlagRequired("role")
	killCmd.MarkFlagRequired("name")
}

var killCmd = &cobra.Command{
	Use:   "kill",
	Short: "kill an Aurora Job",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s %s %s", env, role, name)
	},
}
