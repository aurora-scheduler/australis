package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/paypal/gorealis"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:              "australis",
	Short:            "australis is a client for Apache Aurora",
	Long:             `A light-weight command line client for use with Apache Aurora built using gorealis.`,
	PersistentPreRun: connect,
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		// Make all children close the client by default upon terminating
		client.Close()
	},
}

var username, password, zkurl string
var env, role, name string
var client realis.Realis
var monitor *realis.Monitor

func init() {
	rootCmd.PersistentFlags().StringVarP(&zkurl, "zookeeper", "z", "", "Zookeeper node(s) where Aurora stores information.")
	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "Username to use for API authentification")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "Password to use for API authentification")

	// TODO(rdelvalle): Add plain URL support. Enforce that url or zookeepr is set. Consider adding support for clusters.json.
	// For now make ZK mandatory
	rootCmd.MarkPersistentFlagRequired("zookeeper")
}

func Execute() {
	rootCmd.Execute()
}

func connect(cmd *cobra.Command, args []string) {

	var err error

	// Connect to Aurora Scheduler and create a client object
	client, err = realis.NewRealisClient(
		realis.BasicAuth(username, password),
		realis.ThriftJSON(),
		realis.TimeoutMS(20000),
		realis.BackOff(&realis.Backoff{
			Steps:    2,
			Duration: 10 * time.Second,
			Factor:   2.0,
			Jitter:   0.1,
		}),
		realis.ZKUrl(zkurl))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	monitor = &realis.Monitor{Client: client}

}
