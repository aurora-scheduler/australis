package cmd

import (
	"fmt"
	"os"

	"github.com/paypal/gorealis"
	"github.com/paypal/gorealis/gen-go/apache/aurora"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(fetchCmd)

	// Sub-commands

	// Fetch Task Config
	fetchCmd.AddCommand(taskConfigCmd)
	taskConfigCmd.Flags().StringVarP(&env, "environment", "e", "", "Aurora Environment")
	taskConfigCmd.Flags().StringVarP(&role, "role", "r", "", "Aurora Role")
	taskConfigCmd.Flags().StringVarP(&name, "name", "n", "", "Aurora Name")

	// Fetch Leader
	fetchCmd.AddCommand(leaderCmd)
}

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch information from Aurora",
}

var taskConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Fetch a list of task configurations from Aurora.",
	Long:  `To be written.`,
	Run:   fetchTasks,
}

var leaderCmd = &cobra.Command{
	Use:               "leader",
	PersistentPreRun:  func(cmd *cobra.Command, args []string) {}, //We don't need a realis client for this cmd
	PersistentPostRun: func(cmd *cobra.Command, args []string) {}, //We don't need a realis client for this cmd
	Short:             "Fetch current Aurora leader given Zookeeper nodes. Pass Zookeeper nodes separated by a space as an argument to this command.",
	Long:              `To be written.`,
	Run:               fetchLeader,
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

func fetchLeader(cmd *cobra.Command, args []string) {
	fmt.Printf("Fetching leader from %v \n", args)

	if len(args) < 1 {
		fmt.Println("At least one Zookeper node address must be passed in.")
		os.Exit(1)
	}

	url, err := realis.LeaderFromZKOpts(realis.ZKEndpoints(args...), realis.ZKPath("/aurora/scheduler"))

	if err != nil {
		fmt.Printf("error: %+v\n", err.Error())
		os.Exit(1)
	}

	fmt.Print(url)
}
