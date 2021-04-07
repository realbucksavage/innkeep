package eureka

import (
	"io/ioutil"
	"net/http"

	"github.com/realbucksavage/innkeep"
	"github.com/realbucksavage/innkeep/registry"
	"k8s.io/klog/v2"
)

type eurekaConnector struct {
	registry registry.Registry
	srv      *http.Server
}

func (e *eurekaConnector) Start() error {
	klog.V(2).Infof("starting eureka server at %s", e.srv.Addr)
	return e.srv.ListenAndServe()
}

func NewConnector(reg registry.Registry, bindAddr string) innkeep.Connector {
	klog.V(3).Infof("creating eureka transport")

	srv := &http.Server{
		Addr:    bindAddr,
		Handler: makeHandler(reg),
	}

	return &eurekaConnector{
		registry: reg,
		srv:      srv,
	}
}

type logAllHandler struct{}

func (l *logAllHandler) ServeHTTP(_ http.ResponseWriter, r *http.Request) {
	bytes, _ := ioutil.ReadAll(r.Body)
	klog.Infof("%s %s\n%v\n\n%%s", r.Method, r.URL, r.Header, string(bytes))
}
