package innkeep

type Application struct {
	Name      string
	Instances []Instance
}

type Instance struct {
	Host     HostInfo
	Ports    map[string]Port
	Metadata MetadataMap
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
