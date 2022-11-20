/*
Copyright 2022.

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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

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

// PodCheckpointSpec defines the desired state of PodCheckpoint
type PodCheckpointSpec struct {
	PodName    string `json:"podName"`
	Storage    string `json:"storage"`
	SecretName string `json:"secretName"`
}

// PodCheckpointStatus defines the observed state of PodCheckpoint
type PodCheckpointStatus struct {
	Phase               PodCheckpointPhase   `json:"phase"`
	ContainerConditions []ContainerCondition `json:"containerConditions"`
}

type ContainerCondition struct {
	ContainerName   string                   `json:"containerName"`
	ContainerID     string                   `json:"containerID"`
	CheckpointPhase ContainerCheckpointPhase `json:"checkpointPhase"`
}

// +genclient
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PodCheckpoint is the Schema for the podcheckpoints API
type PodCheckpoint struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PodCheckpointSpec   `json:"spec,omitempty"`
	Status PodCheckpointStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PodCheckpointList contains a list of PodCheckpoint
type PodCheckpointList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PodCheckpoint `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PodCheckpoint{}, &PodCheckpointList{})
}
