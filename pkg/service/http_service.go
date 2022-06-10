package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/open-feature/flagd/pkg/eval"
	gen "github.com/open-feature/flagd/pkg/generated"
)

type HttpServiceConfiguration struct {
	Port int32
}

type HttpService struct {
	HttpServiceConfiguration *HttpServiceConfiguration
}

var defaultReason = "DEFAULT"
var errorReason = "ERROR"

type Server struct {
	eval eval.IEvaluator
}

// implement the generated ServerInterface.
// TODO: might be able to simplify some of this with generics.
// TODO: add improved, more RESTful error handling
func (s Server) ResolveBoolean(w http.ResponseWriter, r *http.Request, flagKey gen.FlagKey, params gen.ResolveBooleanParams) {
	result, err := s.eval.ResolveBooleanValue(flagKey, params.DefaultValue)
	if (err != nil) {
		fmt.Println(err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(gen.ResolutionDetailsBoolean{
			Value: &params.DefaultValue,
			Reason: &errorReason,
		})
		return;
	}
	json.NewEncoder(w).Encode(gen.ResolutionDetailsBoolean{
		Value: &result,
		Reason: &defaultReason,
	})
}

func (s Server) ResolveString(w http.ResponseWriter, r *http.Request, flagKey gen.FlagKey, params gen.ResolveStringParams) {
	result, err := s.eval.ResolveStringValue(flagKey, params.DefaultValue)
	if (err != nil) {
		fmt.Println(err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(gen.ResolutionDetailsString{
			Value: &params.DefaultValue,
			Reason: &errorReason,
		})
		return;
	}
	json.NewEncoder(w).Encode(gen.ResolutionDetailsString{
		Value: &result,
		Reason: &defaultReason,
	})
}

func (s Server) ResolveNumber(w http.ResponseWriter, r *http.Request, flagKey gen.FlagKey, params gen.ResolveNumberParams) {
	result, err := s.eval.ResolveNumberValue(flagKey, params.DefaultValue)
	if (err != nil) {
		fmt.Println(err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(gen.ResolutionDetailsNumber{
			Value: &params.DefaultValue,
			Reason: &errorReason,
		})
		return;
	}
	json.NewEncoder(w).Encode(gen.ResolutionDetailsNumber{
		Value: &result,
		Reason: &defaultReason,
	})
}

func (s Server) ResolveObject(w http.ResponseWriter, r *http.Request, flagKey gen.FlagKey, params gen.ResolveObjectParams) {
	result, err := s.eval.ResolveObjectValue(flagKey, params.DefaultValue.AdditionalProperties)
	if (err != nil) {
		fmt.Println(err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(gen.ResolutionDetailsObject{
			Value: &gen.ResolutionDetailsObject_Value{
				AdditionalProperties: params.DefaultValue.AdditionalProperties,
			},
			Reason: &defaultReason,
		})
		return;
	}
	json.NewEncoder(w).Encode(gen.ResolutionDetailsObject{
		Value: &gen.ResolutionDetailsObject_Value{
			AdditionalProperties: result,
		},
		Reason: &defaultReason,
	})
}

func (h *HttpService) Serve(eval eval.IEvaluator, ctx context.Context) error {
	if h.HttpServiceConfiguration == nil {
		return errors.New("http service configuration has not been initialised")
	}
	http.Handle("/", gen.Handler(Server{ eval }))
	http.ListenAndServe(fmt.Sprintf(":%d", h.HttpServiceConfiguration.Port), nil)

	<- ctx.Done()
	return nil
}
