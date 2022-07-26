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
	"fmt"

	"github.com/aurora-scheduler/australis/internal"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(simulateCmd)

	simulateCmd.AddCommand(fitCmd)
}

var simulateCmd = &cobra.Command{
	Use:   "simulate",
	Short: "Simulate some work based on the current cluster condition, and return the output",
}

var fitCmd = &cobra.Command{
	Use:   "fit",
	Short: "Compute how many tasks can we fit to a cluster",
	Run:   fit,
	Args:  cobra.RangeArgs(1, 2),
}

func fit(cmd *cobra.Command, args []string) {
	log.Infof("Compute how many tasks can be fit in the remaining cluster capacity")

	taskConfig, err := internal.UnmarshalTaskConfig(args[0])
	if err != nil {
		log.Fatalln(err)
	}

	offers, err := client.Offers()
	if err != nil {
		log.Fatal("error: %+v", err)
	}

	numTasks, err := client.FitTasks(taskConfig, offers)
	if err != nil {
		log.Fatal("error: %+v", err)
	}

	fmt.Println(numTasks)
}
