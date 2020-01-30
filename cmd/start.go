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
	"time"

	"github.com/paypal/gorealis/v2/gen-go/apache/aurora"
	"github.com/spf13/cobra"
)

const countFlag = "count"
const percentageFlag = "percentage"

func init() {
	rootCmd.AddCommand(startCmd)

	// Sub-commands
	startCmd.AddCommand(startDrainCmd.cmd)
	startDrainCmd.cmd.Run = drain

	// Maintenance specific flags
	startDrainCmd.cmd.Flags().DurationVar(&startDrainCmd.monitorInterval, "interval", time.Second*5, "Interval at which to poll scheduler.")
	startDrainCmd.cmd.Flags().DurationVar(&startDrainCmd.monitorTimeout, "timeout", time.Minute*10, "Time after which the monitor will stop polling and throw an error.")

	/* SLA Aware commands */
	startCmd.AddCommand(startSLADrainCmd.cmd)
	startSLADrainCmd.cmd.Run = slaDrain

	// SLA Maintenance specific flags
	startSLADrainCmd.cmd.Flags().Int64Var(&count, countFlag, 5, "Instances count that should be running to meet SLA.")
	startSLADrainCmd.cmd.Flags().Float64Var(&percent, percentageFlag, 80.0, "Percentage of instances that should be running to meet SLA.")
	startSLADrainCmd.cmd.Flags().DurationVar(&duration, "duration", time.Minute*1, "Minimum time duration a task needs to be `RUNNING` to be treated as active.")
	startSLADrainCmd.cmd.Flags().DurationVar(&forceDrainTimeout, "sla-limit", time.Minute*60, "Time limit after which SLA-Aware drain sheds SLA Awareness.")
	startSLADrainCmd.cmd.Flags().DurationVar(&startSLADrainCmd.monitorInterval, "interval", time.Second*10, "Interval at which to poll scheduler.")
	startSLADrainCmd.cmd.Flags().DurationVar(&startSLADrainCmd.monitorTimeout, "timeout", time.Minute*20, "Time after which the monitor will stop polling and throw an error.")

	startCmd.AddCommand(startMaintenanceCmd.cmd)
	startMaintenanceCmd.cmd.Run = maintenance

	// SLA Maintenance specific flags
	startMaintenanceCmd.cmd.Flags().DurationVar(&startMaintenanceCmd.monitorInterval, "interval", time.Second*5, "Interval at which to poll scheduler.")
	startMaintenanceCmd.cmd.Flags().DurationVar(&startMaintenanceCmd.monitorTimeout, "timeout", time.Minute*10, "Time after which the monitor will stop polling and throw an error.")
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a service, maintenance on a host (DRAIN), a snapshot, or a backup.",
}

var startDrainCmd = monitorCmdConfig{
	cmd: &cobra.Command{
		Use:   "drain [space separated host list]",
		Short: "Place a list of space separated Mesos Agents into draining mode.",
		Long: `Adds a Mesos Agent to Aurora's Drain list. Agents in this list
are not allowed to schedule new tasks and any tasks already running on this Agent
are killed and rescheduled in an Agent that is not in maintenance mode. Command
expects a space separated list of hosts to place into maintenance mode.`,
		Args: cobra.MinimumNArgs(1),
	},
}

var startSLADrainCmd = monitorCmdConfig{
	cmd: &cobra.Command{
		Use:   "sla-drain [space separated host list]",
		Short: "Place a list of space separated Mesos Agents into maintenance mode using SLA aware strategies.",
		Long: `Adds a Mesos Agent to Aurora's Drain list. Agents in this list
are not allowed to schedule new tasks and any tasks already running on this Agent
are killed and rescheduled in an Agent that is not in maintenance mode. Command
expects a space separated list of hosts to place into maintenance mode.
If the --count argument is passed, tasks will be drained using the count SLA policy as a fallback
when a Job does not have a defined SLA policy.
If the --percentage argument is passed, tasks will be drained using the percentage SLA policy as a fallback
when a Job does not have a defined SLA policy.`,
		Args: cobra.MinimumNArgs(1),
	},
}

var startMaintenanceCmd = monitorCmdConfig{
	cmd: &cobra.Command{
		Use:   "maintenance [space separated host list]",
		Short: "Place a list of space separated Mesos Agents into maintenance mode.",
		Long: `Places Mesos Agent into Maintenance mode. Agents in this list
are de-prioritized for scheduling a task. Command
expects a space separated list of hosts to place into maintenance mode.`,
		Args: cobra.MinimumNArgs(1),
	},
}

func drain(cmd *cobra.Command, args []string) {
	log.Infoln("Setting hosts to DRAINING")
	log.Infoln(args)
	result, err := client.DrainHosts(args...)
	if err != nil {
		log.Fatalf("error: %+v", err)
	}

	log.Debugln(result)

	log.Infof("Monitoring for %v at %v intervals", monitorHostCmd.monitorTimeout, monitorHostCmd.monitorInterval)
	// Monitor change to DRAINING and DRAINED mode
	hostResult, err := client.HostMaintenanceMonitor(
		args,
		[]aurora.MaintenanceMode{aurora.MaintenanceMode_DRAINED},
		startDrainCmd.monitorInterval,
		startDrainCmd.monitorTimeout)

	maintenanceMonitorPrint(hostResult, []aurora.MaintenanceMode{aurora.MaintenanceMode_DRAINED})

	if err != nil {
		log.Fatalln(err)
	}
}

func slaDrainHosts(policy *aurora.SlaPolicy, interval, timeout time.Duration, hosts ...string) {

	result, err := client.SLADrainHosts(policy, int64(forceDrainTimeout.Seconds()), hosts...)
	if err != nil {
		log.Fatalf("error: %+v\n", err)
	}

	log.Debugln(result)

	log.Infof("Monitoring for %v at %v intervals", timeout, interval)
	// Monitor change to DRAINING and DRAINED mode
	hostResult, err := client.HostMaintenanceMonitor(
		hosts,
		[]aurora.MaintenanceMode{aurora.MaintenanceMode_DRAINED},
		interval,
		timeout)

	maintenanceMonitorPrint(hostResult, []aurora.MaintenanceMode{aurora.MaintenanceMode_DRAINED})

	if err != nil {
		log.Fatalf("error: %+v", err)
	}
}
func slaDrain(cmd *cobra.Command, args []string) {
	// This check makes sure only a single flag is set.
	// If they're both set or both not set, the statement will evaluate to true.
	if cmd.Flags().Changed(percentageFlag) == cmd.Flags().Changed(countFlag) {
		log.Fatal("Either percentage or count must be set exclusively.")
	}

	policy := &aurora.SlaPolicy{}

	if cmd.Flags().Changed(percentageFlag) {
		log.Infoln("Setting hosts to DRAINING with the Percentage SLA policy.")
		policy.PercentageSlaPolicy = &aurora.PercentageSlaPolicy{
			Percentage: percent,
			DurationSecs: int64(duration.Seconds()),
		}
	}

	if cmd.Flags().Changed(countFlag) {
		log.Infoln("Setting hosts to DRAINING with the Count SLA policy.")
		policy.CountSlaPolicy = &aurora.CountSlaPolicy{Count: count, DurationSecs: int64(duration.Seconds())}
	}

	log.Infoln("Hosts affected: ", args)
	slaDrainHosts(policy, startDrainCmd.monitorInterval, startDrainCmd.monitorTimeout, args...)
}

func maintenance(cmd *cobra.Command, args []string) {
	log.Infoln("Setting hosts to Maintenance mode")
	log.Infoln(args)
	result, err := client.StartMaintenance(args...)
	if err != nil {
		log.Fatalf("error: %+v\n", err)
	}

	log.Debugln(result)

	log.Infof("Monitoring for %v at %v intervals", monitorHostCmd.monitorTimeout, monitorHostCmd.monitorInterval)

	// Monitor change to DRAINING and DRAINED mode
	hostResult, err := client.HostMaintenanceMonitor(
		args,
		[]aurora.MaintenanceMode{aurora.MaintenanceMode_SCHEDULED},
		startMaintenanceCmd.monitorInterval,
		startMaintenanceCmd.monitorTimeout)

	maintenanceMonitorPrint(hostResult, []aurora.MaintenanceMode{aurora.MaintenanceMode_SCHEDULED})

	if err != nil {
		log.Fatalln("error: %+v", err)
	}
}
