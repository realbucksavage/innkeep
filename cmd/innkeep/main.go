package main

import "k8s.io/klog/v2"

func main() {
	klog.InitFlags(nil)
	opts := readOptions()

	if opts.Eureka.Enable {
		klog.V(3).Infof("creating eureka transport")
	}

	if opts.GRPC.Enable {
		klog.V(3).Infof("creating grpc transport")
	}

	if opts.AMQP.Enable {
		klog.V(3).Infof("creating AMQP transport")
	}
}
