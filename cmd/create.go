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

	"github.com/aurora-scheduler/australis/internal"
	realis "github.com/aurora-scheduler/gorealis/v2"
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

	auroraJob := realis.NewJob().
		Environment(job.Environment).
		Role(job.Role).
		Name(job.Name).
		CPU(job.CPU).
		RAM(job.RAM).
		Disk(job.Disk).
		IsService(job.Service).
		InstanceCount(job.Instances)

	// Adding URIs.
	for _, uri := range job.URIs {
		auroraJob.AddURIs(uri.Extract, uri.Cache, uri.URI)
	}

	// Adding Metadata.
	for key, value := range job.Metadata {
		auroraJob.AddLabel(key, value)
	}

	// If thermos jobs processes are provided, use them
	if len(job.Thermos) > 0 {
		thermosExec := realis.ThermosExecutor{}
		for _, process := range job.Thermos {
			thermosExec.AddProcess(realis.NewThermosProcess(process.Name, process.Cmd))
		}
		auroraJob.ThermosExecutor(thermosExec)
	} else if job.Executor.Name != "" {
		// Non-Thermos executor
		if job.Executor.Name == "" {
			log.Fatal("no executor name provided")
		}

		auroraJob.ExecutorName(job.Executor.Name)
		auroraJob.ExecutorData(job.Executor.Data)
	} else if job.Container != nil {
		if job.Container.Docker == nil {
			log.Fatal("no container specified")
		}

		if job.Container.Docker.Tag != "" && !strings.ContainsRune(job.Container.Docker.Name, ':') {
			job.Container.Docker.Name += ":" + job.Container.Docker.Tag
		}
		auroraJob.Container(realis.NewDockerContainer().Image(job.Container.Docker.Name))

	} else {
		log.Fatal("task does not contain a thermos definition, a custom executor name, or a container to launch")
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
