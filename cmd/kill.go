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
	realis "github.com/paypal/gorealis/v2"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(killCmd)

	/* Sub-Commands */

	// Kill Job
	killCmd.AddCommand(killJobCmd)

	killJobCmd.Flags().StringVarP(env, "environment", "e", "", "Aurora Environment")
	killJobCmd.Flags().StringVarP(role, "role", "r", "", "Aurora Role")
	killJobCmd.Flags().StringVarP(name, "name", "n", "", "Aurora Name")
	killJobCmd.MarkFlagRequired("environment")
	killJobCmd.MarkFlagRequired("role")
	killJobCmd.MarkFlagRequired("name")
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

	if ok, err := client.MonitorInstances(job.JobKey(), 0, 5, 50); !ok || err != nil {
		log.Fatalln("Unable to kill all instances of job")
	}
}
