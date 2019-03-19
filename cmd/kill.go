package cmd

import (
	realis "github.com/paypal/gorealis/v2"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(killCmd)

	/* Sub-Commands */

	// Kill Job
	killCmd.AddCommand(killJobCmd)

	killJobCmd.Flags().StringVarP(env, "environment", "e", "", "Aurora Environment")
	killJobCmd.Flags().StringVarP(role, "role", "r", "", "Aurora Role")
	killJobCmd.Flags().StringVarP(name, "name", "n", "", "Aurora Name")
	killJobCmd.MarkFlagRequired("environment")
	killJobCmd.MarkFlagRequired("role")
	killJobCmd.MarkFlagRequired("name")

	// Kill every task in the Aurora cluster
	killCmd.AddCommand(killEntireClusterCmd)
}

var killCmd = &cobra.Command{
	Use:   "kill",
	Short: "Kill an Aurora Job",
}

var killJobCmd = &cobra.Command{
	Use:   "job",
	Short: "Kill an Aurora Job",
	Run:   killJob,
}

var killEntireClusterCmd = &cobra.Command{
	Use:   "entire-cluster",
	Short: "Kill every task in the cluster.",
	Long:  `To be written.`,
	Run:   killEntireCluster,
}

func killJob(cmd *cobra.Command, args []string) {
	log.Infof("Killing job [Env:%s Role:%s Name:%s]\n", *env, *role, *name)

	job := realis.NewJob().
		Environment(*env).
		Role(*role).
		Name(*name)
	err := client.KillJob(job.JobKey())
	if err != nil {
		log.Fatalln(err)
	}

	if ok, err := client.InstancesMonitor(job.JobKey(), 0, 5, 50); !ok || err != nil {
		log.Fatalln("Unable to kill all instances of job")
	}
}

func killEntireCluster(cmd *cobra.Command, args []string) {
	log.Println("This command will kill every single task inside of a cluster.")
	log.Println("Not implemented yet.")
}
