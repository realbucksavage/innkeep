package main

import (
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/realbucksavage/innkeep"
	"github.com/realbucksavage/innkeep/eureka"
	"github.com/realbucksavage/innkeep/registry"
	"k8s.io/klog/v2"
)

func main() {
	klog.InitFlags(nil)
	opts := readOptions()

	reg := registry.NewRegistry()

	var wg sync.WaitGroup
	shutdown := make(chan bool)

	if opts.Eureka.Enable {

		klog.V(3).Infof("creating eureka transport")

		conn := eureka.NewConnector(reg, opts.Eureka.BindAddr)
		wg.Add(1)

		go func() {
			if err := conn.Start(); err != nil && err != http.ErrServerClosed {
				klog.Fatalf("start eureka failed: %v", err)
			}
		}()
		go gracefulShutdown(conn, shutdown, wg.Done)
	}

	if opts.GRPC.Enable {
		klog.V(3).Infof("creating grpc transport")
	}

	if opts.AMQP.Enable {
		klog.V(3).Infof("creating AMQP transport")
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

	<-sig
	klog.V(2).Infof("interrupt received, beginning graceful shutdown")

	close(shutdown)
	wg.Wait()
}

func gracefulShutdown(conn innkeep.Connector, shutdown chan bool, doneFunc func()) {
	<-shutdown

	klog.V(3).Infof("attempting to gracefully shutdown connector %s", conn.Name())
	if err := conn.Stop(); err != nil {
		klog.Errorf("%s didn't shut down gracefully: %v", conn.Name(), err)
	}

	doneFunc()
}
