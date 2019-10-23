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
	rootCmd.AddCommand(resumeJobUpdateCmd)

	resumeJobUpdateCmd.Flags().StringVarP(env, "environment", "e", "", "Aurora Environment")
	resumeJobUpdateCmd.Flags().StringVarP(role, "role", "r", "", "Aurora Role")
	resumeJobUpdateCmd.Flags().StringVarP(name, "name", "n", "", "Aurora Name")
	resumeJobUpdateCmd.Flags().StringVar(&updateID, "id", "", "Update ID")
	resumeJobUpdateCmd.Flags().StringVar(message, "message", "", "Message to store along resume.")
}

var resumeJobUpdateCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resume a Job update",
	Run:   resumeJobUpdate,
}

func resumeJobUpdate(cmd *cobra.Command, args []string) {
	err := client.ResumeJobUpdate(
		&aurora.JobUpdateKey{
			Job: &aurora.JobKey{Environment: *env, Role: *role, Name: *name},
			ID:  updateID,
		},
		*message)

	if err != nil {
		log.Fatal(err)
	}
}
