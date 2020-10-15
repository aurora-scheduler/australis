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
	"github.com/aurora-scheduler/gorealis/v2/gen-go/apache/aurora"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(restartCmd)

	restartCmd.AddCommand(restartJobCmd)
	restartJobCmd.Flags().StringVarP(env, "environment", "e", "", "Aurora Environment")
	restartJobCmd.Flags().StringVarP(role, "role", "r", "", "Aurora Role")
	restartJobCmd.Flags().StringVarP(name, "name", "n", "", "Aurora Name")
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

func restartJob(cmd *cobra.Command, args []string) {
	key := aurora.JobKey{Environment: *env, Role: *role, Name: *name}
	if err := client.RestartJob(key); err != nil {
		log.Fatal("unable to create Aurora job: ", err)
	}
}
