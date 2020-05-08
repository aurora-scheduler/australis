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
	"fmt"
	"time"

	realis "github.com/aurora-scheduler/gorealis/v2"
)

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
	MinTimeInRunning       time.Duration   `yaml:"minTimeInRunning"`
	RollbackOnFailure      bool            `yaml:"rollbackOnFailure"`
	InstanceRanges         []InstanceRange `yaml:"instanceRanges"`
	InstanceCount          int32           `yaml:"instanceCount"`
	PulseTimeout           time.Duration   `yaml:"pulseTimeout"`
	SLAAware               bool            `yaml:"slaAware"`
	Strategy               UpdateStrategy  `yaml:"strategy"`
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

type UpdateJob struct {
	JobConfig      Job            `yaml:"jobConfig"`
	UpdateSettings UpdateSettings `yaml:"updateSettings"`
}

func (u *UpdateJob) ToRealis() (*realis.JobUpdate, error) {
	jobConfig, err := u.JobConfig.ToRealis()
	if err != nil {
		return nil, fmt.Errorf("invalid job configuration %w", err)
	}

	update := realis.JobUpdateFromAuroraTask(jobConfig.AuroraTask())

	update.MaxPerInstanceFailures(u.UpdateSettings.MaxPerInstanceFailures).
		MaxFailedInstances(u.UpdateSettings.MaxFailedInstances).
		WatchTime(u.UpdateSettings.MinTimeInRunning).
		RollbackOnFail(u.UpdateSettings.RollbackOnFailure).
		PulseIntervalTimeout(u.UpdateSettings.PulseTimeout).
		SlaAware(u.UpdateSettings.SLAAware).
		InstanceCount(u.UpdateSettings.InstanceCount)

	strategy := u.UpdateSettings.Strategy
	switch {
	case strategy.VariableBatch != nil:
		update.VariableBatchStrategy(strategy.VariableBatch.AutoPause, strategy.VariableBatch.GroupSizes...)
	case strategy.Batch != nil:
		update.BatchUpdateStrategy(strategy.Batch.AutoPause, strategy.Batch.GroupSize)
	case strategy.Queue != nil:
		update.QueueUpdateStrategy(strategy.Queue.GroupSize)
	default:
		update.QueueUpdateStrategy(1)
	}

	for _, r := range u.UpdateSettings.InstanceRanges {
		update.AddInstanceRange(r.First, r.Last)
	}

	return update, nil

}
