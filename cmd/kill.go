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
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(killCmd)

	/* Sub-Commands */

	// Kill Job
	killCmd.AddCommand(killJobCmd)
	killCmd.AddCommand(killTasksCmd)

	killJobCmd.Flags().StringVarP(env, "environment", "e", "", "Aurora Environment")
	killJobCmd.Flags().StringVarP(role, "role", "r", "", "Aurora Role")
	killJobCmd.Flags().StringVarP(name, "name", "n", "", "Aurora Name")
	killJobCmd.Flags().BoolVarP(&monitor, "monitor", "m", true, "monitor the result after sending the command")
	killJobCmd.MarkFlagRequired("environment")
	killJobCmd.MarkFlagRequired("role")
	killJobCmd.MarkFlagRequired("name")

	//Set flags for killTask sub-command
	killTasksCmd.Flags().StringVarP(env, "environment", "e", "", "Aurora Environment")
	killTasksCmd.Flags().StringVarP(role, "role", "r", "", "Aurora Role")
	killTasksCmd.Flags().StringVarP(name, "name", "n", "", "Aurora Name")
	killTasksCmd.Flags().StringVarP(instance, "instance", "i", "", "Instance Number")
	killTasksCmd.Flags().BoolVarP(&monitor, "monitor", "m", true, "monitor the result after sending the command")
	killTasksCmd.MarkFlagRequired("environment")
	killTasksCmd.MarkFlagRequired("role")
	killTasksCmd.MarkFlagRequired("name")
	killTasksCmd.MarkFlagRequired("instance")
}

var killCmd = &cobra.Command{
	Use:   "kill",
	Short: "Kill an Aurora Job",
}

var killJobCmd = &cobra.Command{
	Use:   "job",
	Short: "Kill an Aurora Job",
	Run:   killJob,
}

var killTasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Kill Aurora Tasks",
	Run:   killTasks,
}

func killJob(cmd *cobra.Command, args []string) {
	log.Infof("Killing job [Env:%s Role:%s Name:%s]\n", *env, *role, *name)

	job := realis.NewJob().
		Environment(*env).
		Role(*role).
		Name(*name)
	err := client.KillJob(job.JobKey())
	if err != nil {
		log.Fatalln(err)
	}
	if monitor {
		if ok, err := client.MonitorInstances(job.JobKey(), 0, 5, 50); !ok || err != nil {
			log.Fatalln("Unable to kill all instances of job")
		}
	}
}

func killTasks(cmd *cobra.Command, args []string) {
	log.Infof("Killing task [Env:%s Role:%s Name:%s Instance:%s]\n", *env, *role, *name, *instance)

	//Set jobKey for the tasks to be killed.
	task := realis.NewTask().
		Environment(*env).
		Role(*role).
		Name(*name)

	//Check that instance number is passed
	if instance == nil {
		log.Fatalln("Instance number not found. Please pass a valid instance number. If you want to pass multiple instances, please pass them as space separated integer values")
	}

	/*
	* In the following block, we convert instance numbers, which were passed as strings, to integer values
	* After converting them to integers, we add them to a slice of type int32.
	 */
	var intErr error
	var instanceNumber int

	splitString := strings.Split(*instance, " ")
	instanceList := make([]int32, len(splitString))

	for i := range instanceList {
		instanceNumber, intErr = strconv.Atoi(splitString[i])
		if intErr != nil {
			log.Fatalln("Instance passed should be a number. Error: " + intErr.Error())
			break
		} else {
			instanceList[i] = int32(instanceNumber)
		}
	}

	//Call the killtasks function, passing the instanceList as the list of instances to be killed.
	_, err := client.KillInstances(task.JobKey(), instanceList...)

	if err != nil {
		log.Fatalln(err)
	}
	if monitor {
		if ok, err := client.MonitorInstances(task.JobKey(), 0, 5, 50); !ok || err != nil {
			log.Fatalln("Unable to kill the given task")
		}
	}
}
