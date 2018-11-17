package cmd

import (
	"github.com/paypal/gorealis/gen-go/apache/aurora"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

func init() {
	rootCmd.AddCommand(startCmd)


	// Sub-commands
	startCmd.AddCommand(startDrainCmd)

	// Maintenance specific flags
	startDrainCmd.Flags().IntVar(&monitorInterval,"interval", 5, "Interval at which to poll scheduler.")
	startDrainCmd.Flags().IntVar(&monitorTimeout,"timeout", 50, "Time after which the monitor will stop polling and throw an error.")

	startCmd.AddCommand(startSLADrainCmd)

	// SLA Maintenance specific flags
	startSLADrainCmd.Flags().IntVar(&monitorInterval,"interval", 5, "Interval at which to poll scheduler.")
	startSLADrainCmd.Flags().IntVar(&monitorTimeout,"timeout", 50, "Time after which the monitor will stop polling and throw an error.")

	startCmd.AddCommand(startMaintenanceCmd)

	// SLA Maintenance specific flags
	startMaintenanceCmd.Flags().IntVar(&monitorInterval,"interval", 5, "Interval at which to poll scheduler.")
	startMaintenanceCmd.Flags().IntVar(&monitorTimeout,"timeout", 50, "Time after which the monitor will stop polling and throw an error.")
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a service, maintenance on a host (DRAIN), a snapshot, or a backup.",
}

var startDrainCmd = &cobra.Command{
	Use:   "drain [space separated host list]",
	Short: "Place a list of space separated Mesos Agents into draining mode.",
	Long: `Adds a Mesos Agent to Aurora's Drain list. Agents in this list
are not allowed to schedule new tasks and any tasks already running on this Agent
are killed and rescheduled in an Agent that is not in maintenance mode. Command
expects a space separated list of hosts to place into maintenance mode.`,
	Args: cobra.MinimumNArgs(1),
	Run:  drain,
}

var startSLADrainCmd = &cobra.Command{
	Use:   "sla-drain [space separated host list]",
	Short: "Place a list of space separated Mesos Agents into maintenance mode using SLA awareness.",
	Long: `Adds a Mesos Agent to Aurora's Drain list. Agents in this list
are not allowed to schedule new tasks and any tasks already running on this Agent
are killed and rescheduled in an Agent that is not in maintenance mode. Command
expects a space separated list of hosts to place into maintenance mode.`,
	Args: cobra.MinimumNArgs(1),
	Run:  SLAdrain,
}

var startMaintenanceCmd = &cobra.Command{
	Use:   "maintenance [space separated host list]",
	Short: "Place a list of space separated Mesos Agents into maintenance mode.",
	Long: `Places Mesos Agent into Maintenance mode. Agents in this list
are de-prioritized for scheduling a task. Command
expects a space separated list of hosts to place into maintenance mode.`,
	Args: cobra.MinimumNArgs(1),
	Run:  maintenance,
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


	maintenanceMonitorPrint(hostResult, []aurora.MaintenanceMode{aurora.MaintenanceMode_DRAINED})

	if err != nil {
		log.Fatalln("error: %+v", err)
	}
}

func SLAdrain(cmd *cobra.Command, args []string) {
	log.Infoln("Setting hosts to DRAINING with SLA awareness")
	log.Infoln(args)

	policy := aurora.SlaPolicy{CountSlaPolicy: &aurora.CountSlaPolicy{Count: 1, DurationSecs: 60*30}}
	result, err := client.SLADrainHosts(&policy, 60*60, args...)
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

	maintenanceMonitorPrint(hostResult, []aurora.MaintenanceMode{aurora.MaintenanceMode_DRAINED})

	if err != nil {
		log.Fatalln("error: %+v", err)
	}

}

func maintenance(cmd *cobra.Command, args []string) {
	log.Infoln("Setting hosts to Maintenance mode")
	log.Infoln(args)
	_, result, err := client.StartMaintenance(args...)
	if err != nil {
		log.Fatalf("error: %+v\n", err)
	}

	log.Debugln(result)

	// Monitor change to DRAINING and DRAINED mode
	hostResult, err := monitor.HostMaintenance(
		args,
		[]aurora.MaintenanceMode{aurora.MaintenanceMode_SCHEDULED},
		monitorInterval,
		monitorTimeout)


	maintenanceMonitorPrint(hostResult, []aurora.MaintenanceMode{aurora.MaintenanceMode_SCHEDULED})

	if err != nil {
		log.Fatalln("error: %+v", err)
	}
}
