package ofrep

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/open-feature/flagd/core/pkg/evaluator"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/service/ofrep"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/open-feature/flagd/core/pkg/telemetry"
	"github.com/open-feature/flagd/flagd/pkg/service"
	evalservice "github.com/open-feature/flagd/flagd/pkg/service/flag-evaluation"
	metricsmw "github.com/open-feature/flagd/flagd/pkg/service/middleware/metrics"
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
	Logger                     *logger.Logger
	evaluator                  evaluator.IEvaluator
	contextValues              map[string]any
	headerToContextKeyMappings map[string]string
	tracer                     trace.Tracer
}

func NewOfrepHandler(
	logger *logger.Logger,
	evaluator evaluator.IEvaluator,
	contextValues map[string]any,
	headerToContextKeyMappings map[string]string,
	metricsRecorder telemetry.IMetricsRecorder,
	serviceName string,
) http.Handler {
	h := handler{
		Logger:                     logger,
		evaluator:                  evaluator,
		contextValues:              contextValues,
		headerToContextKeyMappings: headerToContextKeyMappings,
		tracer:                     otel.Tracer("flagd.ofrep.v1"),
	}

	router := mux.NewRouter()
	router.Handle(singleEvaluation,
		metricsmw.NewHTTPMetric(metricsmw.Config{
			Service:        serviceName,
			MetricRecorder: metricsRecorder,
			Logger:         logger,
			HandlerID:      singleEvaluation,
		}).Handler(http.HandlerFunc(h.HandleFlagEvaluation)),
	).Methods("POST")

	router.Handle(bulkEvaluation,
		metricsmw.NewHTTPMetric(metricsmw.Config{
			Service:        serviceName,
			MetricRecorder: metricsRecorder,
			Logger:         logger,
			HandlerID:      bulkEvaluation,
		}).Handler(http.HandlerFunc(h.HandleBulkEvaluation)),
	).Methods("POST")

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
		if h.handleExtractionError(w, err, ofrep.ContextErrorResponseFrom(flagKey)) {
			return
		}
	}
	evaluationContext := flagdContext(h.Logger, requestID, request, h.contextValues, r.Header, h.headerToContextKeyMappings)
	selectorExpression := service.SelectorExpressionFromHTTPHeaders(r.Header)
	selector := store.NewSelector(selectorExpression)
	ctx := context.WithValue(r.Context(), store.SelectorContextKey{}, selector)

	evaluation := h.evaluator.ResolveAsAnyValue(ctx, requestID, flagKey, evaluationContext)
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
		if h.handleExtractionError(w, err, ofrep.BulkEvaluationContextError()) {
			return
		}
	}

	evaluationContext := flagdContext(h.Logger, requestID, request, h.contextValues, r.Header, h.headerToContextKeyMappings)
	selectorExpression := service.SelectorExpressionFromHTTPHeaders(r.Header)
	ctx := r.Context()

	evaluations, metadata, err := evalservice.ResolveAllWithSelectorMerge(ctx, requestID, h.evaluator, evaluationContext, selectorExpression)
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

// handleExtractionError checks for errors from extractOfrepRequest and writes an appropriate response.
// It returns true if an error was handled.
func (h *handler) handleExtractionError(w http.ResponseWriter, err error, errorPayload any) bool {
	if err == nil {
		return false
	}
	var maxBytesErr *http.MaxBytesError
	if errors.As(err, &maxBytesErr) {
		h.writeJSONToResponse(http.StatusRequestEntityTooLarge,
			ofrep.InternalError{ErrorDetails: "request body too large"}, w)
		return true
	}
	h.writeJSONToResponse(http.StatusBadRequest, errorPayload, w)
	return true
}

func extractOfrepRequest(req *http.Request) (ofrep.Request, error) {
	request := ofrep.Request{}
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		// Propagate MaxBytesError so callers can return 413.
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			return request, err
		}
		if err.Error() != "EOF" {
			return request, fmt.Errorf("decode error: %w", err)
		}
	}

	return request, nil
}

func flagdContext(
	log *logger.Logger, requestID string, request ofrep.Request, staticContextValues map[string]any, headers http.Header, headerToContextKeyMappings map[string]string,
) map[string]any {
	context := make(map[string]any)
	if res, ok := request.Context.(map[string]any); ok {
		context = res
	} else {
		log.WarnWithID(requestID, "provided context does not comply with flagd, continuing ignoring the context")
	}

	return evalservice.MergeContextsAndHeaders(context, staticContextValues, headers, headerToContextKeyMappings)
}
