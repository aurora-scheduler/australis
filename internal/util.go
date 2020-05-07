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

package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aurora-scheduler/gorealis/v2/gen-go/apache/aurora"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

type MonitorCmdConfig struct {
	Cmd                             *cobra.Command
	MonitorInterval, MonitorTimeout time.Duration
	StatusList                      []string
}

var log *logrus.Logger

// Logger sets the logger available to the internal package
func Logger(l *logrus.Logger) {
	log = l
}

// ToJSON converts an interface to a JSON formatted string
func ToJSON(v interface{}) string {
	output, err := json.Marshal(v)

	if err != nil {
		log.Fatalf("Unable to serialize Aurora response: %+v", v)
	}

	return string(output)
}

func GetLoggingLevels() string {
	var buffer bytes.Buffer

	for _, level := range logrus.AllLevels {
		buffer.WriteString(level.String())
		buffer.WriteString(" ")
	}

	buffer.Truncate(buffer.Len() - 1)

	return buffer.String()
}

func MaintenanceMonitorPrint(hostResult map[string]bool, desiredStates []aurora.MaintenanceMode, toJson bool) {
	if len(hostResult) > 0 {
		// Create anonymous struct for JSON formatting
		output := struct {
			DesiredStates   []string `json:"desired_states"`
			Transitioned    []string `json:"transitioned"`
			NonTransitioned []string `json:"non-transitioned"`
		}{
			make([]string, 0),
			make([]string, 0),
			make([]string, 0),
		}

		for _, state := range desiredStates {
			output.DesiredStates = append(output.DesiredStates, state.String())
		}

		for host, ok := range hostResult {
			if ok {
				output.Transitioned = append(output.Transitioned, host)
			} else {
				output.NonTransitioned = append(output.NonTransitioned, host)
			}
		}

		if toJson {
			fmt.Println(ToJSON(output))
		} else {
			fmt.Printf("Entered %v status: %v\n", output.DesiredStates, output.Transitioned)
			fmt.Printf("Did not enter %v status: %v\n", output.DesiredStates, output.NonTransitioned)
		}
	}
}

func UnmarshalJob(filename string) (Job, error) {

	job := Job{}

	if jobsFile, err := os.Open(filename); err != nil {
		return job, errors.Wrap(err, "unable to read the job config file")
	} else {
		if err := yaml.NewDecoder(jobsFile).Decode(&job); err != nil {
			return job, errors.Wrap(err, "unable to parse job config file")
		}

		if err := job.Validate(); err != nil {
			return job, fmt.Errorf("invalid job config %w", err)
		}
	}

	return job, nil
}

func (j *Job) Validate() error {
	if j.Name == "" {
		return errors.New("job name not specified")
	}

	if j.Role == "" {
		return errors.New("job role not specified")
	}

	if j.Environment == "" {
		return errors.New("job environment not specified")
	}

	if j.Instances <= 0 {
		return errors.New("number of instances in job cannot be less than or equal to 0")
	}

	if j.CPU <= 0.0 {
		return errors.New("CPU must be greater than 0")
	}

	if j.RAM <= 0 {
		return errors.New("RAM must be greater than 0")
	}

	if j.Disk <= 0 {
		return errors.New("Disk must be greater than 0")
	}

	if len(j.Thermos) == 0 && j.Executor.Name == "" && j.Container == nil {
		return errors.New("task does not contain a thermos definition, a custom executor name, or a container to launch")
	}
	return nil
}

func UnmarshalUpdate(filename string) (UpdateJob, error) {

	updateJob := UpdateJob{}

	if jobsFile, err := os.Open(filename); err != nil {
		return updateJob, errors.Wrap(err, "unable to read the job config file")
	} else {
		if err := yaml.NewDecoder(jobsFile).Decode(&updateJob); err != nil {
			return updateJob, errors.Wrap(err, "unable to parse job config file")
		}

		if err := updateJob.JobConfig.Validate(); err != nil {
			return updateJob, fmt.Errorf("invalid job config %w", err)
		}
		if err := updateJob.UpdateSettings.Validate(); err != nil {
			return updateJob, fmt.Errorf("invalid update configuration %w", err)
		}
	}

	return updateJob, nil
}

func (u *UpdateSettings) Validate() error {
	if u.InstanceCount <= 0 {
		return errors.New("instance count must be larger than 0")
	}

	if u.Strategy.VariableBatch != nil {
		if len(u.Strategy.VariableBatch.GroupSizes) == 0 {
			return errors.New("variable batch strategy must specify at least one batch size")
		}
		for _, batch := range u.Strategy.VariableBatch.GroupSizes {
			if batch <= 0 {
				return errors.New("all groups in a variable batch strategy must be larger than 0")
			}
		}
	} else if u.Strategy.Batch != nil {
		if u.Strategy.Batch.GroupSize <= 0 {
			return errors.New("batch strategy must specify a group larger than 0")
		}
	} else if u.Strategy.Queue != nil {
		if u.Strategy.Queue.GroupSize <= 0 {
			return errors.New("queue strategy must specify a group larger than 0")
		}
	} else {
		log.Info("No strategy set, falling back on queue strategy with a group size 1")
	}
	return nil
}
