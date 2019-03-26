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
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(completionCmd)

	completionCmd.Hidden = true
	completionCmd.Flags().StringVar(&filename, "filename", "australis.completion.sh", "Path and name of the autocompletion file.")
}

var completionCmd = &cobra.Command{
	Use:   "autocomplete",
	Short: "Create auto completion for bash.",
	Long: `Create auto completion bash file for australis. Auto completion file must be placed in the correct
directory in order for bash to pick up the definitions.

Copy australis.completion.sh into the correct folder and rename to australis

In Linux, this directory is usually /etc/bash_completion.d/
In MacOS this directory is $(brew --prefix)/etc/bash_completion.d if auto completion was install through brew.
`,
	PersistentPreRun:  func(cmd *cobra.Command, args []string) {}, // We don't need a realis client for this cmd
	PersistentPostRun: func(cmd *cobra.Command, args []string) {}, // We don't need a realis client for this cmd
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenBashCompletionFile(filename)
	},
}

