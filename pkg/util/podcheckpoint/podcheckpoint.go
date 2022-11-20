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
Copyright 2018 The Kubernetes Authors.

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

package pod

import (
	"context"
	"encoding/json"
	"fmt"
	"k8s.io/podcheckpoint/pkg/apis/podcheckpointcontroller/v1alpha1"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	podcheckpointclientset "k8s.io/podcheckpoint/pkg/client/clientset/versioned"
)

// PatchPodCheckpointStatus patches pod status.
func PatchPodCheckpointStatus(c podcheckpointclientset.Interface, namespace, name string, oldPodCheckpointStatus, newPodCheckpointStatus v1alpha1.PodCheckpointStatus) (*v1alpha1.PodCheckpoint, []byte, error) {
	patchBytes, err := preparePatchBytesforPodCheckpointStatus(namespace, name, oldPodCheckpointStatus, newPodCheckpointStatus)
	if err != nil {
		return nil, nil, err
	}

	updatedPodCheckpoint, err := c.PodcheckpointcontrollerV1alpha1().PodCheckpoints(namespace).Patch(context.TODO(), name, types.StrategicMergePatchType, patchBytes, v1.PatchOptions{}, "status")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to patch status %q for podcheckpoint %q/%q: %v", patchBytes, namespace, name, err)
	}
	return updatedPodCheckpoint, patchBytes, nil
}

func preparePatchBytesforPodCheckpointStatus(namespace, name string, oldPodCheckpointStatus, newPodCheckpointStatus v1alpha1.PodCheckpointStatus) ([]byte, error) {
	oldData, err := json.Marshal(v1alpha1.PodCheckpoint{
		Status: oldPodCheckpointStatus,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to Marshal oldData for pod %q/%q: %v", namespace, name, err)
	}

	newData, err := json.Marshal(v1alpha1.PodCheckpoint{
		Status: newPodCheckpointStatus,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to Marshal newData for pod %q/%q: %v", namespace, name, err)
	}

	patchBytes, err := strategicpatch.CreateTwoWayMergePatch(oldData, newData, v1alpha1.PodCheckpoint{})
	if err != nil {
		return nil, fmt.Errorf("failed to CreateTwoWayMergePatch for pod %q/%q: %v", namespace, name, err)
	}
	return patchBytes, nil
}
