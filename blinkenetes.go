package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IanLewis/launchpad/mk2"
	"github.com/golang/glog"
	"github.com/rakyll/portmidi"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/release_1_2"
)

const (
	resyncPeriod = 1 * time.Minute
)

type blinkenetes struct {
	client *clientset.Clientset
	pad    *mk2.Launchpad
}

func main() {
	client, err := newKubeClient()
	if err != nil {
		glog.Fatalf("Failed to create kubernetes client: %v", err)
	}

	blink := blinkenetes{
		pad:    newPad(),
		client: client,
	}

	watchPods(&blink)
	watchNodes(&blink)

	blink.pad.Reset()
	blink.showPacMan()
	blink.catchSignals()
}

func newPad() *mk2.Launchpad {
	var err error
	defer func() {
		if err != nil {
			log.Fatalln("Error while initializing Launchpad", err)
		}
	}()

	if err = portmidi.Initialize(); err != nil {
		return nil
	}

	pad, err := mk2.Open()
	if err != nil {
		return nil
	}
	pad.Reset()

	return pad
}

func (blink *blinkenetes) catchSignals() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		done <- true
	}()

	<-done
	blink.pad.Reset()
}
