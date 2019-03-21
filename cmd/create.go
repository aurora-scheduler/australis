package cmd

import (
	"os"

	yaml "gopkg.in/yaml.v2"

	realis "github.com/paypal/gorealis/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an Aurora Job",
	Run:   createJob,
	Args:  cobra.ExactArgs(1),
}

type URI struct {
	URI     string `yaml:"uri"`
	Extract bool   `yaml:"extract"`
	Cache   bool   `yaml:"cache"`
}

type Executor struct {
	Name string `yaml:"name"`
	Data string `yaml:"data"`
}

type ThermosProcess struct {
	Name string `yaml:"name"`
	Cmd  string `yaml:"cmd"`
}

type Job struct {
	Environment string            `yaml:"environment"`
	Role        string            `yaml:"role"`
	Name        string            `yaml:"name"`
	CPU         float64           `yaml:"cpu"`
	RAM         int64             `yaml:"ram"`
	Disk        int64             `yaml:"disk"`
	Executor    Executor          `yaml:"executor"`
	Instances   int32             `yaml:"instances"`
	URIs        []URI             `yaml:"uris"`
	Metadata    map[string]string `yaml:"labels"`
	Service     bool              `yaml:"service"`
	Thermos     []ThermosProcess  `yaml:",flow"`
}

func (j *Job) Validate() bool {
	if j.Name == "" {
		return false
	}

	if j.Role == "" {
		return false
	}

	if j.Environment == "" {
		return false
	}

	if j.Instances <= 0 {
		return false
	}

	if j.CPU <= 0.0 {
		return false
	}

	if j.RAM <= 0 {
		return false
	}

	if j.Disk <= 0 {
		return false
	}

	return true
}

func unmarshalJob(filename string) (Job, error) {

	job := Job{}

	if jobsFile, err := os.Open(filename); err != nil {
		return job, errors.Wrap(err, "unable to read the job config file")
	} else {
		if err := yaml.NewDecoder(jobsFile).Decode(&job); err != nil {
			return job, errors.Wrap(err, "unable to parse job config file")
		}

		if !job.Validate() {
			return job, errors.New("invalid job config")
		}
	}

	return job, nil
}

func createJob(cmd *cobra.Command, args []string) {

	job, err := unmarshalJob(args[0])

	if err != nil {
		log.Fatalln(err)
	}

	auroraJob := realis.NewJob().
		Environment(job.Environment).
		Role(job.Role).
		Name(job.Name).
		CPU(job.CPU).
		RAM(job.RAM).
		Disk(job.Disk).
		IsService(job.Service).
		InstanceCount(job.Instances)

	// Adding URIs.
	for _, uri := range job.URIs {
		auroraJob.AddURIs(uri.Extract, uri.Cache, uri.URI)
	}

	// Adding Metadata.
	for key, value := range job.Metadata {
		auroraJob.AddLabel(key, value)
	}

	// If thermos jobs processes are provided, use them
	if len(job.Thermos) > 0 {
		thermosExec := realis.ThermosExecutor{}
		for _, process := range job.Thermos {
			thermosExec.AddProcess(realis.NewThermosProcess(process.Name, process.Cmd))
		}
		auroraJob.ThermosExecutor(thermosExec)
	} else {
		// Non-Thermos executor
		if job.Executor.Name == "" {
			log.Fatal("no executor provided")
		}

		auroraJob.ExecutorName(job.Executor.Name)
		auroraJob.ExecutorData(job.Executor.Data)
	}

	if err := client.CreateJob(auroraJob); err != nil {
		log.Fatal("unable to create Aurora job: ", err)
	} else {
		if ok, monitorErr := client.InstancesMonitor(auroraJob.JobKey(), auroraJob.GetInstanceCount(), 5, 50); !ok || monitorErr != nil {
			if err := client.KillJob(auroraJob.JobKey()); err != nil {
				log.Fatal(monitorErr, err)
			}
			log.Fatal(monitorErr)
		}
	}

}
