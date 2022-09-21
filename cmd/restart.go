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
	"strconv"
	"strings"

	realis "github.com/aurora-scheduler/gorealis/v2"

	"github.com/aurora-scheduler/gorealis/v2/gen-go/apache/aurora"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(restartCmd)

	restartCmd.AddCommand(restartJobCmd)
	restartJobCmd.Flags().StringVarP(env, "environment", "e", "", "Aurora Environment")
	restartJobCmd.Flags().StringVarP(role, "role", "r", "", "Aurora Role")
	restartJobCmd.Flags().StringVarP(name, "name", "n", "", "Aurora Name")

	restartCmd.AddCommand(restartTasksCmd)
	restartTasksCmd.Flags().StringVarP(env, "environment", "e", "", "Aurora Environment")
	restartTasksCmd.Flags().StringVarP(role, "role", "r", "", "Aurora Role")
	restartTasksCmd.Flags().StringVarP(name, "name", "n", "", "Aurora Name")
	restartTasksCmd.Flags().StringVarP(instances, "instances", "I", "", "Instances e.g. 1, 2, 5")
	restartTasksCmd.Flags().BoolVarP(&monitor, "monitor", "m", true, "monitor the result after sending the command")
	restartTasksCmd.MarkFlagRequired("environment")
	restartTasksCmd.MarkFlagRequired("role")
	restartTasksCmd.MarkFlagRequired("name")
	restartTasksCmd.MarkFlagRequired("instances")
}

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart an Aurora Job.",
}

var restartJobCmd = &cobra.Command{
	Use:   "job",
	Short: "Restart a Job.",
	Run:   restartJob,
}

var restartTasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Restart tasks for a Job.",
	Run:   restartTasks,
}

func restartJob(cmd *cobra.Command, args []string) {
	key := aurora.JobKey{Environment: *env, Role: *role, Name: *name}
	if err := client.RestartJob(key); err != nil {
		log.Fatal("unable to create Aurora job: ", err)
	}
}

func restartTasks(cmd *cobra.Command, args []string) {
	log.Infof("Restarts task [Env:%s Role:%s Name:%s Instance:%s Monitor:%s]\n", *env, *role, *name, *instances, strconv.FormatBool(monitor))

	//Set jobKey for the tasks to be killed.
	task := realis.NewTask().
		Environment(*env).
		Role(*role).
		Name(*name)

	/*
	* In the following block, we convert instance numbers, which were passed as strings, to integer values
	* After converting them to integers, we add them to a slice of type int32.
	 */

	splitString := strings.Split(*instances, ",")
	instanceList := make([]int32, len(splitString))

	for i := range instanceList {
		splitString[i] = strings.TrimSpace(splitString[i])
		var instanceNumber int
		var err error
		if instanceNumber, err = strconv.Atoi(splitString[i]); err != nil {
			log.Fatalln("Instance passed should be a number. Error: " + err.Error())
			return
		}
		instanceList[i] = int32(instanceNumber)
	}

	//Call the RestartInstances function, passing the instanceList as the list of instances to be restarted.
	if err := client.RestartInstances(task.JobKey(), instanceList...); err != nil {
		log.Fatalln(err)
	}

	if monitor {
		if ok, err := client.MonitorInstances(task.JobKey(), int32(len(instanceList)), 5, 50); !ok || err != nil {
			log.Fatalln("Monitor failed to monitor the given task after restart. Error: " + err.Error())
		}
	}

}
