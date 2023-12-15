/*
Copyright 2019 The Volcano Authors.

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

package gpupredicates

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	extenderv1 "k8s.io/kube-scheduler/extender/v1"
	"volcano.sh/volcano/pkg/scheduler/api"
	"volcano.sh/volcano/pkg/scheduler/framework"
	"volcano.sh/volcano/pkg/scheduler/plugins/gpupredicates/predicate"
)

const (
	// PluginName indicates name of volcano scheduler plugin.
	PluginName = "gpupredicates"
)

type gpuPlugin struct {
	kubeClient kubernetes.Interface
}

// New function returns prioritizePlugin object
func New(aruguments framework.Arguments) framework.Plugin {
	return &gpuPlugin{}
}

func (gp *gpuPlugin) Name() string {
	return PluginName
}

func (gp *gpuPlugin) OnSessionOpen(ssn *framework.Session) {
	klog.V(5).Infof("Enter gpu plugin ...")
	gp.kubeClient = ssn.KubeClient()
	defer func() {
		klog.V(5).Infof("Leaving gpu plugin ...")
	}()
	ssn.AddPredicateFn(gp.Name(), func(task *api.TaskInfo, node *api.NodeInfo) ([]*api.Status, error) {
		response := []*api.Status{}

		gpuFilter, err := predicate.NewGPUFilter(gp.kubeClient)
		if err != nil {
			klog.Fatalf("Failed to new gpu quota filter: %s", err.Error())
			return append(response, &api.Status{Reason: err.Error(), Code: api.Error}), err
		}
		// Call gpu-admission predicate.Filter
		var admissionRequest extenderv1.ExtenderArgs
		admissionRequest.Nodes = &v1.NodeList{Items: []v1.Node{*(node.Node)}}
		admissionRequest.Pod = task.Pod
		admssionResponse := gpuFilter.Filter(admissionRequest)

		if admssionResponse.Error != "" {
			response = append(response, &api.Status{Reason: admssionResponse.Error, Code: api.Error})
		} else {
			response = append(response, &api.Status{Reason: "success", Code: api.Success})
		}
		klog.V(4).Infof("gpuPlugin result: %s,pod %s node %s", response[0].Reason, task.Pod.Name, node.Name)
		return response, nil
	})

}

func (bp *gpuPlugin) OnSessionClose(ssn *framework.Session) {
}
