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
	"github.com/paypal/gorealis/v2/gen-go/apache/aurora"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pulseJobUpdateCmd)

	pulseJobUpdateCmd.Flags().StringVarP(env, "environment", "e", "", "Aurora Environment")
	pulseJobUpdateCmd.Flags().StringVarP(role, "role", "r", "", "Aurora Role")
	pulseJobUpdateCmd.Flags().StringVarP(name, "name", "n", "", "Aurora Name")
	pulseJobUpdateCmd.Flags().StringVar(&updateID, "id", "", "Update ID")
}

var pulseJobUpdateCmd = &cobra.Command{
	Use:   "pulse",
	Short: "Pulse a Job update",
	Run:   pulseJobUpdate,
}

func pulseJobUpdate(cmd *cobra.Command, args []string) {
	_, err := client.PulseJobUpdate(
		aurora.JobUpdateKey{
			Job: &aurora.JobKey{Environment: *env, Role: *role, Name: *name},
			ID:  updateID,
		})

	if err != nil {
		log.Fatal(err)
	}
}
