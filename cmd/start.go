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
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"time"

	"github.com/paypal/gorealis/v2/gen-go/apache/aurora"
	"github.com/spf13/cobra"
)

const countFlag = "count"
const percentageFlag = "percentage"
const jsonFlag = "json"
const jsonFileFlag = "json-file"

func init() {
	rootCmd.AddCommand(startCmd)

	// Sub-commands
	startCmd.AddCommand(startDrainCmd.cmd)
	startDrainCmd.cmd.Run = drain

	// Maintenance specific flags
	startDrainCmd.cmd.Flags().DurationVar(&startDrainCmd.monitorInterval, "interval", time.Second*5, "Interval at which to poll scheduler.")
	startDrainCmd.cmd.Flags().DurationVar(&startDrainCmd.monitorTimeout, "timeout", time.Minute*10, "Time after which the monitor will stop polling and throw an error.")
	startDrainCmd.cmd.Flags().StringVar(&fromJsonFile, jsonFileFlag, "", "JSON file to read list of agents from.")
	startDrainCmd.cmd.Flags().BoolVar(&fromJson, jsonFlag, false, "Read JSON list of agents from the STDIN.")

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
	startSLADrainCmd.cmd.Flags().StringVar(&fromJsonFile, jsonFileFlag, "", "JSON file to read list of agents from.")
	startSLADrainCmd.cmd.Flags().BoolVar(&fromJson, jsonFlag, false, "Read JSON list of agents from the STDIN.")

	startCmd.AddCommand(startMaintenanceCmd.cmd)
	startMaintenanceCmd.cmd.Run = maintenance

	// SLA Maintenance specific flags
	startMaintenanceCmd.cmd.Flags().DurationVar(&startMaintenanceCmd.monitorInterval, "interval", time.Second*5, "Interval at which to poll scheduler.")
	startMaintenanceCmd.cmd.Flags().DurationVar(&startMaintenanceCmd.monitorTimeout, "timeout", time.Minute*10, "Time after which the monitor will stop polling and throw an error.")
	startMaintenanceCmd.cmd.Flags().StringVar(&fromJsonFile, jsonFileFlag, "", "JSON file to read list of agents from.")
	startMaintenanceCmd.cmd.Flags().BoolVar(&fromJson, jsonFlag, false, "Read JSON list of agents from the STDIN.")
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a service, maintenance on a host (DRAIN), a snapshot, or a backup.",
}

var startDrainCmd = monitorCmdConfig{
	cmd: &cobra.Command{
		Use:   "drain [space separated host list or use JSON flags]",
		Short: "Place a list of space separated Mesos Agents into draining mode.",
		Long: `Adds a Mesos Agent to Aurora's Drain list. Agents in this list
are not allowed to schedule new tasks and any tasks already running on this Agent
are killed and rescheduled in an Agent that is not in maintenance mode. Command
expects a space separated list of hosts to place into maintenance mode.`,
		Args: argsValidateJSONFlags,
	},
}

var startSLADrainCmd = monitorCmdConfig{
	cmd: &cobra.Command{
		Use:   "sla-drain [space separated host list or use JSON flags]",
		Short: "Place a list of space separated Mesos Agents into maintenance mode using SLA aware strategies.",
		Long: `Adds a Mesos Agent to Aurora's Drain list. Agents in this list
are not allowed to schedule new tasks and any tasks already running on this Agent
are killed and rescheduled in an Agent that is not in maintenance mode. Command
expects a space separated list of hosts to place into maintenance mode.
If the --count argument is passed, tasks will be drained using the count SLA policy as a fallback
when a Job does not have a defined SLA policy.
If the --percentage argument is passed, tasks will be drained using the percentage SLA policy as a fallback
when a Job does not have a defined SLA policy.`,
		Args: argsValidateJSONFlags,
	},
}

var startMaintenanceCmd = monitorCmdConfig{
	cmd: &cobra.Command{
		Use:   "maintenance [space separated host list or use JSON flags]",
		Short: "Place a list of space separated Mesos Agents into maintenance mode.",
		Long: `Places Mesos Agent into Maintenance mode. Agents in this list
are de-prioritized for scheduling a task. Command
expects a space separated list of hosts to place into maintenance mode.`,
		Args: argsValidateJSONFlags,
	},
}

func argsValidateJSONFlags(cmd *cobra.Command, args []string) error {
	if cmd.Flags().Changed(jsonFlag) && cmd.Flags().Changed(jsonFileFlag) {
		return errors.New("only json file or json stdin must be set")
	}
	// These two flags are mutually exclusive
	if cmd.Flags().Changed(jsonFlag) != cmd.Flags().Changed(jsonFileFlag) {
		return nil
	}

	if len(args) < 1 {
		return errors.New("at least one host must be specified")
	}
	return nil
}

func hostList(cmd *cobra.Command, args []string) []string {
	var hosts []string
	if cmd.Flags().Changed(jsonFlag) {
		err := json.NewDecoder(os.Stdin).Decode(&hosts)
		if err != nil {
			log.Fatal(err)
		}
	} else if cmd.Flags().Changed(jsonFileFlag) {
		data, err := ioutil.ReadFile(fromJsonFile)
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(data, &hosts)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		hosts = args
	}

	return hosts
}

func drain(cmd *cobra.Command, args []string) {
	hosts := hostList(cmd, args)

	log.Infoln("Setting hosts to DRAINING")
	log.Infoln(hosts)
	result, err := client.DrainHosts(hosts...)
	if err != nil {
		log.Fatalf("error: %+v", err)
	}

	log.Debugln(result)

	log.Infof("Monitoring for %v at %v intervals", monitorHostCmd.monitorTimeout, monitorHostCmd.monitorInterval)
	// Monitor change to DRAINING and DRAINED mode
	hostResult, err := client.MonitorHostMaintenance(
		hosts,
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
	hostResult, err := client.MonitorHostMaintenance(
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
	hosts := hostList(cmd, args)

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
	slaDrainHosts(policy, startDrainCmd.monitorInterval, startDrainCmd.monitorTimeout, hosts...)
}

func maintenance(cmd *cobra.Command, args []string) {
	hosts := hostList(cmd, args)

	log.Infoln("Setting hosts to Maintenance mode")
	log.Infoln(hosts)
	result, err := client.StartMaintenance(hosts...)
	if err != nil {
		log.Fatalf("error: %+v", err)
	}

	log.Debugln(result)

	log.Infof("Monitoring for %v at %v intervals", monitorHostCmd.monitorTimeout, monitorHostCmd.monitorInterval)

	// Monitor change to DRAINING and DRAINED mode
	hostResult, err := client.MonitorHostMaintenance(
		hosts,
		[]aurora.MaintenanceMode{aurora.MaintenanceMode_SCHEDULED},
		startMaintenanceCmd.monitorInterval,
		startMaintenanceCmd.monitorTimeout)

	maintenanceMonitorPrint(hostResult, []aurora.MaintenanceMode{aurora.MaintenanceMode_SCHEDULED})

	if err != nil {
		log.Fatalln("error: %+v", err)
	}
}
