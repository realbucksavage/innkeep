package eureka

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/realbucksavage/innkeep"
	"github.com/realbucksavage/innkeep/registry"
	"k8s.io/klog/v2"
)

func makeHandler(reg registry.Registry) http.Handler {
	m := mux.NewRouter().PathPrefix("/eureka/apps").Subrouter()

	m.Methods("GET").Path("/").Handler(kithttp.NewServer(
		makeApplicationsEndpoint(reg),
		func(c context.Context, r *http.Request) (request interface{}, err error) { return nil, nil },
		kithttp.EncodeJSONResponse,
	))

	m.Methods("POST").Path("/{name}").Handler(kithttp.NewServer(
		makeRegisterApplicationsEndpoint(reg),
		decodeRegisterRequest,
		kithttp.EncodeJSONResponse,
	))

	m.Methods("DELETE").Path("/{name}/{instanceID}").Handler(kithttp.NewServer(
		makeDeleteInstanceEndpoint(reg),
		decodeDeleteInstance,
		kithttp.EncodeJSONResponse,
	))

	return m
}

func decodeRegisterRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		klog.V(3).ErrorS(err, "cannot read body of registration request")
		return nil, err
	}

	klog.V(4).Infof("attempting registration with body:\n%s", string(body))

	var req registerRequest
	if err := json.Unmarshal(body, &req); err != nil {
		klog.V(3).ErrorS(err, "cannot unmarshal registration request")
		return nil, err
	}

	klog.V(2).Infof("parsed registration request: %v", req)

	if req.Instance.Metadata == nil {
		req.Instance.Metadata = make(innkeep.MetadataMap)
	}

	req.Instance.Metadata["__Innkeep_Orig_Body"] = string(body)
	return req, nil
}

func decodeDeleteInstance(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	return dergisterRequest{
		app:        vars["name"],
		instanceID: vars["instanceID"],
	}, nil
}
