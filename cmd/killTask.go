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
	realis "github.com/aurora-scheduler/gorealis/v2"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(killTask)

	/* Sub-Commands */

	// Kill Job
	killTask.AddCommand(killTaskCmd)

	killTaskCmd.Flags().StringVarP(env, "environment", "e", "", "Aurora Environment")
	killTaskCmd.Flags().StringVarP(role, "role", "r", "", "Aurora Role")
	killTaskCmd.Flags().StringVarP(name, "name", "n", "", "Aurora Name")
	killTaskCmd.Flags().BoolVarP(&monitor, "monitor", "m", true, "monitor the result after sending the command")
	killTaskCmd.MarkFlagRequired("environment")
	killTaskCmd.MarkFlagRequired("role")
	killTaskCmd.MarkFlagRequired("name")
}

var killTask = &cobra.Command{
	Use:   "killTask",
	Short: "Kill an Aurora Task",
}

var killTaskCmd = &cobra.Command{
	Use:   "task",
	Short: "Kill an Aurora Task",
	Run:   killTask,
}

func killTask(cmd *cobra.Command, args []string) {
	log.Infof("Killing task [Env:%s Role:%s Name:%s]\n", *env, *role, *name)

	task := realis.NewTask().
		Environment(*env).
		Role(*role).
		Name(*name)
	err := client.KillTask(task.JobKey())
	if err != nil {
		log.Fatalln(err)
	}
	if monitor {
		if ok, err := client.MonitorInstances(task.JobKey(), 0, 5, 50); !ok || err != nil {
			log.Fatalln("Unable to kill the given task")
		}
	}
}
