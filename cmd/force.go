package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

func init() {
	rootCmd.AddCommand(forceCmd)

	// Sub-commands
	forceCmd.AddCommand(forceSnapshotCmd)
    forceCmd.AddCommand(forceBackupCmd)

	// Recon sub-commands
	forceCmd.AddCommand(reconCmd)
	reconCmd.AddCommand(forceImplicitReconCmd)
	reconCmd.AddCommand(forceExplicitReconCmd)
}

var forceCmd = &cobra.Command{
	Use:   "force",
	Short: "Force the scheduler to do a snapshot, a backup, or a task reconciliation.",
}

var forceSnapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Force the leading scheduler to perform a Snapshot.",
	Long: `Takes a Snapshot of the in memory state of the Apache Aurora cluster and
writes it to the Mesos replicated log. This should NOT be confused with a backup.`,
	Run:  snapshot,
}

func snapshot(cmd *cobra.Command, args []string) {
	fmt.Println("Forcing scheduler to write snapshot to Mesos replicated log")
	err := client.Snapshot()
	if err != nil {
		log.Fatalf("error: %+v\n", err)
	} else {
		log.Println("Snapshot started successfully")
	}
}

var forceBackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Force the leading scheduler to perform a Backup.",
	Long: `Force the Aurora Scheduler to write a backup of the latest snapshot to the filesystem 
of the leading scheduler.`,
	Run:  backup,
}

func backup(cmd *cobra.Command, args []string) {
	fmt.Println("Forcing scheduler to write a Backup of latest Snapshot to file system")
	err := client.PerformBackup()
	if err != nil {
		log.Fatalf("error: %+v\n", err)
	} else {
		log.Println("Backup started successfully")
	}
}

var reconCmd = &cobra.Command{
	Use:   "recon",
	Short: "Force the leading scheduler to perform a reconciliation.",
	Long: `Force the Aurora Scheduler to perform a task reconciliation.
Explicit Recon:
Aurora will send a list of non-terminal task IDs and the master responds
with the latest state for each task, if possible.
Implicit Recon:
Aurora will send an empty list of tasks and the master responds with the latest
state for all currently known non-terminal tasks.
`,
}

var forceExplicitReconCmd = &cobra.Command{
	Use:   "explicit [batch_size]",
	Short: "Force the leading scheduler to perform an explicit recon.",
	Long: `Aurora will send a list of non-terminal task IDs and the master
responds with the latest state for each task, if possible.
`,
	Run:  explicitRecon,
	Args: cobra.MaximumNArgs(1),
}

func explicitRecon(cmd *cobra.Command, args []string) {
	var batchSize *int32

	log.Println("Forcing scheduler to perform an explicit reconciliation with Mesos")

	switch len(args) {
	case 0:
		log.Infoln("Using default batch size for explicit recon.")
	case 1:
		log.Infof("Using %v as batch size for explicit recon.\n", args[0])

	    // Get batch size from args and convert it to the right format
		batchInt, err := strconv.Atoi(args[0])
        if err != nil {
            log.Fatalf("error: %+v\n", err)
        }

		batchInt32 := int32(batchInt)
		batchSize = &batchInt32
    default:
        log.Fatalln("Provide 0 arguments to use default batch size or one argument to use a custom batch size.")
	}

	err := client.ForceExplicitTaskReconciliation(batchSize)
	if err != nil {
		log.Fatalf("error: %+v\n", err.Error())
	} else {
		fmt.Println("Explicit reconciliation started successfully")
	}
}

var forceImplicitReconCmd = &cobra.Command{
	Use:   "implicit",
	Short: "Force the leading scheduler to perform an implicit recon.",
	Long: `Forces leading scheduler to ask Mesos Master for a list of the latest state for
all currently known non-terminal tasks being run by Aurora.`,
	Run:  implicitRecon,
}

func implicitRecon(cmd *cobra.Command, args []string) {

	log.Println("Forcing scheduler to perform an implicit reconciliation with Mesos")
	err := client.ForceImplicitTaskReconciliation()
	if err != nil {
		log.Fatalf("error: %+v\n", err)
	} else {
		fmt.Println("Implicit reconciliation started successfully")
	}
}
