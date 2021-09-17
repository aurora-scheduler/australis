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
	"errors"
	"strings"

	realis "github.com/aurora-scheduler/gorealis/v2"
	"github.com/aurora-scheduler/gorealis/v2/gen-go/apache/aurora"
)

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

type DockerContainer struct {
	Name string `yaml:"name"`
	Tag  string `yaml:"tag"`
}

type Container struct {
	Docker *DockerContainer `yaml:"docker"`
}

type ValueConstraint struct {
	Name    string   `yaml:"name"`
	Values  []string `yaml:"values"`
	Negated bool     `yaml:"negated"`
}

type LimitConstraint struct {
	Name  string `yaml:"name"`
	Limit int32  `yaml:"limit"`
}

type Job struct {
	Environment         string            `yaml:"environment"`
	Role                string            `yaml:"role"`
	Name                string            `yaml:"name"`
	CPU                 float64           `yaml:"cpu"`
	RAM                 int64             `yaml:"ram"`
	Disk                int64             `yaml:"disk"`
	Executor            Executor          `yaml:"executor"`
	Instances           int32             `yaml:"instances"`
	MaxFailures         int32             `yaml:"maxFailures"`
	URIs                []URI             `yaml:"uris"`
	Metadata            map[string]string `yaml:"labels"`
	Service             bool              `yaml:"service"`
	Priority            int32             `yaml:"priority"`
	Thermos             []ThermosProcess  `yaml:",flow,omitempty"`
	Container           *Container        `yaml:"container,omitempty"`
	CronSchedule        *string           `yaml:"cronSchedule,omitempty"`
	CronCollisionPolicy *string           `yaml:"cronCollisionPolicy,omitempty"`
	ValueConstraints    []ValueConstraint `yaml:"valueConstraints,flow,omitempty"`
	LimitConstraints    []LimitConstraint `yaml:"limitConstraints,flow,omitempty"`
}

func (j *Job) ToRealis() (*realis.AuroraJob, error) {
	auroraJob := realis.NewJob().
		Environment(j.Environment).
		Role(j.Role).
		Name(j.Name).
		CPU(j.CPU).
		RAM(j.RAM).
		Disk(j.Disk).
		IsService(j.Service).
		Priority(j.Priority).
		InstanceCount(j.Instances).
		MaxFailure(j.MaxFailures)

	if j.CronSchedule != nil {
		auroraJob.CronSchedule(*j.CronSchedule)
	}

	if j.CronCollisionPolicy != nil {
		// Ignoring error because we have already checked for it in the validate function
		policy, _ := aurora.CronCollisionPolicyFromString(*j.CronCollisionPolicy)
		auroraJob.CronCollisionPolicy(policy)
	}

	// Adding URIs.
	for _, uri := range j.URIs {
		auroraJob.AddURIs(uri.Extract, uri.Cache, uri.URI)
	}

	// Adding Metadata.
	for key, value := range j.Metadata {
		auroraJob.AddLabel(key, value)
	}

	// If thermos jobs processes are provided, use them
	if len(j.Thermos) > 0 {
		thermosExec := realis.ThermosExecutor{}
		for _, process := range j.Thermos {
			thermosExec.AddProcess(realis.NewThermosProcess(process.Name, process.Cmd))
		}
		auroraJob.ThermosExecutor(thermosExec)
	} else if j.Executor.Name != "" {
		// Non-Thermos executor
		if j.Executor.Name == "" {
			return nil, errors.New("no executor name provided")
		}

		auroraJob.ExecutorName(j.Executor.Name)
		auroraJob.ExecutorData(j.Executor.Data)
	} else if j.Container != nil {
		if j.Container.Docker == nil {
			return nil, errors.New("no container specified")
		}

		if j.Container.Docker.Tag != "" && !strings.ContainsRune(j.Container.Docker.Name, ':') {
			j.Container.Docker.Name += ":" + j.Container.Docker.Tag
		}
		auroraJob.Container(realis.NewDockerContainer().Image(j.Container.Docker.Name))

	}

	// Setting Constraints
	for _, valConstraint := range j.ValueConstraints {
		auroraJob.AddValueConstraint(valConstraint.Name, valConstraint.Negated, valConstraint.Values...)
	}

	for _, limit := range j.LimitConstraints {
		auroraJob.AddLimitConstraint(limit.Name, limit.Limit)
	}

	return auroraJob, nil
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
		return errors.New("disk must be greater than 0")
	}

	if len(j.Thermos) == 0 && j.Executor.Name == "" && j.Container == nil {
		return errors.New("task does not contain a thermos definition, a custom executor name, or a container to launch")
	}
	return nil
}

func (j *Job) ValidateCron() error {
	if j.CronSchedule == nil {
		return errors.New("cron schedule must be set")
	}

	if j.CronCollisionPolicy != nil {
		if _, err := aurora.CronCollisionPolicyFromString(*j.CronCollisionPolicy); err != nil {
			return err
		}
	} else {
		killExisting := aurora.CronCollisionPolicy_KILL_EXISTING.String()
		j.CronCollisionPolicy = &killExisting
	}

	return nil
}
