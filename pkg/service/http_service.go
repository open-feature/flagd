package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/open-feature/flagd/pkg/eval"
	gen "github.com/open-feature/flagd/pkg/generated"
	"github.com/open-feature/flagd/pkg/model"
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
func (s Server) ResolveBoolean(
	w http.ResponseWriter,
	r *http.Request,
	flagKey gen.FlagKey,
	params gen.ResolveBooleanParams,
) {
	var contextObj gen.Context
	err := json.NewDecoder(r.Body).Decode(&contextObj)
	if err != nil {
		handleParseError(err, w)
		return
	}

	result, reason, err := s.eval.ResolveBooleanValue(flagKey, params.DefaultValue, contextObj)
	if err != nil {
		handleEvaluationError(err, reason, w)
		return
	}
	_ = json.NewEncoder(w).Encode(gen.ResolutionDetailsBoolean{
		Value:  result,
		Reason: &reason,
	})
}

func (s Server) ResolveString(
	w http.ResponseWriter,
	r *http.Request,
	flagKey gen.FlagKey,
	params gen.ResolveStringParams,
) {
	var contextObj gen.Context
	err := json.NewDecoder(r.Body).Decode(&contextObj)
	if err != nil {
		handleParseError(err, w)
		return
	}

	result, reason, err := s.eval.ResolveStringValue(flagKey, params.DefaultValue, contextObj)
	if err != nil {
		handleEvaluationError(err, reason, w)
		return
	}
	_ = json.NewEncoder(w).Encode(gen.ResolutionDetailsString{
		Value:  result,
		Reason: &reason,
	})
}

func (s Server) ResolveNumber(
	w http.ResponseWriter,
	r *http.Request,
	flagKey gen.FlagKey,
	params gen.ResolveNumberParams,
) {
	var contextObj gen.Context
	err := json.NewDecoder(r.Body).Decode(&contextObj)
	if err != nil {
		handleParseError(err, w)
		return
	}

	result, reason, err := s.eval.ResolveNumberValue(flagKey, params.DefaultValue, contextObj)
	if err != nil {
		handleEvaluationError(err, reason, w)
		return
	}
	_ = json.NewEncoder(w).Encode(gen.ResolutionDetailsNumber{
		Value:  result,
		Reason: &reason,
	})
}

func (s Server) ResolveObject(
	w http.ResponseWriter,
	r *http.Request,
	flagKey gen.FlagKey,
	params gen.ResolveObjectParams,
) {
	var contextObj gen.Context
	err := json.NewDecoder(r.Body).Decode(&contextObj)
	if err != nil {
		handleParseError(err, w)
		return
	}

	result, reason, err := s.eval.ResolveObjectValue(flagKey, params.DefaultValue.AdditionalProperties, contextObj)
	if err != nil {
		handleEvaluationError(err, reason, w)
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

func handleParseError(err error, w http.ResponseWriter) {
	log.Errorf("Error parsing context from body: %s", err.Error())
	reason := model.ErrorReason
	errorCode := model.ParseErrorCode
	w.WriteHeader(400)
	if err = json.NewEncoder(w).Encode(gen.ResolutionDetailsWithError{
		ErrorCode: &errorCode,
		Reason:    &reason,
	}); err != nil {
		log.Errorf("Error encoding response: %s", err.Error())
	}
}

// some basic mapping of errors from model to HTTP
func handleEvaluationError(err error, reason string, w http.ResponseWriter) {
	// TODO: we should consider creating a custom error that includes a code instead of using the message for this.
	message := err.Error()
	switch message {
	case model.FlagNotFoundErrorCode:
		w.WriteHeader(404)
	case model.TypeMismatchErrorCode:
		w.WriteHeader(400)
	default:
		w.WriteHeader(500)
	}
	log.Error(message)
	if err = json.NewEncoder(w).Encode(gen.ResolutionDetailsWithError{
		ErrorCode: &message,
		Reason:    &reason,
	}); err != nil {
		log.Errorf("Error encoding response: %s", err.Error())
	}
}
