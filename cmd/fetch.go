package cmd

import (
	"fmt"
	"github.com/spf13/pflag"

	"github.com/paypal/gorealis/v2"
	"github.com/paypal/gorealis/v2/gen-go/apache/aurora"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
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

	fetchCmd.AddCommand(leaderCmd)

	// Hijack help function to hide unnecessary global flags
	help := leaderCmd.HelpFunc()
	leaderCmd.SetHelpFunc(func(cmd *cobra.Command, s []string){
		if cmd.HasInheritedFlags(){
			cmd.InheritedFlags().VisitAll(func(f *pflag.Flag){
				if f.Name != "logLevel" {
					f.Hidden = true
				}
			})
		}
		help(cmd, s)
	})

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
	PreRun:            setConfig,
	Args: cobra.MinimumNArgs(1),
	Short:             "Fetch current Aurora leader given Zookeeper nodes. ",
	Long: `Gets the current leading aurora scheduler instance using information from Zookeeper path.
Pass Zookeeper nodes separated by a space as an argument to this command.`,
	Run: fetchLeader,
}

var fetchJobsCmd = &cobra.Command{
	Use:   "jobs",
	Short: "Fetch a list of task Aurora running under a role.",
	Long:  `To be written.`,
	Run:   fetchJobs,
}

var fetchStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Fetch the maintenance status of a node from Aurora",
	Long:  `This command will print the actual status of the mesos agent nodes in Aurora server`,
	Run:   fetchStatus,
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
		log.Fatalf("error: %+v\n", err)
	}

	if toJson {
		fmt.Println(toJSON(tasks))
	} else {
		for _, t := range tasks {
			fmt.Println(t)
		}
	}
}

func fetchStatus(cmd *cobra.Command, args []string) {
	log.Infof("Fetching maintenance status for %v \n", args)
	result, err := client.MaintenanceStatus(args...)
	if err != nil {
		log.Fatalf("error: %+v\n", err)
	}

	if toJson {
		fmt.Println(toJSON(result.Statuses))
	} else {
		for _, k := range result.GetStatuses() {
			fmt.Printf("Result: %s:%s\n", k.Host, k.Mode)
		}
	}
}

func fetchLeader(cmd *cobra.Command, args []string) {
	log.Infof("Fetching leader from %v \n", args)

	if len(args) < 1 {
		log.Fatalln("At least one Zookeeper node address must be passed in.")
	}

	url, err := realis.LeaderFromZKOpts(realis.ZKEndpoints(args...), realis.ZKPath(cmd.Flag("zkPath").Value.String()))

	if err != nil {
		log.Fatalf("error: %+v\n", err)
	}

	fmt.Println(url)
}

// TODO: Expand this to be able to filter by job name and environment.
func fetchJobs(cmd *cobra.Command, args []string) {
	log.Infof("Fetching tasks under role: %s \n", *role)

	if *role == "" {
		log.Fatalln("Role must be specified.")
	}

	if *role == "*" {
		log.Warnln("This is an expensive operation.")
		*role = ""
	}

	result, err := client.GetJobs(*role)

	if err != nil {
		log.Fatalf("error: %+v\n", err)
	}

	if toJson {
		var configSlice []*aurora.JobConfiguration

		for _, config := range result.GetConfigs() {
			configSlice = append(configSlice, config)
		}

		fmt.Println(toJSON(configSlice))
	} else {
		for jobConfig := range result.GetConfigs() {
			fmt.Println(jobConfig)
		}
	}
}
