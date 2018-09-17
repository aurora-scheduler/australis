package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/paypal/gorealis"
	"github.com/paypal/gorealis/gen-go/apache/aurora"
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
	log.Printf("Killing job [Env:%s Role:%s Name:%s]\n", *env, *role, *name)

	job := realis.NewJob().
		Environment(*env).
		Role(*role).
		Name(*name)
	resp, err := client.KillJob(job.JobKey())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if resp.ResponseCode == aurora.ResponseCode_OK {
		if ok, err := monitor.Instances(job.JobKey(), 0, 5, 50); !ok || err != nil {
			log.Println("Unable to kill all instances of job")
			os.Exit(1)
		}
	}
	fmt.Println(resp.String())
}

func killEntireCluster(cmd *cobra.Command, args []string) {
	log.Println("This command will kill every single task inside of a cluster.")
	log.Println("Not implemented yet.")
}
