package cmd

import (
	"fmt"
	"os"

	"github.com/paypal/gorealis/gen-go/apache/aurora"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(stopCmd)

	// Stop subcommands
	stopCmd.AddCommand(stopMaintCmd)
	stopMaintCmd.Flags().IntVar(&monitorInterval,"interval", 5, "Interval at which to poll scheduler.")
	stopMaintCmd.Flags().IntVar(&monitorTimeout,"timeout", 50, "Time after which the monitor will stop polling and throw an error.")

	// Stop update

	stopCmd.AddCommand(stopUpdateCmd)
	stopUpdateCmd.Flags().StringVarP(env, "environment", "e", "", "Aurora Environment")
	stopUpdateCmd.Flags().StringVarP(role, "role", "r", "", "Aurora Role")
	stopUpdateCmd.Flags().StringVarP(name, "name", "n", "", "Aurora Name")

}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a service or maintenance on a host (DRAIN).",
}

var stopMaintCmd = &cobra.Command{
	Use:   "drain [space separated host list]",
	Short: "Stop maintenance on a host (move to NONE).",
	Long:  `Transition a list of hosts currently in a maintenance status out of it.`,
	Run:   endMaintenance,
}

var stopUpdateCmd = &cobra.Command{
	Use:   "update [update ID]",
	Short: "Stop update",
	Long:  `To be written.`,
	Run:   stopUpdate,
}

func endMaintenance(cmd *cobra.Command, args []string) {
	fmt.Println("Setting hosts to NONE maintenance status.")
	fmt.Println(args)
	_, result, err := client.EndMaintenance(args...)
	if err != nil {
		fmt.Printf("error: %+v\n", err.Error())
		os.Exit(1)
	}

	// Monitor change to NONE mode
	hostResult, err := monitor.HostMaintenance(
		args,
		[]aurora.MaintenanceMode{aurora.MaintenanceMode_NONE},
		monitorInterval,
		monitorTimeout)
	if err != nil {
		for host, ok := range hostResult {
			if !ok {
				fmt.Printf("Host %s did not transtion into desired mode(s)\n", host)
			}
		}

		fmt.Printf("error: %+v\n", err.Error())
		return
	}

	fmt.Println(result.String())
}

func stopUpdate(cmd *cobra.Command, args []string) {

	if len(args) != 1 {
		fmt.Println("Only a single update ID must be provided.")
		os.Exit(1)
	}

	fmt.Printf("Stopping (aborting) update [%s/%s/%s] %s\n", *env, *role, *name, args[0])

	resp, err := client.AbortJobUpdate(aurora.JobUpdateKey{
		Job: &aurora.JobKey{Environment: *env, Role: *role, Name: *name},
		ID:  args[0],
	},
		"")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(resp.String())

}
