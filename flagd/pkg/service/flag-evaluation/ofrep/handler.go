package ofrep

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/open-feature/flagd/core/pkg/evaluator"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/service/ofrep"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/open-feature/flagd/core/pkg/telemetry"
	"github.com/open-feature/flagd/flagd/pkg/service"
	evalservice "github.com/open-feature/flagd/flagd/pkg/service/flag-evaluation"
	"github.com/open-feature/flagd/flagd/pkg/service/flag-evaluation/ofrep/sse"
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

// configVersioner resolves when the flag configuration behind an SSE channel last changed, for the
// ADR-0008 flagConfigLastModified metadata. Implemented by the SSE change tracker; nil when SSE is
// disabled, and ok is false until a stream for the channel is live.
type configVersioner interface {
	Version(channel string) (etag string, lastModified int64, ok bool)
}

type handler struct {
	Logger                     *logger.Logger
	evaluator                  evaluator.IEvaluator
	contextValues              map[string]any
	headerToContextKeyMappings map[string]string
	metricsRecorder            telemetry.IMetricsRecorder
	tracer                     trace.Tracer

	versioner             configVersioner
	sseEnabled            bool
	sseInactivityDelaySec int
	ssePublicURL          string
}

// SSEConfig carries the SSE advertisement settings the bulk handler needs to expose the
// `eventStreams` block and conditional-evaluation ETags.
type SSEConfig struct {
	Enabled            bool
	Versioner          configVersioner
	InactivityDelaySec int
	PublicURL          string
}

func NewOfrepHandler(
	logger *logger.Logger,
	evaluator evaluator.IEvaluator,
	contextValues map[string]any,
	headerToContextKeyMappings map[string]string,
	metricsRecorder telemetry.IMetricsRecorder,
	serviceName string,
	sseCfg SSEConfig,
) http.Handler {
	h := handler{
		Logger:                     logger,
		evaluator:                  evaluator,
		contextValues:              contextValues,
		headerToContextKeyMappings: headerToContextKeyMappings,
		metricsRecorder:            metricsRecorder,
		tracer:                     otel.Tracer("flagd.ofrep.v1"),
		versioner:                  sseCfg.Versioner,
		sseEnabled:                 sseCfg.Enabled,
		sseInactivityDelaySec:      sseCfg.InactivityDelaySec,
		ssePublicURL:               sseCfg.PublicURL,
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

	// conditionalETag sits inside the metrics middleware so a 304 is recorded as a 304.
	router.Handle(bulkEvaluation,
		metricsmw.NewHTTPMetric(metricsmw.Config{
			Service:        serviceName,
			MetricRecorder: metricsRecorder,
			Logger:         logger,
			HandlerID:      bulkEvaluation,
		}).Handler(conditionalETag(http.HandlerFunc(h.HandleBulkEvaluation))),
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
	selectorExpression := r.Header.Get(service.FLAGD_SELECTOR_HEADER)
	selector, err := store.NewSelector(selectorExpression)
	if err != nil {
		h.writeJSONToResponse(http.StatusBadRequest, ofrep.EvaluationError{
			Key:          flagKey,
			ErrorCode:    model.GeneralErrorCode,
			ErrorDetails: fmt.Sprintf("invalid selector: %v", err),
		}, w)
		return
	}
	ctx := context.WithValue(r.Context(), store.SelectorContextKey{}, selector)

	evaluation := h.evaluator.ResolveAsAnyValue(ctx, requestID, flagKey, evaluationContext)
	if h.metricsRecorder != nil {
		h.metricsRecorder.RecordEvaluation(ctx, evaluation.Error, evaluation.Reason, evaluation.Variant, evaluation.FlagKey)
	}
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
	selectorExpression := r.Header.Get(service.FLAGD_SELECTOR_HEADER)
	selector, err := store.NewSelector(selectorExpression)
	if err != nil {
		res := ofrep.BulkEvaluationContextErrorFrom(model.GeneralErrorCode, fmt.Sprintf("invalid selector: %v", err))
		h.writeJSONToResponse(http.StatusBadRequest, res, w)
		return
	}
	ctx := context.WithValue(r.Context(), store.SelectorContextKey{}, selector)

	if trigger := r.URL.Query().Get(flagConfigEtagParam); trigger != "" {
		h.Logger.Debug(fmt.Sprintf("bulk refetch triggered by %s=%s", flagConfigEtagParam, trigger))
	}

	evaluations, metadata, err := h.evaluator.ResolveAllValues(ctx, requestID, evaluationContext)
	if h.metricsRecorder != nil {
		for _, evaluation := range evaluations {
			h.metricsRecorder.RecordEvaluation(ctx, evaluation.Error, evaluation.Reason, evaluation.Variant, evaluation.FlagKey)
		}
	}
	if err != nil {
		h.Logger.WarnWithID(requestID, fmt.Sprintf("error from resolver: %v", err))

		res := ofrep.BulkEvaluationContextErrorFrom(model.GeneralErrorCode,
			fmt.Sprintf("Bulk evaluation failed. Tracking ID: %s", requestID))
		h.writeJSONToResponse(http.StatusInternalServerError, res, w)
		return
	}

	h.writeJSONToResponse(http.StatusOK,
		h.bulkResponse(selectorExpression, evaluations, metadata, h.configLastModified(selectorExpression)), w)
}

// configLastModified reports when the flag configuration behind the selector last changed, or 0
// when nothing is tracking it. The tracker's config ETag is deliberately unused: it describes the
// configuration, not the response, so it cannot validate this endpoint's cache.
func (h *handler) configLastModified(channel string) int64 {
	if h.versioner == nil {
		return 0
	}
	_, lastModified, ok := h.versioner.Version(channel)
	if !ok {
		return 0
	}
	return lastModified
}

// bulkResponse assembles the OFREP bulk response, adding the ADR-0008 eventStreams block and
// lastModified metadata when SSE is enabled.
func (h *handler) bulkResponse(selectorExpression string, evaluations []evaluator.AnyValue, metadata model.Metadata, lastModified int64) ofrep.BulkEvaluationResponse {
	response := ofrep.BulkEvaluationResponseFrom(evaluations, metadata)
	if !h.sseEnabled {
		return response
	}
	response.EventStreams = h.eventStreams(selectorExpression)
	if lastModified > 0 {
		if response.Metadata == nil {
			response.Metadata = model.Metadata{}
		}
		response.Metadata["flagConfigLastModified"] = lastModified
	}
	return response
}

// eventStreams builds the ADR-0008 eventStreams advertisement pointing OFREP clients back at
// this flagd's SSE endpoint. The advertised channel is the request's own selector expression,
// carried as the final path segment, so the stream covers exactly the flags the client just
// evaluated.
//
// It uses the structured `endpoint` form and omits origin unless a public URL is configured, so
// the client resolves the requestUri against the OFREP base URL it is already talking to.
func (h *handler) eventStreams(selectorExpression string) []ofrep.EventStream {
	requestURI := sse.ChannelPath(ssePath, selectorExpression)

	return []ofrep.EventStream{{
		Type:               "sse",
		InactivityDelaySec: h.sseInactivityDelaySec,
		Endpoint: &ofrep.EventStreamEndpoint{
			Origin:     strings.TrimSuffix(h.ssePublicURL, "/"),
			RequestUri: requestURI,
		},
	}}
}

// flagConfigEtagParam is the ADR-0008 query parameter carrying the config version that an SSE
// refetch event advertised. It is trigger metadata only; the conditional response is decided by
// the If-None-Match header.
const flagConfigEtagParam = "flagConfigEtag"

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
