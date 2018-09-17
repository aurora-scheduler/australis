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
	taskConfigCmd.Flags().StringVarP(env, "environment", "e", "", "Aurora Environment")
	taskConfigCmd.Flags().StringVarP(role, "role", "r", "", "Aurora Role")
	taskConfigCmd.Flags().StringVarP(name, "name", "n", "", "Aurora Name")

	// Fetch Leader
	fetchCmd.AddCommand(leaderCmd)

	// Fetch jobs
	fetchCmd.AddCommand(fetchJobsCmd)
	fetchJobsCmd.Flags().StringVarP(role, "role", "r", "", "Aurora Role")
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

var fetchJobsCmd = &cobra.Command{
	Use:   "jobs",
	Short: "Fetch a list of task Aurora running under a role.",
	Long:  `To be written.`,
	Run:   fetchJobs,
}

func fetchTasks(cmd *cobra.Command, args []string) {
	fmt.Printf("Fetching job configuration for [%s/%s/%s] \n", env, role, name)

	// Task Query takes nil for values it shouldn't need to match against.
	// This allows us to potentially more expensive calls for specific environments, roles, or job names.
	if *env == "" {
		env = nil
	}
	if *role == "" {
		role = nil
	}
	if *role == "" {
		role = nil
	}
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

// TODO: Expand this to be able to filter by job name and environment.
func fetchJobs(cmd *cobra.Command, args []string)  {
	fmt.Printf("Fetching tasks under role: %s \n", *role)

	if *role == "" {
		fmt.Println("Role must be specified.")
		os.Exit(1)
	}

	if *role == "*" {
		fmt.Println("Warning: This is an expensive operation.")
		*role = ""
	}

	_, result, err := client.GetJobs(*role)

	if err != nil {
		fmt.Printf("error: %+v\n", err.Error())
		os.Exit(1)
	}

	for jobConfig, _ := range result.GetConfigs() {
		fmt.Println(jobConfig)
	}

}
