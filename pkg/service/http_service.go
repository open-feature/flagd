package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/open-feature/flagd/pkg/eval"
	gen "github.com/open-feature/flagd/pkg/generated"
	log "github.com/sirupsen/logrus"
)

type HTTPServiceConfiguration struct {
	Port int32
}

type HTTPService struct {
	HTTPServiceConfiguration *HTTPServiceConfiguration
}

type Server struct {
	eval eval.IEvaluator
}

// implement the generated ServerInterface.
// TODO: might be able to simplify some of this with generics.
// TODO: add improved, more RESTful error handling, we should inspect
// the returned ErrorCode and respond with a matching HTTP status code.
func (s Server) ResolveBoolean(
	w http.ResponseWriter, r *http.Request, flagKey gen.FlagKey, params gen.ResolveBooleanParams,
) {
	result, reason, err := s.eval.ResolveBooleanValue(flagKey, params.DefaultValue)
	if err != nil {
		message := err.Error()
		log.Error(message)
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(gen.ResolutionDetailsWithError{
			ErrorCode: &message,
			Reason:    &reason,
		})
		return
	}
	_ = json.NewEncoder(w).Encode(gen.ResolutionDetailsBoolean{
		Value:  result,
		Reason: &reason,
	})
}

func (s Server) ResolveString(
	w http.ResponseWriter, r *http.Request, flagKey gen.FlagKey, params gen.ResolveStringParams,
) {
	result, reason, err := s.eval.ResolveStringValue(flagKey, params.DefaultValue)
	if err != nil {
		message := err.Error()
		log.Error(message)
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(gen.ResolutionDetailsWithError{
			ErrorCode: &message,
			Reason:    &reason,
		})
		return
	}
	_ = json.NewEncoder(w).Encode(gen.ResolutionDetailsString{
		Value:  result,
		Reason: &reason,
	})
}

func (s Server) ResolveNumber(
	w http.ResponseWriter, r *http.Request, flagKey gen.FlagKey, params gen.ResolveNumberParams,
) {
	result, reason, err := s.eval.ResolveNumberValue(flagKey, params.DefaultValue)
	if err != nil {
		message := err.Error()
		log.Error(message)
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(gen.ResolutionDetailsWithError{
			ErrorCode: &message,
			Reason:    &reason,
		})
		return
	}
	_ = json.NewEncoder(w).Encode(gen.ResolutionDetailsNumber{
		Value:  result,
		Reason: &reason,
	})
}

func (s Server) ResolveObject(
	w http.ResponseWriter, r *http.Request, flagKey gen.FlagKey, params gen.ResolveObjectParams,
) {
	result, reason, err := s.eval.ResolveObjectValue(flagKey, params.DefaultValue.AdditionalProperties)
	if err != nil {
		message := err.Error()
		log.Error(message)
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(gen.ResolutionDetailsWithError{
			ErrorCode: &message,
			Reason:    &reason,
		})
		return
	}
	_ = json.NewEncoder(w).Encode(gen.ResolutionDetailsObject{
		Value: gen.ResolutionDetailsObject_Value{
			AdditionalProperties: result,
		},
		Reason: &reason,
	})
}

func (h *HTTPService) Serve(ctx context.Context, eval eval.IEvaluator) error {
	if h.HTTPServiceConfiguration == nil {
		return errors.New("http service configuration has not been initialised")
	}
	http.Handle("/", gen.Handler(Server{eval}))
	_ = http.ListenAndServe(fmt.Sprintf(":%d", h.HTTPServiceConfiguration.Port), nil)

	<-ctx.Done()
	return nil
}
