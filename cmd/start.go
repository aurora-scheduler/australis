package cmd

import (
	"fmt"
	"os"

	"github.com/paypal/gorealis/gen-go/apache/aurora"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)

	// Sub-commands
	startCmd.AddCommand(startMaintCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a service or maintenance on a host (DRAIN).",
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
	fmt.Println("Setting hosts to DRAINING")
	fmt.Println(args)
	_, result, err := client.DrainHosts(args...)
	if err != nil {
		fmt.Printf("error: %+v\n", err.Error())
		os.Exit(1)
	}

	// Monitor change to DRAINING and DRAINED mode
	hostResult, err := monitor.HostMaintenance(
		args,
		[]aurora.MaintenanceMode{aurora.MaintenanceMode_DRAINED, aurora.MaintenanceMode_DRAINING},
		5,
		10)
	if err != nil {
		for host, ok := range hostResult {
			if !ok {
				fmt.Printf("Host %s did not transtion into desired mode(s)\n", host)
			}
		}

		fmt.Printf("error: %+v\n", err.Error())
		return
	}

	fmt.Print(result.String())
}
