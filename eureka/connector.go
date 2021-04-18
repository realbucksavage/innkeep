package eureka

import (
	"net/http"

	"github.com/realbucksavage/innkeep"
	"github.com/realbucksavage/innkeep/registry"
	"k8s.io/klog/v2"
)

type eurekaConnector struct {
	registry registry.Registry
	srv      *http.Server
}

func (e *eurekaConnector) Name() string {
	return "Innkeep/Eureka/v1"
}

func (e *eurekaConnector) Start() error {
	klog.V(2).Infof("starting eureka server at %s", e.srv.Addr)
	return e.srv.ListenAndServe()
}

func (e *eurekaConnector) Stop() error {
	return e.srv.Close()
}

func NewConnector(reg registry.Registry, bindAddr string) innkeep.Connector {
	srv := &http.Server{
		Addr:    bindAddr,
		Handler: makeHandler(reg),
	}

	return &eurekaConnector{
		registry: reg,
		srv:      srv,
	}
}
