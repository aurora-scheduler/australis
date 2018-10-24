package cmd

import (
    "fmt"
    "github.com/spf13/viper"
    "os"
    "strings"
    "time"

    "github.com/paypal/gorealis"
    "github.com/spf13/cobra"
)

var username, password, zkAddr, schedAddr string
var env, role, name = new(string), new(string), new(string)
var client realis.Realis
var monitor *realis.Monitor
var skipCertVerification bool
var caCertsPath string
var clientKey, clientCert string
var configFile string

const australisVer = "v0.0.5"

var monitorInterval, monitorTimeout int

func init() {


	rootCmd.SetVersionTemplate(`{{printf "%s\n" .Version}}`)

	rootCmd.PersistentFlags().StringVarP(&zkAddr, "zookeeper", "z", "", "Zookeeper node(s) where Aurora stores information. (comma separated list)")
	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "Username to use for API authentication")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "Password to use for API authentication")
	rootCmd.PersistentFlags().StringVarP(&schedAddr, "scheduler_addr", "s", "", "Aurora Scheduler's address.")
	rootCmd.PersistentFlags().StringVarP(&clientKey, "clientKey", "k", "", "Client key to use to connect to Aurora.")
	rootCmd.PersistentFlags().StringVarP(&clientCert, "clientCert", "c", "", "Client certificate to use to connect to Aurora.")
	rootCmd.PersistentFlags().StringVarP(&caCertsPath, "caCertsPath", "a", "", "Path where CA certificates can be found.")
	rootCmd.PersistentFlags().BoolVarP(&skipCertVerification, "skipCertVerification", "i", false, "Skip CA certificate hostname verification.")
	rootCmd.PersistentFlags().StringVar(&configFile, "config",  "/etc/aurora/australis.yml", "Config file to use.")

}

var rootCmd = &cobra.Command{
	Use:              "australis",
	Short:            "australis is a client for Apache Aurora",
	Long:             `A light-weight command line client for use with Apache Aurora built using gorealis.`,
	PersistentPreRun: connect,
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		// Make all children close the client by default upon terminating
		client.Close()
	},
	Version: australisVer,
}

func Execute() {
	rootCmd.Execute()
}

func connect(cmd *cobra.Command, args []string) {
    var err error

    zkAddrSlice := strings.Split(zkAddr, ",")

    viper.SetConfigFile(configFile)
    err = viper.ReadInConfig()
    if err == nil {
        // Best effort load configuration. Will only set config values when flags have not set them already.
        if viper.IsSet("zk") && len(zkAddrSlice) == 1 && zkAddrSlice[0] == "" {
            zkAddrSlice = viper.GetStringSlice("zk")
        }

        if viper.IsSet("username") && username == "" {
            username = viper.GetString("username")
        }

        if viper.IsSet("password") && password == "" {
            password = viper.GetString("password")
        }

        if viper.IsSet("clientKey") && clientKey == "" {
            clientKey = viper.GetString("clientKey")
        }

        if viper.IsSet("clientCert") && clientCert == "" {
            clientCert = viper.GetString("clientCert")
        }

        if viper.IsSet("caCertsPath") && caCertsPath == "" {
            caCertsPath = viper.GetString("caCertsPath")
        }

        if viper.IsSet("skipCertVerification") && !skipCertVerification {
            skipCertVerification = viper.GetBool("skipCertVerification")
        }
    }

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
	if len(zkAddrSlice) > 0 && zkAddrSlice[0] != "" {
		// Configure Zookeeper to connect
		zkOptions := []realis.ZKOpt{ realis.ZKEndpoints(zkAddrSlice...), realis.ZKPath("/aurora/scheduler")}
		realisOptions = append(realisOptions, realis.ZookeeperOptions(zkOptions...))
	} else if schedAddr != "" {
		realisOptions = append(realisOptions, realis.SchedulerUrl(schedAddr))
	} else {
		fmt.Println("Zookeeper address or Scheduler URL must be provided.")
		os.Exit(1)
	}

	// Client certificate configuration if available
	if clientKey != "" || clientCert != "" || caCertsPath != "" {
		realisOptions = append(realisOptions,
			realis.Certspath(caCertsPath),
			realis.ClientCerts(clientKey, clientCert),
			realis.InsecureSkipVerify(skipCertVerification))
	}

	// Connect to Aurora Scheduler and create a client object
	client, err = realis.NewRealisClient(realisOptions...)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	monitor = &realis.Monitor{Client: client}

}
