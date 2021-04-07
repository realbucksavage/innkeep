package eureka

import "github.com/realbucksavage/innkeep"

type appsResponse struct {
	VerisonDelta string    `json:"versions__delta"`
	AppsHashcode string    `json:"apps__hashcode"`
	Applications []appInfo `json:"application"`
}

type appInfo struct {
	Name string `json:"name"`
}

type registerRequest struct {
	Instance instanceRequest `json:"instance"`
}

type instanceRequest struct {
	InstanceID string              `json:"instanceId"`
	App        string              `json:"app"`
	IpAddr     string              `json:"ipAddr"`
	Hostname   string              `json:"hostname"`
	Port       portInfo            `json:"port"`
	SecurePort portInfo            `json:"securePort"`
	Metadata   innkeep.MetadataMap `json:"metadata"`
}

type portInfo struct {
	Number  int    `json:"$"`
	Enabled string `json:"@enabled"`
}

type statusCodeResponse struct {
	status int
}

func (s statusCodeResponse) StatusCode() int {
	return s.status
}
