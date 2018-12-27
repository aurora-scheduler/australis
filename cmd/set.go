package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

func init() {
	rootCmd.AddCommand(setCmd)

	// Set Sub-commands
	setCmd.AddCommand(setQuotaCmd)
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a value in the Aurora Scheduler.",
}

var setQuotaCmd = &cobra.Command{
	Use:   "quota <role> cpu:<value> ram:<value> disk:<value>",
	Short: "Set Quota resources for a role.",
	Long:  `Quotas can be set for roles in Aurora. Using this command we can set the resources reserved a role.`,
	Run:   setQuota,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 4 {
			return fmt.Errorf("role, cpu, ram, and disk resources must be provided")
		}

		*role = args[0]

		for i:=1; i < len(args); i++ {
			resourcePair := strings.Split(args[i], ":")

			if len(resourcePair) != 2 {
				return fmt.Errorf("all resources must be provided in <resource>:<value> format")
			}

			var err error
			switch resourcePair[0] {

			case "cpu":
				cpu, err = strconv.ParseFloat(resourcePair[1], 64)
				if err != nil {
					return errors.Wrap(err, "unable to convert CPU value provided to a floating point number")
				}
			case "ram":
				ram, err = strconv.ParseInt(resourcePair[1], 10, 64)
				if err != nil {
					return errors.Wrap(err, "unable to convert RAM value provided to a integer number")
				}

            case "disk":
				disk, err = strconv.ParseInt(resourcePair[1], 10, 64)
				if err != nil {
					return errors.Wrap(err, "unable to convert DISK value provided to a integer number")
				}
			default:
				return fmt.Errorf("unknown resource value provided, only cpu, ram, and disk are supported")
			}
		}

		return nil
	},
}

func setQuota(cmd *cobra.Command, args []string) {
	log.Println("Setting Quota resources for role.")
	log.Println(args)

	err := client.SetQuota(*role, &cpu, &ram, &disk)

	if err != nil {
		log.Fatal(err)
	}
}
