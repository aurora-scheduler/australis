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

	"github.com/aurora-scheduler/gorealis/v2/gen-go/apache/aurora"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rollbackCmd)

	rollbackCmd.AddCommand(rollbackUpdateCmd)
	rollbackUpdateCmd.Flags().StringVarP(env, "environment", "e", "", "Aurora Environment")
	rollbackUpdateCmd.Flags().StringVarP(role, "role", "r", "", "Aurora Role")
	rollbackUpdateCmd.Flags().StringVarP(name, "name", "n", "", "Aurora Name")
	rollbackUpdateCmd.Flags().StringVar(&updateID, "id", "", "Update ID")
	rollbackUpdateCmd.Flags().StringVar(message, "message", "", "Message to store alongside resume event")
	rollbackUpdateCmd.MarkFlagRequired("environment")
	rollbackUpdateCmd.MarkFlagRequired("role")
	rollbackUpdateCmd.MarkFlagRequired("name")
	rollbackUpdateCmd.MarkFlagRequired("id")
}

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback an operation such as an Update",
}

var rollbackUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Rollback an update",
	Run:   rollbackUpdate,
}

func rollbackUpdate(cmd *cobra.Command, args []string) {
	var updateMessage string
	if message != nil {
		updateMessage = *message
	}
	err := client.RollbackJobUpdate(aurora.JobUpdateKey{
		Job: &aurora.JobKey{Environment: *env, Role: *role, Name: *name},
		ID:  updateID,
	}, updateMessage)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Rollback update for update ID %v sent successfully\n", updateID)
}
