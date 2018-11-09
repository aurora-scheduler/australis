package cmd

import (
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(env, "environment", "e", "", "Aurora Environment")
	createCmd.Flags().StringVarP(role, "role", "r", "", "Aurora Role")
	createCmd.Flags().StringVarP(name, "name", "n", "", "Aurora Name")
	createCmd.MarkFlagRequired("environment")
	createCmd.MarkFlagRequired("role")
	createCmd.MarkFlagRequired("name")
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an Aurora Job",
	Run:   createJob,
}

func createJob(cmd *cobra.Command, args []string) {
	log.Println("Not implemented yet.")
}
