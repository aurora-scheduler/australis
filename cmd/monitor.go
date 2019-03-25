package cmd

import (
	"strings"
	"time"

	"github.com/paypal/gorealis/v2/gen-go/apache/aurora"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(monitorCmd)

	monitorCmd.AddCommand(monitorHostCmd.cmd)

	monitorHostCmd.cmd.Run = monitorHost
	monitorHostCmd.cmd.Flags().DurationVar(&monitorHostCmd.monitorInterval, "interval", time.Second*5, "Interval at which to poll scheduler.")
	monitorHostCmd.cmd.Flags().DurationVar(&monitorHostCmd.monitorTimeout, "timeout", time.Minute*10, "Time after which the monitor will stop polling and throw an error.")
	monitorHostCmd.cmd.Flags().StringSliceVar(&monitorHostCmd.statusList, "statuses", []string{aurora.MaintenanceMode_DRAINED.String()}, "List of acceptable statuses for a host to be in. (case-insensitive) [NONE, SCHEDULED, DRAINED, DRAINING]")
}

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Watch for a specific state change",
}

var monitorHostCmd = MonitorCmdConfig{
	cmd: &cobra.Command{
		Use:   "hosts",
		Short: "Watch a host maintenance status until it enters one of the desired statuses.",
		Long: `Provide a list of hosts to monitor for desired statuses. Statuses may be passed using the --statuses
flag with a list of comma separated statuses. Statuses include [NONE, SCHEDULED, DRAINED, DRAINING]`,
	},
	statusList: make([]string, 0),
}

func monitorHost(cmd *cobra.Command, args []string) {
	maintenanceModes := make([]aurora.MaintenanceMode, 0)

	for _, status := range monitorHostCmd.statusList {
		mode, err := aurora.MaintenanceModeFromString(strings.ToUpper(status))
		if err != nil {
			log.Fatal(err)
		}

		maintenanceModes = append(maintenanceModes, mode)
	}

	log.Infof("Monitoring for %v at %v intervals", monitorHostCmd.monitorTimeout, monitorHostCmd.monitorInterval)
	hostResult, err := client.HostMaintenanceMonitor(args, maintenanceModes, monitorHostCmd.monitorInterval, monitorHostCmd.monitorTimeout)

	maintenanceMonitorPrint(hostResult, maintenanceModes)

	if err != nil {
		log.Fatal(err)
	}
}
