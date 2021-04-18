package eureka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/realbucksavage/innkeep"
	"github.com/realbucksavage/innkeep/registry"
	"k8s.io/klog/v2"
)

func makeHandler(reg registry.Registry) http.Handler {
	m := mux.NewRouter().PathPrefix("/eureka/apps").Subrouter()

	m.Methods("GET").PathPrefix("").Handler(kithttp.NewServer(
		makeApplicationsEndpoint(reg),
		func(c context.Context, r *http.Request) (request interface{}, err error) { return nil, nil },
		kithttp.EncodeJSONResponse,
	))

	m.Methods("POST").PathPrefix("/{name}").Handler(kithttp.NewServer(
		makeRegisterApplicationsEndpoint(reg),
		decodeRegisterRequest,
		kithttp.EncodeJSONResponse,
	))

	return m
}

func makeApplicationsEndpoint(reg registry.Registry) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		apps, err := reg.Applications()
		if err != nil {
			return nil, err
		}

		resp := appsResponse{
			VerisonDelta: "1",
			Applications: make([]appInfo, 0),
		}

		for _, a := range apps {
			instances := make([]instanceInfo, 0)

			for _, i := range a.Instances {
				in := instanceInfo{
					InstanceID:       i.Id,
					Hostname:         i.Host.Hostname,
					App:              a.Name,
					IpAddr:           i.Host.GetHost(),
					Status:           i.Status,
					OverriddenStatus: i.Status,
					Metadata:         i.Metadata,
				}

				for kind, p := range i.Ports {
					switch kind {
					case "http":
						{
							in.Port = portInfo{Number: p.NonSecure, Enabled: "true"}
							in.SecurePort = portInfo{Number: p.Secure, Enabled: "true"}
						}
					default: // set ports set by other transports in metadata
						{
							i.Metadata["ports."+kind+".secure"] = fmt.Sprintf("%d", p.Secure)
							i.Metadata["ports."+kind+".nonSecure"] = fmt.Sprintf("%d", p.NonSecure)
						}
					}
				}

				instances = append(instances, in)
			}

			resp.Applications = append(resp.Applications, appInfo{
				Name:     a.Name,
				Instance: instances,
			})
		}

		type appResp struct {
			Apps appsResponse `json:"applications"`
		}
		return appResp{resp}, nil
	}
}

func decodeRegisterRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		klog.V(3).Error(err, "cannot read body of registration request")
		return nil, err
	}

	klog.V(4).Infof("attempting registration with body:\n%s", string(body))

	var req registerRequest
	if err := json.Unmarshal(body, &req); err != nil {
		klog.V(3).Error(err, "cannot unmarshal registration request")
		return nil, err
	}

	klog.V(2).Infof("parsed registration request: %v", req)

	if req.Instance.Metadata == nil {
		req.Instance.Metadata = make(innkeep.MetadataMap)
	}

	req.Instance.Metadata["__Innkeep_Orig_Body"] = string(body)
	return req, nil
}

func makeRegisterApplicationsEndpoint(reg registry.Registry) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(registerRequest)
		if !ok {
			klog.Errorf("request is not registerRequest struct")
			return nil, errors.New("bad input type")
		}

		inst := req.Instance
		app := &innkeep.Application{
			Name:           inst.App,
			HealthCheckURL: inst.HealthCheckURL,
			Instances: []innkeep.Instance{
				{
					Id: inst.InstanceID,
					Host: innkeep.HostInfo{
						PreferHostname: false,
						IPv4:           inst.Hostname,
						Hostname:       inst.Hostname,
					},
					Ports: map[string]innkeep.Port{
						"http": {
							Secure:    inst.Port.Number,
							NonSecure: inst.SecurePort.Number,
						},
					},
					Metadata:        inst.Metadata,
					Status:          inst.Status,
					LastUpdatedTime: time.Now().UnixNano(),
				},
			},
		}
		reg.Register(app)
		return statusCodeResponse{http.StatusNoContent}, nil
	}
}
