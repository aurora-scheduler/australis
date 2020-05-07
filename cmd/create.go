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
	"github.com/aurora-scheduler/australis/internal"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an Aurora Job",
	Run:   createJob,
	Args:  cobra.ExactArgs(1),
}

func createJob(cmd *cobra.Command, args []string) {
	job, err := internal.UnmarshalJob(args[0])

	if err != nil {
		log.Fatalln(err)
	}

	auroraJob, err := job.ToRealis()
	if err != nil {
		log.Fatalln(err)
	}

	if err := client.CreateJob(auroraJob); err != nil {
		log.Fatal("unable to create Aurora job: ", err)
	}

	if ok, monitorErr := client.MonitorInstances(auroraJob.JobKey(),
		auroraJob.GetInstanceCount(),
		5,
		50); !ok || monitorErr != nil {
		if err := client.KillJob(auroraJob.JobKey()); err != nil {
			log.Fatal(monitorErr, err)
		}
		log.Fatal(monitorErr)
	}

}
