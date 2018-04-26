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

var username, password, zkAddr, schedAddr string
var env, role, name string
var client realis.Realis
var monitor *realis.Monitor
var insecureSkipVerify bool
var caCertsPath string
var clientKey, clientCert string

func init() {
	rootCmd.PersistentFlags().StringVarP(&zkAddr, "zookeeper", "z", "", "Zookeeper node(s) where Aurora stores information.")
	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "Username to use for API authentication")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "Password to use for API authentication")
	rootCmd.PersistentFlags().StringVarP(&schedAddr, "scheduler_addr", "s", "", "Aurora Scheduler's address.")
	rootCmd.PersistentFlags().StringVarP(&clientKey, "clientKey", "k", "", "Client key to use to connect to Aurora.")
	rootCmd.PersistentFlags().StringVarP(&clientCert, "clientCert", "c", "", "Client certificate to use to connect to Aurora.")
	rootCmd.PersistentFlags().StringVarP(&caCertsPath, "caCertsPath", "a", "", "CA certificates path to use.")
}

func Execute() {
	rootCmd.Execute()
}

func connect(cmd *cobra.Command, args []string) {

	var err error

	realisOptions := []realis.ClientOption{realis.BasicAuth(username, password),
		realis.ThriftJSON(),
		realis.TimeoutMS(20000),
		realis.BackOff(realis.Backoff{
			Steps:    2,
			Duration: 10 * time.Second,
			Factor:   2.0,
			Jitter:   0.1,
		})}

	// Prefer zookeeper if both ways of connecting are provided
	if zkAddr != "" {

		// Configure Zookeeper to connect
		zkOptions := []realis.ZKOpt{ realis.ZKEndpoints(zkAddr), realis.ZKPath("/aurora/scheduler")}

		if clientKey != "" || clientCert != "" || caCertsPath != "" {
			zkOptions = append(zkOptions, realis.ZKAuroraPortOverride(8081), realis.ZKAuroraSchemeOverride("https"))

			realisOptions = append(realisOptions, realis.Certspath(caCertsPath), realis.ClientCerts(clientKey, clientCert))
		}

		realisOptions = append(realisOptions, realis.ZookeeperOptions(zkOptions...))
	} else if schedAddr != "" {
		realisOptions = append(realisOptions, realis.SchedulerUrl(schedAddr))
	} else {
		fmt.Println("Zookeeper address or Scheduler URL must be provided.")
		os.Exit(1)
	}

	// Connect to Aurora Scheduler and create a client object
	client, err = realis.NewRealisClient(realisOptions...)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	monitor = &realis.Monitor{Client: client}

}
