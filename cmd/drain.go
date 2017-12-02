package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(drainCmd)
}

var drainCmd = &cobra.Command{
	Use:   "drain [space separated host list]",
	Short: "Place a list of space separated Mesos Agents into maintenance mode.",
	Long: `Adds a Mesos Agent to Aurora's Drain list. Agents in this list
    are not allowed to schedule new tasks and any tasks already running on this Agent
    are killed and rescheduled in an Agent that is not in maintenance mode. Command
    expects a space separated list of hosts to place into maintenance mode.`,
	Args: cobra.MinimumNArgs(1),
	Run:  drain,
}

func drain(cmd *cobra.Command, args []string) {

	fmt.Print(args)
}
