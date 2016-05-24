package main

import (
	"fmt"

	"github.com/golang/glog"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/cache"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/release_1_2"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
	"k8s.io/kubernetes/pkg/controller/framework"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/util/wait"
)

func newKubeClient() (*clientset.Clientset, error) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	overrides := &clientcmd.ConfigOverrides{}

	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides).ClientConfig()
	if err != nil {
		return nil, err
	}

	glog.Infof("Using %s for kubernetes master", config.Host)
	glog.Infof("Using kubernetes API %v", config.GroupVersion)

	return clientset.NewForConfig(config)
}

func watchPods(blink *blinkenetes) cache.Store {
	lw := cache.NewListWatchFromClient(blink.client.CoreClient, "pods", api.NamespaceAll, fields.Everything())
	hf := framework.ResourceEventHandlerFuncs{
		AddFunc: blink.handlePodCreate,
		UpdateFunc: func(oldObj, newObj interface{}) {
			blink.handlePodUpdate(oldObj, newObj)
		},
		DeleteFunc: blink.handlePodDelete,
	}

	store, controller := framework.NewInformer(lw, &api.Pod{}, resyncPeriod, hf)

	go controller.Run(wait.NeverStop)
	return store
}

func watchNodes(blink *blinkenetes) cache.Store {
	lw := cache.NewListWatchFromClient(blink.client.CoreClient, "nodes", api.NamespaceAll, fields.Everything())
	hf := framework.ResourceEventHandlerFuncs{
		AddFunc: blink.handleNodeCreate,
		UpdateFunc: func(oldObj, newObj interface{}) {
			blink.handleNodeUpdate(oldObj, newObj)
		},
		DeleteFunc: blink.handleNodeDelete,
	}

	store, controller := framework.NewInformer(lw, &api.Node{}, resyncPeriod, hf)

	go controller.Run(wait.NeverStop)
	return store
}

func (blink *blinkenetes) handlePodCreate(obj interface{}) {
	if e, ok := obj.(*api.Pod); ok {
		fmt.Println("Pod created:", e.GetName(), e.Status.Phase)
	}
}

func (blink *blinkenetes) handleNodeCreate(obj interface{}) {
	if e, ok := obj.(*api.Node); ok {
		fmt.Println("Node created:", e.GetName(), e.Status.Phase)
	}
}

func (blink *blinkenetes) handlePodUpdate(old interface{}, new interface{}) {
	oldPod, okOld := old.(*api.Pod)
	newPod, okNew := new.(*api.Pod)

	if okOld && okNew {
		if oldPod.Status.PodIP != newPod.Status.PodIP {
			blink.handlePodDelete(oldPod)
			blink.handlePodCreate(newPod)
		}
	} else if okNew {
		blink.handlePodCreate(newPod)
	} else if okOld {
		blink.handlePodDelete(oldPod)
	}
}

func (blink *blinkenetes) handleNodeUpdate(old interface{}, new interface{}) {
	oldPod, okOld := old.(*api.Node)
	newPod, okNew := new.(*api.Node)

	if okOld && okNew {
		fmt.Println("Node updated:", newPod.GetName())
	} else if okNew {
		blink.handleNodeCreate(newPod)
	} else if okOld {
		blink.handleNodeDelete(oldPod)
	}
}

func (blink *blinkenetes) handlePodDelete(obj interface{}) {
	if e, ok := obj.(*api.Pod); ok {
		fmt.Println("Pod deleted:", e)
	}
}

func (blink *blinkenetes) handleNodeDelete(obj interface{}) {
	if e, ok := obj.(*api.Node); ok {
		fmt.Println("Node deleted:", e)
	}
}
