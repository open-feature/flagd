package ofrep

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/open-feature/flagd/core/pkg/evaluator"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/service/ofrep"
	"github.com/rs/xid"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const (
	key              = "key"
	singleEvaluation = "/ofrep/v1/evaluate/flags/{key}"
	bulkEvaluation   = "/ofrep/v1/evaluate/{path:flags\\/|flags}"
)

type handler struct {
	Logger        *logger.Logger
	evaluator     evaluator.IEvaluator
	contextValues map[string]any
	tracer        trace.Tracer
}

func NewOfrepHandler(logger *logger.Logger, evaluator evaluator.IEvaluator, contextValues map[string]any) http.Handler {
	h := handler{
		Logger:        logger,
		evaluator:     evaluator,
		contextValues: contextValues,
		tracer:        otel.Tracer("flagd.ofrep.v1"),
	}

	router := mux.NewRouter()
	router.HandleFunc(singleEvaluation, h.HandleFlagEvaluation).Methods("POST")
	router.HandleFunc(bulkEvaluation, h.HandleBulkEvaluation).Methods("POST")
	return otelhttp.NewHandler(router, "flagd.ofrep")
}

func (h *handler) HandleFlagEvaluation(w http.ResponseWriter, r *http.Request) {
	requestID := xid.New().String()
	defer h.Logger.ClearFields(requestID)

	// obtain flag key
	vars := mux.Vars(r)
	if vars == nil {
		h.writeJSONToResponse(
			http.StatusInternalServerError,
			ofrep.InternalError{ErrorDetails: "failed to obtain the flag key from the request"}, w)
		return
	}

	flagKey := vars[key]
	request, err := extractOfrepRequest(r)
	if err != nil {
		h.writeJSONToResponse(http.StatusBadRequest, ofrep.ContextErrorResponseFrom(flagKey), w)
		return
	}
	context := flagdContext(h.Logger, requestID, request, h.contextValues)
	evaluation := h.evaluator.ResolveAsAnyValue(r.Context(), requestID, flagKey, context)
	if evaluation.Error != nil {
		status, evaluationError := ofrep.EvaluationErrorResponseFrom(evaluation)
		h.writeJSONToResponse(status, evaluationError, w)
	} else {
		h.writeJSONToResponse(http.StatusOK, ofrep.SuccessResponseFrom(evaluation), w)
	}
}

func (h *handler) HandleBulkEvaluation(w http.ResponseWriter, r *http.Request) {
	requestID := xid.New().String()
	defer h.Logger.ClearFields(requestID)

	request, err := extractOfrepRequest(r)
	if err != nil {
		h.writeJSONToResponse(http.StatusBadRequest, ofrep.BulkEvaluationContextError(), w)
		return
	}

	context := flagdContext(h.Logger, requestID, request, h.contextValues)
	evaluations, metadata, err := h.evaluator.ResolveAllValues(r.Context(), requestID, context)
	if err != nil {
		h.Logger.WarnWithID(requestID, fmt.Sprintf("error from resolver: %v", err))

		res := ofrep.BulkEvaluationContextErrorFrom(model.GeneralErrorCode,
			fmt.Sprintf("Bulk evaluation failed. Tracking ID: %s", requestID))
		h.writeJSONToResponse(http.StatusInternalServerError, res, w)
	} else {
		h.writeJSONToResponse(http.StatusOK, ofrep.BulkEvaluationResponseFrom(evaluations, metadata), w)
	}
}

func (h *handler) writeJSONToResponse(status int, payload interface{}, w http.ResponseWriter) {
	// first marshal payload
	marshal, err := json.Marshal(payload)
	if err != nil {
		// always a 500
		h.Logger.Warn(fmt.Sprintf("error marshelling the response: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(marshal)
	if err != nil {
		h.Logger.Warn(fmt.Sprintf("error while writing response: %v", err))
	}
}

func extractOfrepRequest(req *http.Request) (ofrep.Request, error) {
	request := ofrep.Request{}
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil && err.Error() != "EOF" {
		return request, fmt.Errorf("decode error: %w", err)
	}

	return request, nil
}

func flagdContext(
	log *logger.Logger, requestID string, request ofrep.Request, staticContextValues map[string]any,
) map[string]any {
	context := make(map[string]any)
	if res, ok := request.Context.(map[string]any); ok {
		for k, v := range res {
			context[k] = v
		}
	} else {
		log.WarnWithID(requestID, "provided context does not comply with flagd, continuing ignoring the context")
	}

	for k, v := range staticContextValues {
		context[k] = v
	}

	return context
}
