/*
Copyright 2015 The Kubernetes Authors.

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

package config

import (
	"flag"
	"fmt"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/podcheckpoint/pkg/apis/podcheckpointcontroller/v1alpha1"
	"os"
	"path/filepath"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"
	api "k8s.io/kubernetes/pkg/apis/core"
	kubetypes "k8s.io/kubernetes/pkg/kubelet/types"
	podcheckpoint_clientset "k8s.io/podcheckpoint/pkg/client/clientset/versioned"
	informers "k8s.io/podcheckpoint/pkg/client/informers/externalversions"
	"k8s.io/podcheckpoint/pkg/signals"
)

// WaitForAPIServerSyncPeriod is the period between checks for the node list/watch initial sync
const WaitForAPIServerSyncPeriod = 1 * time.Second

// NewSourceApiserver creates a config source that watches and pulls from the apiserver.
func NewSourceApiserver(c clientset.Interface, nodeName types.NodeName, nodeHasSynced func() bool, updates chan<- interface{}, configFile string) {
	lw := cache.NewListWatchFromClient(c.CoreV1().RESTClient(), "pods", metav1.NamespaceAll, fields.OneTermEqualSelector(api.PodHostField, string(nodeName)))

	// The Reflector responsible for watching pods at the apiserver should be run only after
	// the node sync with the apiserver has completed.
	klog.Info("Waiting for node sync before watching apiserver pods")
	go func() {
		for {
			if nodeHasSynced() {
				klog.V(4).Info("node sync completed")
				break
			}
			time.Sleep(WaitForAPIServerSyncPeriod)
			klog.V(4).Info("node sync has not completed yet")
		}
		klog.Info("Watching apiserver")
		newSourceApiserverFromLW(lw, updates)

		fmt.Println("Start PodCheckpointInformer!!!")
		newPodCheckpointInformer(c, updates, configFile)
	}()
}

// newSourceApiserverFromLW holds creates a config source that watches and pulls from the apiserver.
func newSourceApiserverFromLW(lw cache.ListerWatcher, updates chan<- interface{}) {
	send := func(objs []interface{}) {
		var pods []*v1.Pod
		for _, o := range objs {
			pods = append(pods, o.(*v1.Pod))
		}
		updates <- kubetypes.PodUpdate{Pods: pods, Op: kubetypes.SET, Source: kubetypes.ApiserverSource}
	}
	r := cache.NewReflector(lw, &v1.Pod{}, cache.NewUndeltaStore(send, cache.MetaNamespaceKeyFunc), 0)
	go r.Run(wait.NeverStop)
}

func newPodCheckpointInformer(c clientset.Interface, updates chan<- interface{}, configFile string) {
	stopCh := signals.SetupSignalHandler()
	send := func(obj interface{}) {
		fmt.Printf("exec podcheckpoint add func :%v ", obj)
		podcheckpoint := obj.(*v1alpha1.PodCheckpoint)
		if podcheckpoint.Status.Phase == "" {
			podcheckpoint.Status.Phase = v1alpha1.PodPrepareCheckpoint
		}
		updates <- kubetypes.PodUpdate{PodCheckpoint: podcheckpoint, Op: kubetypes.CHECKPOINT, Source: kubetypes.ApiserverSource}
	}

	var kubeconfig string
	if len(configFile) > 0 {
		kubeconfig = configFile
		fmt.Println("kubeletconfigFile : %+v", kubeconfig)
	} else {
		if home := os.Getenv("HOME"); home != "" {
			flag.StringVar(&kubeconfig, "kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) Absolute path to the kubeconfig file")
			fmt.Println("kubeletconfig : %+v", kubeconfig)
		} else {
			flag.StringVar(&kubeconfig, "kubeconfig", "", "Absolute path to the kubeconfig file")
		}
	}

	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Errorf("Error building kubeconfig: %s", err.Error())
	}
	podcheckpointClient, err := podcheckpoint_clientset.NewForConfig(cfg)
	if err != nil {
		fmt.Errorf("Error building example clientset: %s", err.Error())
	}
	podcheckpointInformerFactory := informers.NewSharedInformerFactory(podcheckpointClient, time.Second*30)
	podcheckpointInformer := podcheckpointInformerFactory.Podcheckpointcontroller().V1alpha1().PodCheckpoints()
	podcheckpointInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: send,
	})
	fmt.Println("Start podcheckpointInformerFactory")
	podcheckpointInformerFactory.Start(stopCh)
}
