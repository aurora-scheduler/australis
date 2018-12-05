package cmd

import (
	"github.com/paypal/gorealis/gen-go/apache/aurora"
	"github.com/spf13/cobra"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	rootCmd.AddCommand(startCmd)


	// Sub-commands
	startCmd.AddCommand(startDrainCmd)

	// Maintenance specific flags
	startDrainCmd.Flags().DurationVar(&monitorInterval,"interval", time.Second * 5, "Interval at which to poll scheduler.")
	startDrainCmd.Flags().DurationVar(&monitorTimeout,"timeout", time.Minute * 10, "Time after which the monitor will stop polling and throw an error.")

	startCmd.AddCommand(startSLADrainCmd)


	/* SLA Aware commands */
	startSLADrainCmd.AddCommand(startSLACountDrainCmd)

	// SLA Maintenance specific flags
	startSLACountDrainCmd.Flags().DurationVar(&monitorInterval,"interval", time.Second * 5, "Interval at which to poll scheduler.")
	startSLACountDrainCmd.Flags().DurationVar(&monitorTimeout,"timeout", time.Minute * 20, "Time after which the monitor will stop polling and throw an error.")
	startSLACountDrainCmd.Flags().Int64Var(&count, "count", 5, "Instances count that should be running to meet SLA.")
	startSLACountDrainCmd.Flags().DurationVar(&duration, "duration", time.Minute * 10, "Window of time from which we derive the SLA.")

	startSLADrainCmd.AddCommand(startSLAPercentageDrainCmd)

	// SLA Maintenance specific flags
	startSLAPercentageDrainCmd.Flags().DurationVar(&monitorInterval,"interval", time.Second * 5, "Interval at which to poll scheduler.")
	startSLAPercentageDrainCmd.Flags().DurationVar(&monitorTimeout,"timeout", time.Minute * 1, "Time after which the monitor will stop polling and throw an error.")
    startSLAPercentageDrainCmd.Flags().Float64Var(&percent, "percent", 75.0, "Percentage of instances that should be running to meet SLA.")
	startSLAPercentageDrainCmd.Flags().DurationVar(&duration, "duration", time.Minute * 10, "Window of time from which we derive the SLA.")

	startCmd.AddCommand(startMaintenanceCmd)

	// SLA Maintenance specific flags
	startMaintenanceCmd.Flags().DurationVar(&monitorInterval,"interval", time.Second * 5, "Interval at which to poll scheduler.")
	startMaintenanceCmd.Flags().DurationVar(&monitorTimeout,"timeout", time.Minute * 1, "Time after which the monitor will stop polling and throw an error.")
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
	Use:   "sla-drain",
	Short: "Place a list of space separated Mesos Agents into maintenance mode using SLA aware strategies.",
	Long: `Adds a Mesos Agent to Aurora's Drain list. Agents in this list
are not allowed to schedule new tasks and any tasks already running on this Agent
are killed and rescheduled in an Agent that is not in maintenance mode. Command
expects a space separated list of hosts to place into maintenance mode.`,
}

var startSLACountDrainCmd = &cobra.Command{
	Use:   "count [space separated host list]",
	Short: "Place a list of space separated Mesos Agents into maintenance mode using the count SLA aware policy as a fallback.",
	Long: `Adds a Mesos Agent to Aurora's Drain list. Tasks will be drained using the count SLA policy as a fallback
when a Job does not have a defined SLA policy.`,
	Args: cobra.MinimumNArgs(1),
	Run:  SLACountDrain,
}

var startSLAPercentageDrainCmd = &cobra.Command{
	Use:   "percentage [space separated host list]",
	Short: "Place a list of space separated Mesos Agents into maintenance mode using the percentage SLA aware policy as a fallback.",
	Long: `Adds a Mesos Agent to Aurora's Drain list. Tasks will be drained using the percentage SLA policy as a fallback
when a Job does not have a defined SLA policy.`,
	Args: cobra.MinimumNArgs(1),
	Run:  SLAPercentageDrain,
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
		int(monitorInterval.Seconds()),
		int(monitorTimeout.Seconds()))

	maintenanceMonitorPrint(hostResult, []aurora.MaintenanceMode{aurora.MaintenanceMode_DRAINED})

	if err != nil {
		log.Fatalln("error: %+v", err)
	}
}

func slaDrain(policy *aurora.SlaPolicy, hosts ...string) {

	result, err := client.SLADrainHosts(policy, 60*60, hosts...)
	if err != nil {
		log.Fatalf("error: %+v\n", err)
	}

	log.Debugln(result)

	// Monitor change to DRAINING and DRAINED mode
	hostResult, err := monitor.HostMaintenance(
		hosts,
		[]aurora.MaintenanceMode{aurora.MaintenanceMode_DRAINED},
		int(monitorInterval.Seconds()),
		int(monitorTimeout.Seconds()))

	maintenanceMonitorPrint(hostResult, []aurora.MaintenanceMode{aurora.MaintenanceMode_DRAINED})

	if err != nil {
		log.Fatalln("error: %+v", err)
	}
}

func SLACountDrain(cmd *cobra.Command, args []string) {
	log.Infoln("Setting hosts to DRAINING with the Count SLA policy.")
	log.Infoln(args)

	slaDrain(&aurora.SlaPolicy{
		CountSlaPolicy: &aurora.CountSlaPolicy{Count: count, DurationSecs: int64(duration.Seconds())}},
	args...)
}

func SLAPercentageDrain(cmd *cobra.Command, args []string) {
	log.Infoln("Setting hosts to DRAINING with the Percentage SLA policy.")
	log.Infoln(args)

	slaDrain(&aurora.SlaPolicy{
		PercentageSlaPolicy: &aurora.PercentageSlaPolicy{Percentage: percent, DurationSecs: int64(duration.Seconds())}},
	args...)
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
		int(monitorInterval.Seconds()),
		int(monitorTimeout.Seconds()))


	maintenanceMonitorPrint(hostResult, []aurora.MaintenanceMode{aurora.MaintenanceMode_SCHEDULED})

	if err != nil {
		log.Fatalln("error: %+v", err)
	}
}
