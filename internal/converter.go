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
    "strings"

    realis "github.com/aurora-scheduler/gorealis/v2"
)

func (j *Job) ToRealis() (*realis.AuroraJob, error) {

    auroraJob := realis.NewJob().
        Environment(j.Environment).
        Role(j.Role).
        Name(j.Name).
        CPU(j.CPU).
        RAM(j.RAM).
        Disk(j.Disk).
        IsService(j.Service).
        InstanceCount(j.Instances).
        MaxFailure(j.MaxFailures)

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

    return auroraJob, nil
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
        SlaAware(u.UpdateSettings.SLAAware).InstanceCount(u.UpdateSettings.InstanceCount)

    strategy := u.UpdateSettings.Strategy
    switch {
    case strategy.VariableBatch != nil:
        update.VariableBatchStrategy(strategy.VariableBatch.AutoPause, strategy.VariableBatch.GroupSizes...)
    case strategy.Batch != nil:
        update.BatchUpdateStrategy(strategy.Batch.AutoPause,strategy.Batch.GroupSize)
    case strategy.Queue != nil:
        update.QueueUpdateStrategy(strategy.Queue.GroupSize)
    default:
        update.QueueUpdateStrategy(1)
    }

    for _,r := range u.UpdateSettings.InstanceRanges {
        update.AddInstanceRange(r.First, r.Last)
    }

    return update, nil


}
