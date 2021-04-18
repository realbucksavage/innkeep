package registry

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/realbucksavage/innkeep"
	"k8s.io/klog/v2"
)

type Registry interface {
	VersionDelta() string
	Applications() ([]*innkeep.Application, error)
	Register(*innkeep.Application)
	DeregisterInstance(appName, instanceID string) error
}

type defaultRegistry struct {
	version int
	mu      sync.RWMutex
	appMap  map[string]*innkeep.Application
}

func (d *defaultRegistry) VersionDelta() string {
	return fmt.Sprintf("%d", d.version)
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
	defer d.increseVersionDelta()

	if a, found := d.appMap[app.Name]; found {
		klog.V(2).Infof("%s already exists... merging instances", app.Name)

		for _, inst := range app.Instances {
			in, idx, _ := a.FindInstance(inst.Id)

			if idx == -1 {
				klog.V(3).Infof("instance %s is new", inst.Id)
				a.Instances = append(a.Instances, inst)
			} else {
				klog.V(3).Infof("replacing instance %s with %v", in.Id, inst)
				a.Instances[idx] = inst
			}
		}
	} else {
		d.appMap[app.Name] = app
	}
}

func (d *defaultRegistry) DeregisterInstance(appName, instanceID string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	app, ok := d.appMap[appName]
	if !ok {
		klog.V(3).Infof("app not found: %s", appName)
		return innkeep.NewStatusError(fmt.Errorf("app not found: %s", appName), http.StatusNotFound)
	}

	if _, _, err := app.FindInstance(instanceID); err != nil {
		klog.V(3).Infof("instance not found: %s/%s", appName, instanceID)
		return err
	}

	instances := make([]innkeep.Instance, 0)
	for _, i := range app.Instances {
		if i.Id != instanceID {
			instances = append(instances, i)
		}
	}

	d.increseVersionDelta()

	if len(instances) == 0 {
		// dereg application
		delete(d.appMap, appName)
		return nil
	}

	app.Instances = instances
	return nil
}

func (d *defaultRegistry) increseVersionDelta() {
	d.version += 1
}

func NewRegistry() Registry {
	return &defaultRegistry{
		version: 1,
		appMap:  make(map[string]*innkeep.Application),
	}
}
