package cmd

import (
	"fmt"

	"github.com/paypal/gorealis"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
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
	resp, err := client.KillJob(job.JobKey())
	if err != nil {
		log.Fatalln(err)
	}

	if ok, err := monitor.Instances(job.JobKey(), 0, 5, 50); !ok || err != nil {
		log.Fatalln("Unable to kill all instances of job")
	}

	if toJson {
		fmt.Println(toJSON(resp.GetResult_()))
	} else {
		fmt.Println(resp.GetResult_())
	}
}

func killEntireCluster(cmd *cobra.Command, args []string) {
	log.Println("This command will kill every single task inside of a cluster.")
	log.Println("Not implemented yet.")
}
