package ofrep

import (
	"context"
	"encoding/json"
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
	metricsmw "github.com/open-feature/flagd/flagd/pkg/service/middleware/metrics"
	"github.com/rs/xid"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// ISelectorVersionTracker defines the interface for selector version tracking.
// it enables ETag-based cache validation for bulk evaluations.
type ISelectorVersionTracker interface {
	// ETag returns the current ETag for a selector, or empty if not tracked
	ETag(selectorExpression string) string
	// Track starts tracking a selector and returns its initial ETag
	Track(selectorExpression string) string
}

const (
	key               = "key"
	singleEvaluation  = "/ofrep/v1/evaluate/flags/{key}"
	bulkEvaluation    = "/ofrep/v1/evaluate/{path:flags\\/|flags}"
	headerETag        = "ETag"
	headerIfNoneMatch = "If-None-Match"
	headerContentType = "Content-Type"
	contentTypeJSON   = "application/json"
)

type handler struct {
	Logger                     *logger.Logger
	evaluator                  evaluator.IEvaluator
	contextValues              map[string]any
	headerToContextKeyMappings map[string]string
	tracer                     trace.Tracer
	versionTracker             ISelectorVersionTracker
}

func NewOfrepHandler(
	logger *logger.Logger,
	evaluator evaluator.IEvaluator,
	contextValues map[string]any,
	headerToContextKeyMappings map[string]string,
	metricsRecorder telemetry.IMetricsRecorder,
	serviceName string,
	versionTracker ISelectorVersionTracker,
) http.Handler {
	h := handler{
		Logger:                     logger,
		evaluator:                  evaluator,
		contextValues:              contextValues,
		headerToContextKeyMappings: headerToContextKeyMappings,
		tracer:                     otel.Tracer("flagd.ofrep.v1"),
		versionTracker:             versionTracker,
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
		h.writeJSONToResponse(http.StatusBadRequest, ofrep.ContextErrorResponseFrom(flagKey), w)
		return
	}
	evaluationContext := flagdContext(h.Logger, requestID, request, h.contextValues, r.Header, h.headerToContextKeyMappings)
	selectorExpression := r.Header.Get(service.FLAGD_SELECTOR_HEADER)
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
		h.writeJSONToResponse(http.StatusBadRequest, ofrep.BulkEvaluationContextError(), w)
		return
	}

	evaluationContext := flagdContext(h.Logger, requestID, request, h.contextValues, r.Header, h.headerToContextKeyMappings)
	selectorExpression := r.Header.Get(service.FLAGD_SELECTOR_HEADER)

	// check if client's ETag matches current version - return 304 if unchanged
	ifNoneMatch := r.Header.Get(headerIfNoneMatch)
	if h.versionTracker != nil && ifNoneMatch != "" {
		currentETag := h.versionTracker.ETag(selectorExpression)
		if currentETag != "" && ifNoneMatch == currentETag {
h.Logger.Debug("ETag match, returning 304", zap.String("selector", selectorExpression))
			w.Header().Add(headerETag, currentETag)
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	selector := store.NewSelector(selectorExpression)
	ctx := context.WithValue(r.Context(), store.SelectorContextKey{}, selector)

	evaluations, metadata, err := h.evaluator.ResolveAllValues(ctx, requestID, evaluationContext)
	if err != nil {
		h.Logger.WarnWithID(requestID, fmt.Sprintf("error from resolver: %v", err))

		res := ofrep.BulkEvaluationContextErrorFrom(model.GeneralErrorCode,
			fmt.Sprintf("Bulk evaluation failed. Tracking ID: %s", requestID))
		h.writeJSONToResponse(http.StatusInternalServerError, res, w)
	} else {
		response := ofrep.BulkEvaluationResponseFrom(evaluations, metadata)
		h.writeBulkEvaluationResponse(w, r, selectorExpression, response)
	}
}

// writes the bulk evaluation response with ETag support
func (h *handler) writeBulkEvaluationResponse(w http.ResponseWriter, _ *http.Request, selectorExpression string, response ofrep.BulkEvaluationResponse) {
	// marshal the response
	body, err := json.Marshal(response)
	if err != nil {
		h.Logger.Warn("error marshalling response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// track this selector and get ETag for response
	var eTag string
	if h.versionTracker != nil {
		eTag = h.versionTracker.Track(selectorExpression)
	}

	// write response with ETag
	w.Header().Add(headerContentType, contentTypeJSON)
	if eTag != "" {
		w.Header().Add(headerETag, eTag)
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	if err != nil {
		h.Logger.Warn("error while writing response", zap.Error(err))
	}
}

func (h *handler) writeJSONToResponse(status int, payload interface{}, w http.ResponseWriter) {
	// first marshal payload
	marshal, err := json.Marshal(payload)
	if err != nil {
		// always a 500
		h.Logger.Warn("error marshalling the response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add(headerContentType, contentTypeJSON)
	w.WriteHeader(status)
	_, err = w.Write(marshal)
	if err != nil {
		h.Logger.Warn("error while writing response", zap.Error(err))
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

// flagdContext returns combined context values from headers, static context (from cli) and request context.
// highest priority > header-context-from-cli > static-context-from-cli > request-context > lowest priority
func flagdContext(
	log *logger.Logger, requestID string, request ofrep.Request, staticContextValues map[string]any, headers http.Header, headerToContextKeyMappings map[string]string,
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

	for header, contextKey := range headerToContextKeyMappings {
		if values, ok := headers[header]; ok {
			context[contextKey] = values[0]
		}
	}

	return context
}
