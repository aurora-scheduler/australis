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
	"time"
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

type Job struct {
	Environment string            `yaml:"environment"`
	Role        string            `yaml:"role"`
	Name        string            `yaml:"name"`
	CPU         float64           `yaml:"cpu"`
	RAM         int64             `yaml:"ram"`
	Disk        int64             `yaml:"disk"`
	Executor    Executor          `yaml:"executor"`
	Instances   int32             `yaml:"instances"`
	MaxFailures int32             `yaml:"maxFailures"`
	URIs        []URI             `yaml:"uris"`
	Metadata    map[string]string `yaml:"labels"`
	Service     bool              `yaml:"service"`
	Thermos     []ThermosProcess  `yaml:",flow,omitempty"`
	Container   *Container        `yaml:"container,omitempty"`
}
type InstanceRange struct {
	First int32 `yaml:"first"`
	Last  int32 `yaml:"last"`
}

type VariableBatchStrategy struct {
	GroupSizes []int32 `yaml:"groupSizes"`
	AutoPause  bool    `yaml:"autoPause"`
}

type BatchStrategy struct {
	GroupSize int32 `yaml:"groupSize"`
	AutoPause bool  `yaml:"autoPause"`
}

type QueueStrategy struct {
	GroupSize int32 `yaml:"groupSize"`
}

type UpdateStrategy struct {
	VariableBatch *VariableBatchStrategy `yaml:"variableBatch"`
	Batch         *BatchStrategy         `yaml:"batch"`
	Queue         *QueueStrategy         `yaml:"queue"`
}
type UpdateSettings struct {
	MaxPerInstanceFailures int32           `yaml:"maxPerInstanceFailures"`
	MaxFailedInstances     int32           `yaml:"maxFailedInstances"`
	MinTimeInRunning       time.Duration   `yaml:"minTimeRunning"`
	RollbackOnFailure      bool            `yaml:"rollbackOnFailure"`
	InstanceRanges         []InstanceRange `yaml:"instanceRanges"`
	InstanceCount          int32           `yaml:"instanceCount"`
	PulseTimeout           time.Duration   `yaml:"pulseTimeout"`
	SLAAware               bool            `yaml:"slaAware"`
	Strategy               UpdateStrategy  `yaml:"strategy"`
}

type UpdateJob struct {
	JobConfig      Job            `yaml:"jobConfig"`
	UpdateSettings UpdateSettings `yaml:"updateSettings"`
}
