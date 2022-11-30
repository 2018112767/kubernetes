/*
Copyright (c) 2014-2020 CGCL Labs
Container_Migrate is licensed under Mulan PSL v2.
You can use this software according to the terms and conditions of the Mulan PSL v2.
You may obtain a copy of Mulan PSL v2 at:
         http://license.coscl.org.cn/MulanPSL2
THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
See the Mulan PSL v2 for more details.
*/
/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PodCheckpointPhase is a label for the condition of a pod at the current time.
type PodCheckpointPhase string

// These are the valid statuses of podcheckpoints.
const (
	PodPrepareCheckpoint PodCheckpointPhase = "PodPrepareCheckpoint"
	PodCheckpointing     PodCheckpointPhase = "Checkpointing"
	PodSucceeded         PodCheckpointPhase = "Succeeded"
	PodHalfFailed        PodCheckpointPhase = "HalfFailed"
	PodFailed            PodCheckpointPhase = "Failed"
)

// ContainerCheckpointPhase is a label for the condition of a pod at the current time.
type ContainerCheckpointPhase string

// These are the valid statuses of podcheckpoints.
const (
	ContainerPrepareCheckpoint   ContainerCheckpointPhase = "ContainerPrepareCheckpoint"
	ContainerCheckpointing       ContainerCheckpointPhase = "ContainerCheckpointing"
	ContainerCheckpointSucceeded ContainerCheckpointPhase = "ContainerCheckpointSucceeded"
	ContainerCheckpointFailed    ContainerCheckpointPhase = "ContainerCheckpointFailed"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodCheckpoint is a specification for a PodCheckpoint resource
type PodCheckpoint struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PodCheckpointSpec   `json:"spec"`
	Status PodCheckpointStatus `json:"status"`
}

// PodCheckpointSpec is the spec for a PodCheckpoint resource
type PodCheckpointSpec struct {
	PodName string `json:"podName"`
	Storage string `json:"storage"`
	SecretName string `json:"secretName"`
}

// PodCheckpointStatus is the status for a PodCheckpoint resource
type PodCheckpointStatus struct {
	Phase               PodCheckpointPhase   `json:"phase"`
	ContainerConditions []ContainerCondition `json:"containerConditions"`
}

type ContainerCondition struct {
	ContainerName   string                   `json:"containerName"`
	ContainerID     string                   `json:"containerID"`
	CheckpointPhase ContainerCheckpointPhase `json:"checkpointPhase"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodCheckpointList is a list of PodCheckpoint resources
type PodCheckpointList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []PodCheckpoint `json:"items"`
}

