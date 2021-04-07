package main

import (
	"net/http"

	"github.com/realbucksavage/innkeep/eureka"
	"github.com/realbucksavage/innkeep/registry"
	"k8s.io/klog/v2"
)

func main() {
	klog.InitFlags(nil)
	opts := readOptions()

	reg := registry.NewRegistry()

	if opts.Eureka.Enable {
		conn := eureka.NewConnector(reg, opts.Eureka.BindAddr)
		if err := conn.Start(); err != nil && err != http.ErrServerClosed {
			klog.Fatalf("start eureka failed: %v", err)
		}
	}

	if opts.GRPC.Enable {
		klog.V(3).Infof("creating grpc transport")
	}

	if opts.AMQP.Enable {
		klog.V(3).Infof("creating AMQP transport")
	}
}
