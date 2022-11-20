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

package main

import (
	"flag"
	"fmt"
	"k8s.io/klog/v2"
	"time"

	"k8s.io/client-go/tools/clientcmd"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/cache"
	clientset "k8s.io/podcheckpoint/pkg/client/clientset/versioned"
	informers "k8s.io/podcheckpoint/pkg/client/informers/externalversions"
	"k8s.io/podcheckpoint/pkg/signals"
)

var (
	masterURL  string
	kubeconfig string
)

func main() {
	flag.Parse()

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	podcheckpointClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building example clientset: %s", err.Error())
	}
	podcheckpointInformerFactory := informers.NewSharedInformerFactory(podcheckpointClient, time.Second*30)
	// controller := NewController(podcheckpointClient,
	// 	podcheckpointInformerFactory.Podcheckpointcontroller().V1alpha1().PodCheckpoints())
	podcheckpointInformer := podcheckpointInformerFactory.Podcheckpointcontroller().V1alpha1().PodCheckpoints()
	// Set up an event handler for when PodCheckpoint resources change
	podcheckpointInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			klog.Infof("exec podcheckpoint add func :%v ", obj)
			fmt.Printf("exec podcheckpoint add func :%v ", obj)
		},
	})
	klog.Infof("Start podcheckpointInformerFactory")
	fmt.Println("Start podcheckpointInformerFactory")
	podcheckpointInformerFactory.Start(stopCh)
	<-stopCh
	// if err = controller.Run(2, stopCh); err != nil {
	// 	glog.Fatalf("Error running controller: %s", err.Error())
	// }
}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}
