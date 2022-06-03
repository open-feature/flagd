package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	gen "github.com/open-feature/flagd/pkg/generated"
)

type HttpServiceConfiguration struct {
	Port int32
}

type HttpService struct {
	HttpServiceConfiguration *HttpServiceConfiguration
}

var defaultReason = "DEFAULT"

// implement the generated ServerInterface.
type Server struct {}

func (s Server) ResolveBoolean(w http.ResponseWriter, r *http.Request, flagKey gen.FlagKey, params gen.ResolveBooleanParams) {
	json.NewEncoder(w).Encode(gen.ResolutionDetailsBoolean{
		Value: &params.DefaultValue,
		Reason: &defaultReason,
	})
}

func (s Server) ResolveString(w http.ResponseWriter, r *http.Request, flagKey gen.FlagKey, params gen.ResolveStringParams) {
	json.NewEncoder(w).Encode(gen.ResolutionDetailsString{
		Value: &params.DefaultValue,
		Reason: &defaultReason,
	})
}

func (s Server) ResolveNumber(w http.ResponseWriter, r *http.Request, flagKey gen.FlagKey, params gen.ResolveNumberParams) {
	json.NewEncoder(w).Encode(gen.ResolutionDetailsNumber{
		Value: &params.DefaultValue,
		Reason: &defaultReason,
	})
}

func (s Server) ResolveObject(w http.ResponseWriter, r *http.Request, flagKey gen.FlagKey, params gen.ResolveObjectParams) {
	json.NewEncoder(w).Encode(gen.ResolutionDetailsObject{
		Value: &gen.ResolutionDetailsObject_Value{
			AdditionalProperties: params.DefaultValue.AdditionalProperties,
		},
		Reason: &defaultReason,
	})
}

func (h *HttpService) Serve(handlerFunc func(IServiceRequest) IServiceResponse) error {
	if h.HttpServiceConfiguration == nil {
		return errors.New("http service configuration has not been initialised")
	}

	http.Handle("/", gen.Handler(Server{}))
	http.ListenAndServe(fmt.Sprintf(":%d", h.HttpServiceConfiguration.Port), nil)

	return nil
}
