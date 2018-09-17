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
	startCmd.AddCommand(startSnapshotCmd)
    startCmd.AddCommand(startBackupCmd)

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
		[]aurora.MaintenanceMode{aurora.MaintenanceMode_DRAINED},
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

var startSnapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Force the leading scheduler to perform a Snapshot.",
	Long: `Takes a Snapshot of the in memory state of the Apache Aurora cluster and
writes it to the Mesos replicated log. This should NOT be confused with a backup.`,
	Run:  snapshot,
}

func snapshot(cmd *cobra.Command, args []string) {

	fmt.Println("Forcing scheduler to write snapshot to Mesos replicated log")
	err := client.Snapshot()
	if err != nil {
		fmt.Printf("error: %+v\n", err.Error())
		os.Exit(1)
	} else {
		fmt.Println("Snapshot started successfully")
	}
}

var startBackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Force the leading scheduler to perform a Backup.",
	Long: `Force the Aurora Scheduler to write a backup of the latest snapshot to the filesystem 
of the leading scheduler.`,
	Run:  backup,
}

func backup(cmd *cobra.Command, args []string) {

	fmt.Println("Forcing scheduler to write a Backup of latest Snapshot to file system")
	err := client.PerformBackup()
	if err != nil {
		fmt.Printf("error: %+v\n", err.Error())
		os.Exit(1)
	} else {
		fmt.Println("Backup started successfully")
	}
}
