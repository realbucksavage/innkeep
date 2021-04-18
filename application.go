package innkeep

import (
	"fmt"
	"net/http"
)

type Application struct {
	Name           string
	HealthCheckURL string
	Metadata       MetadataMap
	Instances      []Instance
}

type Instance struct {
	Id              string
	Host            HostInfo
	Ports           map[string]Port
	Metadata        MetadataMap
	LastUpdatedTime int64
	LastDirtyTime   int64
	Status          string
}

type HostInfo struct {
	PreferHostname bool
	IPv4           string
	Hostname       string
}

type Port struct {
	Secure    int
	NonSecure int
}

func (h HostInfo) GetHost() string {
	if h.PreferHostname {
		return h.Hostname
	} else {
		return h.IPv4
	}
}

func (a Application) FindInstance(id string) (Instance, int, error) {

	for idx, i := range a.Instances {
		if i.Id == id {
			return i, idx, nil
		}
	}

	return Instance{}, -1, NewStatusError(fmt.Errorf("instance not found: %s", id), http.StatusNotFound)
}
