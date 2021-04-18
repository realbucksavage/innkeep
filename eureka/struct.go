package eureka

import "github.com/realbucksavage/innkeep"

type appsResponse struct {
	VerisonDelta string    `json:"versions__delta"`
	AppsHashcode string    `json:"apps__hashcode"`
	Applications []appInfo `json:"application"`
}

type appInfo struct {
	Name     string         `json:"name"`
	Instance []instanceInfo `json:"instance"`
}

type instanceInfo struct {
	InstanceID       string            `json:"instanceId"`
	Hostname         string            `json:"hostName"`
	App              string            `json:"app"`
	IpAddr           string            `json:"ipAddr"`
	Status           string            `json:"status"`
	OverriddenStatus string            `json:"overriddenstatus"`
	Port             portInfo          `json:"port"`
	SecurePort       portInfo          `json:"securePort"`
	Metadata         map[string]string `json:"metadata"`
}

type dergisterRequest struct {
	app        string
	instanceID string
}

type registerRequest struct {
	Instance instanceRequest `json:"instance"`
}

type instanceRequest struct {
	InstanceID     string              `json:"instanceId"`
	App            string              `json:"app"`
	IpAddr         string              `json:"ipAddr"`
	Hostname       string              `json:"hostname"`
	Port           portInfo            `json:"port"`
	SecurePort     portInfo            `json:"securePort"`
	Metadata       innkeep.MetadataMap `json:"metadata"`
	Status         string              `json:"status"`
	HomePageURL    string              `json:"homePageUrl"`
	HealthCheckURL string              `json:"healthCheckUrl"`
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
