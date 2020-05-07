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
	"strings"
	"time"

	"github.com/aurora-scheduler/australis/internal"
	"github.com/aurora-scheduler/gorealis/v2/gen-go/apache/aurora"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(monitorCmd)

	monitorCmd.AddCommand(monitorHostCmd.Cmd)

	monitorHostCmd.Cmd.Run = monitorHost
	monitorHostCmd.Cmd.Flags().DurationVar(&monitorHostCmd.MonitorInterval, "interval", time.Second*5, "Interval at which to poll scheduler.")
	monitorHostCmd.Cmd.Flags().DurationVar(&monitorHostCmd.MonitorTimeout, "timeout", time.Minute*10, "Time after which the monitor will stop polling and throw an error.")
	monitorHostCmd.Cmd.Flags().StringSliceVar(&monitorHostCmd.StatusList, "statuses", []string{aurora.MaintenanceMode_DRAINED.String()}, "List of acceptable statuses for a host to be in. (case-insensitive) [NONE, SCHEDULED, DRAINED, DRAINING]")
}

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Watch for a specific state change",
}

var monitorHostCmd = internal.MonitorCmdConfig{
	Cmd: &cobra.Command{
		Use:   "hosts",
		Short: "Watch a host maintenance status until it enters one of the desired statuses.",
		Long: `Provide a list of hosts to monitor for desired statuses. Statuses may be passed using the --statuses
flag with a list of comma separated statuses. Statuses include [NONE, SCHEDULED, DRAINED, DRAINING]`,
	},
	StatusList: make([]string, 0),
}

func monitorHost(cmd *cobra.Command, args []string) {
	maintenanceModes := make([]aurora.MaintenanceMode, 0)

	for _, status := range monitorHostCmd.StatusList {
		mode, err := aurora.MaintenanceModeFromString(strings.ToUpper(status))
		if err != nil {
			log.Fatal(err)
		}

		maintenanceModes = append(maintenanceModes, mode)
	}

	log.Infof("Monitoring for %v at %v intervals", monitorHostCmd.MonitorTimeout, monitorHostCmd.MonitorInterval)
	hostResult, err := client.MonitorHostMaintenance(args, maintenanceModes, monitorHostCmd.MonitorInterval, monitorHostCmd.MonitorTimeout)

	internal.MaintenanceMonitorPrint(hostResult, maintenanceModes, toJson)

	if err != nil {
		log.Fatal(err)
	}
}
