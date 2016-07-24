package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/bugroger/kube-blinkenpad/pkg/blinkenpad"
)

var opts blinkenpad.Options

func init() {
}

func main() {
	flag.Parse()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	b := blinkenpad.New(opts)

	go func() {
		b.Start()
	}()

	go func() {
		<-sigs
		b.Stop()
		done <- true
	}()

	<-done
}
