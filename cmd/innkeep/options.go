package main

import (
	"flag"
	"io/ioutil"

	"gopkg.in/yaml.v3"
	"k8s.io/klog/v2"
)

type options struct {
	Eureka   eurekaTransport `yaml:"eureka"`
	GRPC     grpcTransport   `yaml:"grpc"`
	AMQP     amqpTransport   `yaml:"amqp"`
	HTTPPort int             `yaml:"http_port"`
}

type grpcTransport struct {
	Enable   bool   `yaml:"enable"`
	BindAddr string `yaml:"addr"`
}

type eurekaTransport struct {
	Enable   bool   `yaml:"enable"`
	BindAddr string `yaml:"addr"`
	// TODO: peer aware eureka
}

type amqpTransport struct {
	Enable   bool   `yaml:"enable"`
	AMQPAddr string `yaml:"addr"`
	Exchange string `yaml:"exchange"`
}

func readOptions() options {

	var (
		configFile = flag.String("cfg", "config.yml", "Specify the configuration file to be used. Defaults to config.yml")
	)
	flag.Parse()

	bytes, err := ioutil.ReadFile(*configFile)
	if err != nil {
		klog.Fatalf("cannot open config file %s: %v", *configFile, err)
	}

	klog.V(4).Infof("parsing yaml file:\n%s", string(bytes))

	var opts options
	if err := yaml.Unmarshal(bytes, &opts); err != nil {
		klog.Fatalf("cannot parse yaml: %v", err)
	}

	return opts
}
