package cmd

import (
	"fmt"
	"github.com/paypal/gorealis/gen-go/apache/aurora"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

func init() {
	rootCmd.AddCommand(stopCmd)

	// Stop subcommands
	stopCmd.AddCommand(stopMaintCmd)
	stopMaintCmd.Flags().IntVar(&monitorInterval,"interval", 5, "Interval at which to poll scheduler.")
	stopMaintCmd.Flags().IntVar(&monitorTimeout,"timeout", 50, "Time after which the monitor will stop polling and throw an error.")

	// Stop update

	stopCmd.AddCommand(stopUpdateCmd)
	stopUpdateCmd.Flags().StringVarP(env, "environment", "e", "", "Aurora Environment")
	stopUpdateCmd.Flags().StringVarP(role, "role", "r", "", "Aurora Role")
	stopUpdateCmd.Flags().StringVarP(name, "name", "n", "", "Aurora Name")

}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a service or maintenance on a host (DRAIN).",
}

var stopMaintCmd = &cobra.Command{
	Use:   "drain [space separated host list]",
	Short: "Stop maintenance on a host (move to NONE).",
	Long:  `Transition a list of hosts currently in a maintenance status out of it.`,
	Run:   endMaintenance,
}

var stopUpdateCmd = &cobra.Command{
	Use:   "update [update ID]",
	Short: "Stop update",
	Long:  `To be written.`,
	Run:   stopUpdate,
}

func endMaintenance(cmd *cobra.Command, args []string) {
	log.Println("Setting hosts to NONE maintenance status.")
	log.Println(args)
	_, result, err := client.EndMaintenance(args...)
	if err != nil {
		log.Fatalf("error: %+v\n", err)
	}

	log.Debugln(result)

	// Monitor change to NONE mode
	hostResult, err := monitor.HostMaintenance(
		args,
		[]aurora.MaintenanceMode{aurora.MaintenanceMode_NONE},
		monitorInterval,
		monitorTimeout)


	maintenanceMonitorPrint(hostResult,[]aurora.MaintenanceMode{aurora.MaintenanceMode_NONE})

	if err != nil {
		log.Fatalln("error: %+v", err)
	}
}

func stopUpdate(cmd *cobra.Command, args []string) {

	if len(args) != 1 {
		log.Fatalln("Only a single update ID must be provided.")
	}

	log.Infof("Stopping (aborting) update [%s/%s/%s] %s\n", *env, *role, *name, args[0])

	resp, err := client.AbortJobUpdate(aurora.JobUpdateKey{
		Job: &aurora.JobKey{Environment: *env, Role: *role, Name: *name},
		ID:  args[0],
	},
		"")

	if err != nil {
		log.Fatalln(err)
	}

	if toJson{
		fmt.Println(toJSON(resp.GetResult_()))
	} else {
		fmt.Println(resp.GetDetails())
	}
}
