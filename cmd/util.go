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
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aurora-scheduler/gorealis/v2/gen-go/apache/aurora"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type monitorCmdConfig struct {
	cmd                             *cobra.Command
	monitorInterval, monitorTimeout time.Duration
	statusList                      []string
}

func toJSON(v interface{}) string {

	output, err := json.Marshal(v)

	if err != nil {
		log.Fatalln("Unable to serialize Aurora response: %+v", v)
	}

	return string(output)
}

func getLoggingLevels() string {

	var buffer bytes.Buffer

	for _, level := range logrus.AllLevels {
		buffer.WriteString(level.String())
		buffer.WriteString(" ")
	}

	buffer.Truncate(buffer.Len() - 1)

	return buffer.String()

}

func maintenanceMonitorPrint(hostResult map[string]bool, desiredStates []aurora.MaintenanceMode) {
	if len(hostResult) > 0 {
		// Create anonymous struct for JSON formatting
		output := struct {
			DesiredStates   []string `json:desired_states`
			Transitioned    []string `json:transitioned`
			NonTransitioned []string `json:non-transitioned`
		}{
			make([]string, 0),
			make([]string, 0),
			make([]string, 0),
		}

		for _, state := range desiredStates {
			output.DesiredStates = append(output.DesiredStates, state.String())
		}

		for host, ok := range hostResult {
			if ok {
				output.Transitioned = append(output.Transitioned, host)
			} else {
				output.NonTransitioned = append(output.NonTransitioned, host)
			}
		}

		if toJson {
			fmt.Println(toJSON(output))
		} else {
			fmt.Printf("Entered %v status: %v\n", output.DesiredStates, output.Transitioned)
			fmt.Printf("Did not enter %v status: %v\n", output.DesiredStates, output.NonTransitioned)
		}
	}
}
