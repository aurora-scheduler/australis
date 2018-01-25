package cmd

import (
	"fmt"
	"os"

	"github.com/paypal/gorealis/gen-go/apache/aurora"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(fetchCmd)

	// Sub-commands
	fetchCmd.AddCommand(taskConfigCmd)
	taskConfigCmd.Flags().StringVarP(&env, "environment", "e", "", "Aurora Environment")
	taskConfigCmd.Flags().StringVarP(&role, "role", "r", "", "Aurora Role")
	taskConfigCmd.Flags().StringVarP(&name, "name", "n", "", "Aurora Name")
}

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch information from Aurora",
}

var taskConfigCmd = &cobra.Command{
	Use:   "config" ,
	Short: "Fetch a list of task configurations from Aurora.",
	Long: `To be written.`,
	Run:  fetchTasks,
}

func fetchTasks(cmd *cobra.Command, args []string) {
	fmt.Printf("Fetching job configuration for [%s/%s/%s] \n", env, role, name)

	//TODO: Add filtering down by status
	taskQuery := &aurora.TaskQuery{Environment: env, Role: role, JobName: name}

	tasks, err := client.GetTasksWithoutConfigs(taskQuery)
	if err != nil {
		fmt.Printf("error: %+v\n", err.Error())
		os.Exit(1)
	}

	for _, t := range tasks {
		fmt.Println(t)
	}
}
