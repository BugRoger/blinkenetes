package blinkenpad

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/bugroger/kube-blinkenpad/pkg/mk2"
	"github.com/golang/glog"
	"github.com/rakyll/portmidi"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/cache"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
	"k8s.io/kubernetes/pkg/controller/framework"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/util/wait"
)

var (
	VERSION = "0.0.0-dev"
	Green   = [3]int{0, 63, 0}
	Yellow  = [3]int{63, 63, 0}
	Red     = [3]int{63, 0, 0}
	Black   = [3]int{0, 0, 0}
)

const (
	ResyncPeriod = 1 * time.Minute
)

type Options struct {
}

type Blinkenpad struct {
	client *client.Client
	pad    *mk2.Launchpad
	nodes  cache.StoreToNodeLister
	pods   cache.StoreToPodLister

	sync.RWMutex
}

func New(opts Options) *Blinkenpad {
	return &Blinkenpad{}
}

func (b *Blinkenpad) Start() {
	fmt.Printf("Welcome to Blinkenpad %v\n", VERSION)
	b.createClient()
	b.createPad()
	b.watchNodes()
	b.watchPods()
}

func (b *Blinkenpad) Stop() {
	b.pad.Reset()
}

func (b *Blinkenpad) createClient() {
	glog.V(2).Infof("Creating Client")
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	overrides := &clientcmd.ConfigOverrides{}

	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides).ClientConfig()
	b.handleError(err)

	client, err := client.New(config)
	b.handleError(err)

	b.client = client
	glog.V(3).Infof("  using %s", config.Host)
}

func (b *Blinkenpad) createPad() {
	b.handleError(portmidi.Initialize())

	pad, err := mk2.Open()
	b.handleError(err)
	pad.Reset()

	b.pad = pad
}

func (b *Blinkenpad) watchNodes() {
	lw := cache.NewListWatchFromClient(b.client, "nodes", api.NamespaceAll, fields.Everything())
	hf := framework.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) { b.refresh("Node Added") },
		UpdateFunc: func(oldObj, newObj interface{}) { b.refresh("Node Updated") },
		DeleteFunc: func(obj interface{}) { b.refresh("Node deleted") },
	}

	store, controller := framework.NewInformer(lw, &api.Node{}, ResyncPeriod, hf)

	go controller.Run(wait.NeverStop)

	b.nodes = cache.StoreToNodeLister{store}
}

func (b *Blinkenpad) watchPods() {
	lw := cache.NewListWatchFromClient(b.client, "pods", api.NamespaceAll, fields.Everything())
	hf := framework.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) { b.refresh("Pod Added") },
		UpdateFunc: func(oldObj, newObj interface{}) { b.refresh("Pod Updated") },
		DeleteFunc: func(obj interface{}) { b.refresh("Pod deleted") },
	}

	indexer, controller := framework.NewIndexerInformer(lw, &api.Pod{}, ResyncPeriod, hf, cache.Indexers{"node": PodNodeNameIndexFunc})
	b.pods.Indexer = indexer

	go controller.Run(wait.NeverStop)
}

func (b *Blinkenpad) refresh(message string) {
	b.Lock()
	defer b.Unlock()

	fmt.Println(message)
	fmt.Println(b.getMaxPodsOnAnyNode())
	for column := 0; column < 8; column++ {
		b.refreshColumn(column)
	}
}

func (b *Blinkenpad) refreshColumn(i int) {
	nodes, err := b.nodes.List()
	b.handleError(err)

	if len(nodes.Items) < i+1 {
		return
	}

	var nodeNames []string
	for _, node := range nodes.Items {
		nodeNames = append(nodeNames, node.Name)
	}
	sort.Strings(nodeNames)

	status := b.getNodeStatus(nodeNames[i])
	greenPods, yellowPods, redPods := b.getPodStatus(nodeNames[i])

	scale := b.getScale()
	greenButtons := (int)(math.Ceil(float64(greenPods) * scale))
	yellowButtons := (int)(math.Ceil(float64(yellowPods) * scale))
	redButtons := (int)(math.Ceil(float64(redPods) * scale))

	fmt.Printf("Column %v, %v: %v, green: %v/%v, yellow: %v/%v red: %v/%v\n", nodeNames[i], i, status, greenPods, greenButtons, yellowPods, yellowButtons, redPods, redButtons)

	b.pad.Light(i+1, 1, status[0], status[1], status[2])
	for j := 0; j < greenButtons; j++ {
		b.pad.Light(i+1, j+2, Green[0], Green[1], Green[2])
	}
	for j := 0; j < yellowButtons; j++ {
		b.pad.Light(i+1, j+2+greenButtons, Yellow[0], Yellow[1], Yellow[2])
	}
	for j := 0; j < redButtons; j++ {
		b.pad.Light(i+1, j+2+greenButtons+yellowButtons, Red[0], Red[1], Red[2])
	}
	for j := 2 + greenButtons + yellowButtons + redButtons; j < 9; j++ {
		b.pad.Light(i+1, j, Black[0], Black[1], Black[2])
	}
}

func (b *Blinkenpad) getScale() float64 {
	max := b.getMaxPodsOnAnyNode()
	if max <= 5 {
		return 1.0
	}

	if max <= 10 {
		return 0.5
	}

	if max <= 20 {
		return 0.25
	}

	if max <= 40 {
		return 0.125
	}

	return 0.06125
}

func (b *Blinkenpad) getMaxPodsOnAnyNode() int {
	nodes, err := b.nodes.List()
	b.handleError(err)

	max := 1
	for _, node := range nodes.Items {
		pods, err := b.pods.ByIndex("node", node.Name)
		b.handleError(err)

		if len(pods) > max {
			max = len(pods)
		}
	}

	return max
}

func (b *Blinkenpad) getNodeStatus(name string) [3]int {
	nodes, err := b.nodes.List()
	b.handleError(err)

	for _, node := range nodes.Items {
		if node.Name == name {
			if api.IsNodeReady(&node) {
				if node.Spec.Unschedulable {
					return Yellow
				}
				return Green
			}
		}
	}

	return Red
}

func (b *Blinkenpad) getPodStatus(node string) (green, yellow, red int) {
	pods, err := b.pods.ByIndex("node", node)
	b.handleError(err)

	for _, pod := range pods {
		switch b.fishForPodStatus(pod.(*api.Pod)) {
		case "Running":
			green++
		case "Pending":
			yellow++
		case "Terminating":
			yellow++
		case "ContainerCreating":
			yellow++
		default:
			red++
		}
	}

	return
}

func PodNodeNameIndexFunc(obj interface{}) ([]string, error) {
	if pod, ok := obj.(*api.Pod); ok {
		return []string{pod.Spec.NodeName}, nil
	} else {
		return []string{""}, fmt.Errorf("object is not a pod: %v", obj)
	}
}
