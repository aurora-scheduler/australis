/**
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"strings"
	"time"

	"github.com/aurora-scheduler/australis/internal"
	"github.com/spf13/viper"

	realis "github.com/aurora-scheduler/gorealis/v2"
	"github.com/spf13/cobra"

	"github.com/sirupsen/logrus"
)

var username, password, zkAddr, schedAddr string
var env, role, name = new(string), new(string), new(string)
var ram, disk int64
var cpu float64
var client *realis.Client
var skipCertVerification bool
var caCertsPath string
var clientKey, clientCert string
var configFile string
var toJson bool
var fromJson bool
var fromJsonFile string
var logLevel string
var duration time.Duration
var percent float64
var count int64
var filename string
var message = new(string)
var updateID string
var log = logrus.New()

const australisVer = "v0.22.0"

var forceDrainTimeout time.Duration

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
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "/etc/aurora/australis.yml", "Config file to use.")
	rootCmd.PersistentFlags().BoolVar(&toJson, "toJSON", false, "Print output in JSON format.")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "logLevel", "l", "info", "Set logging level ["+internal.GetLoggingLevels()+"].")
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
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// TODO(rdelvalle): Move more from connect into this function
func setConfig(cmd *cobra.Command, args []string) {
	lvl, err := logrus.ParseLevel(logLevel)

	if err != nil {
		log.Fatalf("Log level %v is not valid", logLevel)
	}

	log.SetLevel(lvl)
	internal.Logger(log)
}

func connect(cmd *cobra.Command, args []string) {
	var err error

	setConfig(cmd, args)

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
		realis.Timeout(20 * time.Second),
		realis.BackOff(realis.Backoff{
			Steps:    2,
			Duration: 10 * time.Second,
			Factor:   2.0,
			Jitter:   0.1,
		}),
		realis.SetLogger(log)}

	// Prefer zookeeper if both ways of connecting are provided
	if len(zkAddrSlice) > 0 && zkAddrSlice[0] != "" {
		// Configure Zookeeper to connect
		zkOptions := []realis.ZKOpt{realis.ZKEndpoints(zkAddrSlice...), realis.ZKPath("/aurora/scheduler")}
		realisOptions = append(realisOptions, realis.ZookeeperOptions(zkOptions...))
	} else if schedAddr != "" {
		realisOptions = append(realisOptions, realis.SchedulerUrl(schedAddr))
	} else {
		log.Fatalln("Zookeeper address or Scheduler URL must be provided.")
	}

	// Client certificate configuration if available
	if clientKey != "" || clientCert != "" || caCertsPath != "" {
		realisOptions = append(realisOptions,
			realis.CertsPath(caCertsPath),
			realis.ClientCerts(clientKey, clientCert),
			realis.InsecureSkipVerify(skipCertVerification))
	}

	// Connect to Aurora Scheduler and create a client object
	client, err = realis.NewClient(realisOptions...)

	if err != nil {
		log.Fatal(err)
	}
}
