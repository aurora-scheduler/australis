package cmd

import (
	"fmt"
	"log"
	"os"

	realis "github.com/paypal/gorealis"
	"github.com/paypal/gorealis/gen-go/apache/aurora"
	"github.com/spf13/cobra"
)

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
	Short: "Kill an Aurora Job",
	Run:   killJob,
}

func killJob(cmd *cobra.Command, args []string) {
	log.Printf("Killing job [Env:%s Role:%s Name:%s]\n", env, role, name)

	job := realis.NewJob().
		Environment(env).
		Role(role).
		Name(name)
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
