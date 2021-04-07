package registry

import (
	"sync"

	"github.com/realbucksavage/innkeep"
	"k8s.io/klog/v2"
)

type Registry interface {
	Applications() ([]*innkeep.Application, error)
	Register(*innkeep.Application)
}

type defaultRegistry struct {
	mu     sync.RWMutex
	appMap map[string]*innkeep.Application
}

func (d *defaultRegistry) Applications() ([]*innkeep.Application, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	apps := make([]*innkeep.Application, 0)
	for _, a := range d.appMap {
		apps = append(apps, a)
	}

	return apps, nil
}

func (d *defaultRegistry) Register(app *innkeep.Application) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if a, found := d.appMap[app.Name]; found {
		klog.V(2).Infof("%s already exists... merging instances", app.Name)

		for _, inst := range app.Instances {
			var in innkeep.Instance
			idx := -1

			for i := 0; i < len(a.Instances); i++ {
				m := a.Instances[i]
				if inst.Id == m.Id {
					idx = i
					in = m
					break
				}
			}

			if idx == -1 {
				klog.V(3).Infof("instance %s is new", inst.Id)
				a.Instances = append(a.Instances, in)
			} else {
				klog.V(3).Infof("replacing instance %s with %v", inst.Id, in)
				a.Instances[idx] = in
			}
		}
	} else {
		d.appMap[app.Name] = app
	}
}

func NewRegistry() Registry {
	return &defaultRegistry{
		appMap: make(map[string]*innkeep.Application),
	}
}
