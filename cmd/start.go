package cmd

import (
	"fmt"
	"github.com/paypal/gorealis/gen-go/apache/aurora"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

func init() {
	rootCmd.AddCommand(startCmd)


	// Sub-commands
	startCmd.AddCommand(startMaintCmd)

	// Maintenance specific flags
	startMaintCmd.Flags().IntVar(&monitorInterval,"interval", 5, "Interval at which to poll scheduler.")
	startMaintCmd.Flags().IntVar(&monitorTimeout,"timeout", 50, "Time after which the monitor will stop polling and throw an error.")

}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a service, maintenance on a host (DRAIN), a snapshot, or a backup.",
}

var startMaintCmd = &cobra.Command{
	Use:   "drain [space separated host list]",
	Short: "Place a list of space separated Mesos Agents into maintenance mode.",
	Long: `Adds a Mesos Agent to Aurora's Drain list. Agents in this list
are not allowed to schedule new tasks and any tasks already running on this Agent
are killed and rescheduled in an Agent that is not in maintenance mode. Command
expects a space separated list of hosts to place into maintenance mode.`,
	Args: cobra.MinimumNArgs(1),
	Run:  drain,
}

func drain(cmd *cobra.Command, args []string) {
	log.Infoln("Setting hosts to DRAINING")
	log.Infoln(args)
	_, result, err := client.DrainHosts(args...)
	if err != nil {
		log.Fatalf("error: %+v\n", err)
	}

	log.Debugln(result)

	// Monitor change to DRAINING and DRAINED mode
	hostResult, err := monitor.HostMaintenance(
		args,
		[]aurora.MaintenanceMode{aurora.MaintenanceMode_DRAINED},
		monitorInterval,
		monitorTimeout)

	transitioned := make([]string, 0,0)
	if err != nil {
		nonTransitioned := make([]string, 0,0)

		for host, ok := range hostResult {
			if ok {
				transitioned = append(transitioned, host)
			} else {
				nonTransitioned = append(nonTransitioned, host)
			}
		}

		log.Printf("error: %+v\n", err)
		if toJson {
			fmt.Println(toJSON(nonTransitioned))
		} else {
			fmt.Println("Did not enter DRAINED status: ", nonTransitioned)
		}
	}

    if toJson {
        fmt.Println(toJSON(transitioned))
    } else {
		fmt.Println("Entered DRAINED status: ", transitioned)
    }
}

