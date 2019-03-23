package cmd

import (
	"strings"
	"time"

	"github.com/paypal/gorealis/v2/gen-go/apache/aurora"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(monitorCmd)

	monitorCmd.AddCommand(monitorHostCmd)
	monitorHostCmd.Flags().DurationVar(&monitorInterval, "interval", time.Second*5, "Interval at which to poll scheduler.")
	monitorHostCmd.Flags().DurationVar(&monitorTimeout, "timeout", time.Minute*10, "Time after which the monitor will stop polling and throw an error.")
	monitorHostCmd.Flags().StringSliceVar(&statusList, "statuses", []string{aurora.MaintenanceMode_DRAINED.String()}, "List of acceptable statuses for a host to be in. (case-insensitive) [NONE, SCHEDULED, DRAINED, DRAINING]")
}

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Watch for a specific state change",
}

var monitorHostCmd = &cobra.Command{
	Use:   "hosts",
	Short: "Watch a host maintenance status until it enters one of the desired statuses.",
	Long: `Provide a list of hosts to monitor for desired statuses. Statuses may be passed using the --statuses
flag with a list of comma separated statuses. Statuses include [NONE, SCHEDULED, DRAINED, DRAINING]`,
	PreRun: func(cmd *cobra.Command, args []string) {
		// Manually initializing default values for this command as the default value for shared variables will
		// be dependent on the order in which all commands were initialized
		monitorTimeout = time.Minute * 10
		monitorInterval = time.Second * 5
	},
	Run: monitorHost,
}

func monitorHost(cmd *cobra.Command, args []string) {
	maintenanceModes := make([]aurora.MaintenanceMode, 0)

	for _, status := range statusList {
		mode, err := aurora.MaintenanceModeFromString(strings.ToUpper(status))
		if err != nil {
			log.Fatal(err)
		}

		maintenanceModes = append(maintenanceModes, mode)
	}

	log.Println(monitorTimeout)
	log.Println(monitorInterval)
	hostResult, err := client.HostMaintenanceMonitor(args, maintenanceModes, monitorInterval, monitorTimeout)

	maintenanceMonitorPrint(hostResult, maintenanceModes)
	if err != nil {
		log.Fatalf("error: %+v", err)
	}
}
