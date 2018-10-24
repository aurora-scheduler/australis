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

	/* Fetch Leader */

	leaderCmd.Flags().String("zkPath", "/aurora/scheduler", "Zookeeper node path where leader election happens")

	// Override usage template to hide global flags
	leaderCmd.SetUsageTemplate(`Usage:{{if .Runnable}}
{{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
{{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
	{{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
	{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

	Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
	{{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
	{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`)
	fetchCmd.AddCommand(leaderCmd)

	// Fetch jobs
	fetchJobsCmd.Flags().StringVarP(role, "role", "r", "", "Aurora Role")
	fetchCmd.AddCommand(fetchJobsCmd)

	// Fetch Status
	fetchCmd.AddCommand(fetchStatusCmd)
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
	Use:               "leader [zkNode0, zkNode1, ...zkNodeN]",
	PersistentPreRun:  func(cmd *cobra.Command, args []string) {}, //We don't need a realis client for this cmd
	PersistentPostRun: func(cmd *cobra.Command, args []string) {}, //We don't need a realis client for this cmd
	Short:             "Fetch current Aurora leader given Zookeeper nodes. ",
	Long:              `Gets the current leading aurora scheduler instance using information from Zookeeper path.
Pass Zookeeper nodes separated by a space as an argument to this command.`,
	Run:               fetchLeader,
}

var fetchJobsCmd = &cobra.Command{
	Use:   "jobs",
	Short: "Fetch a list of task Aurora running under a role.",
	Long:  `To be written.`,
	Run:   fetchJobs,
}

var fetchStatusCmd = &cobra.Command{
	Use: "status",
	Short: "Fetch the maintenance status of a node from Aurora",
	Long: `This command will print the actual status of the mesos agent nodes in Aurora server`,
	Run: fetchStatus,
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

func fetchStatus(cmd *cobra.Command, args []string) {
	fmt.Printf("Fetching maintenance status for %v \n", args)
	_, result, err := client.MaintenanceStatus(args...)
	if err != nil {
		fmt.Printf("error: %+v\n", err.Error())
		os.Exit(1)
	}

	for k := range result.Statuses {
		fmt.Printf("Result: %s:%s\n", k.Host, k.Mode)
	}
}

func fetchLeader(cmd *cobra.Command, args []string) {
	fmt.Printf("Fetching leader from %v \n", args)

	if len(args) < 1 {
		fmt.Println("At least one Zookeper node address must be passed in.")
		os.Exit(1)
	}

	url, err := realis.LeaderFromZKOpts(realis.ZKEndpoints(args...), realis.ZKPath(cmd.Flag("zkPath").Value.String()))

	if err != nil {
		fmt.Printf("error: %+v\n", err.Error())
		os.Exit(1)
	}

	fmt.Println(url)
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
