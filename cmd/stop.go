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

	"github.com/aurora-scheduler/gorealis/v2/gen-go/apache/aurora"
	"github.com/spf13/cobra"
)

var stopMaintenanceConfig = monitorCmdConfig{}

func init() {
	rootCmd.AddCommand(stopCmd)

	// Stop subcommands
	stopCmd.AddCommand(stopMaintCmd.cmd)
	stopMaintCmd.cmd.Run = endMaintenance
	stopMaintCmd.cmd.Flags().DurationVar(&stopMaintenanceConfig.monitorInterval, "interval", time.Second*5, "Interval at which to poll scheduler.")
	stopMaintCmd.cmd.Flags().DurationVar(&stopMaintenanceConfig.monitorTimeout, "timeout", time.Minute*1, "Time after which the monitor will stop polling and throw an error.")

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

var stopMaintCmd = monitorCmdConfig{
	cmd: &cobra.Command{
		Use:   "drain [space separated host list]",
		Short: "Stop maintenance on a host (move to NONE).",
		Long:  `Transition a list of hosts currently in a maintenance status out of it.`,
	},
}

var stopUpdateCmd = &cobra.Command{
	Use:   "update [update ID]",
	Short: "Stop update",
	Long:  `To be written.`,
	Run:   stopUpdate,
}

func endMaintenance(cmd *cobra.Command, args []string) {
	log.Println("Setting hosts to NONE maintenance status.")
	log.Println(args)
	result, err := client.EndMaintenance(args...)
	if err != nil {
		log.Fatalf("error: %+v", err)
	}

	log.Debugln(result)

	// Monitor change to NONE mode
	hostResult, err := client.MonitorHostMaintenance(
		args,
		[]aurora.MaintenanceMode{aurora.MaintenanceMode_NONE},
		stopMaintenanceConfig.monitorInterval,
		stopMaintenanceConfig.monitorTimeout)

	maintenanceMonitorPrint(hostResult, []aurora.MaintenanceMode{aurora.MaintenanceMode_NONE})

	if err != nil {
		log.Fatalln("error: %+v", err)
	}
}

func stopUpdate(cmd *cobra.Command, args []string) {

	if len(args) != 1 {
		log.Fatalln("Only a single update ID must be provided.")
	}

	log.Infof("Stopping (aborting) update [%s/%s/%s] %s\n", *env, *role, *name, args[0])

	err := client.AbortJobUpdate(aurora.JobUpdateKey{
		Job: &aurora.JobKey{Environment: *env, Role: *role, Name: *name},
		ID:  args[0],
	},
		"")

	if err != nil {
		log.Fatalln(err)
	}
}
