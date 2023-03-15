//nolint:dupl
package service

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	schemaConnectV1 "buf.build/gen/go/open-feature/flagd/bufbuild/connect-go/schema/v1/schemav1connect"
	schemaV1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/schema/v1"
	"github.com/bufbuild/connect-go"
	"github.com/open-feature/flagd/core/pkg/eval"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	iservice "github.com/open-feature/flagd/core/pkg/service"
	"github.com/open-feature/flagd/core/pkg/service/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"github.com/rs/xid"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/protobuf/types/known/structpb"
)

const ErrorPrefix = "FlagdError:"

type ConnectService struct {
	Logger                      *logger.Logger
	Eval                        eval.IEvaluator
	ConnectServiceConfiguration *ConnectServiceConfiguration
	eventingConfiguration       *eventingConfiguration
	server                      http.Server
}
type ConnectServiceConfiguration struct {
	ServerCertPath   string
	ServerKeyPath    string
	ServerSocketPath string
	CORS             []string
}

type eventingConfiguration struct {
	mu   *sync.RWMutex
	subs map[interface{}]chan iservice.Notification
}

func (s *ConnectService) Serve(ctx context.Context, eval eval.IEvaluator, svcConf iservice.Configuration) error {
	s.Eval = eval
	s.eventingConfiguration = &eventingConfiguration{
		subs: make(map[interface{}]chan iservice.Notification),
		mu:   &sync.RWMutex{},
	}
	lis, err := s.setupServer(svcConf)
	if err != nil {
		return err
	}

	errChan := make(chan error, 1)
	go func() {
		if s.ConnectServiceConfiguration.ServerCertPath != "" && s.ConnectServiceConfiguration.ServerKeyPath != "" {
			if err := s.server.ServeTLS(
				lis,
				s.ConnectServiceConfiguration.ServerCertPath,
				s.ConnectServiceConfiguration.ServerKeyPath,
			); err != nil && !errors.Is(err, http.ErrServerClosed) {
				errChan <- err
			}
		} else {
			if err := s.server.Serve(
				lis,
			); err != nil && !errors.Is(err, http.ErrServerClosed) {
				errChan <- err
			}
		}
		close(errChan)
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return s.server.Shutdown(ctx)
	}
}

func (s *ConnectService) setupServer(svcConf iservice.Configuration) (net.Listener, error) {
	var lis net.Listener
	var err error
	mux := http.NewServeMux()
	if s.ConnectServiceConfiguration.ServerSocketPath != "" {
		lis, err = net.Listen("unix", s.ConnectServiceConfiguration.ServerSocketPath)
	} else {
		address := fmt.Sprintf(":%d", svcConf.Port)
		lis, err = net.Listen("tcp", address)
	}
	if err != nil {
		return nil, err
	}
	path, handler := schemaConnectV1.NewServiceHandler(s)
	mux.Handle(path, handler)
	exporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}

	mdlw := metrics.New(metrics.MiddlewareConfig{
		Service:      "openfeature/flagd",
		MetricReader: exporter,
		Logger:       s.Logger,
	})
	h := metrics.Handler("", mdlw, mux)

	go bindMetrics(s, svcConf)

	if s.ConnectServiceConfiguration.ServerCertPath != "" && s.ConnectServiceConfiguration.ServerKeyPath != "" {
		handler = s.newCORS().Handler(h)
	} else {
		handler = h2c.NewHandler(
			s.newCORS().Handler(h),
			&http2.Server{},
		)
	}
	s.server = http.Server{
		ReadHeaderTimeout: time.Second,
		Handler:           handler,
	}
	return lis, nil
}

func (s *ConnectService) ResolveAll(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveAllRequest],
) (*connect.Response[schemaV1.ResolveAllResponse], error) {
	reqID := xid.New().String()
	defer s.Logger.ClearFields(reqID)
	res := &schemaV1.ResolveAllResponse{
		Flags: make(map[string]*schemaV1.AnyFlag),
	}
	values := s.Eval.ResolveAllValues(reqID, req.Msg.GetContext())
	for _, value := range values {
		switch v := value.Value.(type) {
		case bool:
			res.Flags[value.FlagKey] = &schemaV1.AnyFlag{
				Reason:  value.Reason,
				Variant: value.Variant,
				Value: &schemaV1.AnyFlag_BoolValue{
					BoolValue: v,
				},
			}
		case string:
			res.Flags[value.FlagKey] = &schemaV1.AnyFlag{
				Reason:  value.Reason,
				Variant: value.Variant,
				Value: &schemaV1.AnyFlag_StringValue{
					StringValue: v,
				},
			}
		case float64:
			res.Flags[value.FlagKey] = &schemaV1.AnyFlag{
				Reason:  value.Reason,
				Variant: value.Variant,
				Value: &schemaV1.AnyFlag_DoubleValue{
					DoubleValue: v,
				},
			}
		case map[string]any:
			val, err := structpb.NewStruct(v)
			if err != nil {
				s.Logger.ErrorWithID(reqID, fmt.Sprintf("struct response construction: %v", err))
				continue
			}
			res.Flags[value.FlagKey] = &schemaV1.AnyFlag{
				Reason:  value.Reason,
				Variant: value.Variant,
				Value: &schemaV1.AnyFlag_ObjectValue{
					ObjectValue: val,
				},
			}
		}
	}

	return connect.NewResponse(res), nil
}

func (s *ConnectService) EventStream(
	ctx context.Context,
	req *connect.Request[schemaV1.EventStreamRequest],
	stream *connect.ServerStream[schemaV1.EventStreamResponse],
) error {
	requestNotificationChan := make(chan iservice.Notification, 1)
	s.eventingConfiguration.mu.Lock()
	s.eventingConfiguration.subs[req] = requestNotificationChan
	s.eventingConfiguration.mu.Unlock()
	defer func() {
		s.eventingConfiguration.mu.Lock()
		delete(s.eventingConfiguration.subs, req)
		s.eventingConfiguration.mu.Unlock()
	}()
	requestNotificationChan <- iservice.Notification{
		Type: iservice.ProviderReady,
	}
	for {
		select {
		case <-time.After(20 * time.Second):
			err := stream.Send(&schemaV1.EventStreamResponse{
				Type: string(iservice.KeepAlive),
			})
			if err != nil {
				s.Logger.Error(err.Error())
			}
		case notification := <-requestNotificationChan:
			d, err := structpb.NewStruct(notification.Data)
			if err != nil {
				s.Logger.Error(err.Error())
			}
			err = stream.Send(&schemaV1.EventStreamResponse{
				Type: string(notification.Type),
				Data: d,
			})
			if err != nil {
				s.Logger.Error(err.Error())
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (s *ConnectService) Notify(n iservice.Notification) {
	s.eventingConfiguration.mu.RLock()
	defer s.eventingConfiguration.mu.RUnlock()
	for _, send := range s.eventingConfiguration.subs {
		send <- n
	}
}

func resolve[T constraints](
	logger *logger.Logger,
	resolver func(reqID, flagKey string, ctx *structpb.Struct) (T, string, string, error),
	flagKey string,
	ctx *structpb.Struct,
	resp response[T],
) error {
	reqID := xid.New().String()
	defer logger.ClearFields(reqID)

	logger.WriteFields(
		reqID,
		zap.String("flag-key", flagKey),
		zap.Strings("context-keys", formatContextKeys(ctx)),
	)

	result, variant, reason, evalErr := resolver(reqID, flagKey, ctx)
	if evalErr != nil {
		logger.WarnWithID(reqID, fmt.Sprintf("returning error response, reason: %v", evalErr))
		reason = model.ErrorReason
		evalErr = errFormat(evalErr)
	}

	if err := resp.SetResult(result, variant, reason); err != nil && evalErr == nil {
		logger.ErrorWithID(reqID, err.Error())
		return err
	}

	return evalErr
}

func (s *ConnectService) ResolveBoolean(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveBooleanRequest],
) (*connect.Response[schemaV1.ResolveBooleanResponse], error) {
	res := connect.NewResponse(&schemaV1.ResolveBooleanResponse{})
	err := resolve[bool](
		s.Logger, s.Eval.ResolveBooleanValue, req.Msg.GetFlagKey(), req.Msg.GetContext(), &booleanResponse{res},
	)

	return res, err
}

func (s *ConnectService) ResolveString(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveStringRequest],
) (*connect.Response[schemaV1.ResolveStringResponse], error) {
	res := connect.NewResponse(&schemaV1.ResolveStringResponse{})
	err := resolve[string](
		s.Logger, s.Eval.ResolveStringValue, req.Msg.GetFlagKey(), req.Msg.GetContext(), &stringResponse{res},
	)

	return res, err
}

func (s *ConnectService) ResolveInt(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveIntRequest],
) (*connect.Response[schemaV1.ResolveIntResponse], error) {
	res := connect.NewResponse(&schemaV1.ResolveIntResponse{})
	err := resolve[int64](
		s.Logger, s.Eval.ResolveIntValue, req.Msg.GetFlagKey(), req.Msg.GetContext(), &intResponse{res},
	)

	return res, err
}

func (s *ConnectService) ResolveFloat(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveFloatRequest],
) (*connect.Response[schemaV1.ResolveFloatResponse], error) {
	res := connect.NewResponse(&schemaV1.ResolveFloatResponse{})
	err := resolve[float64](
		s.Logger, s.Eval.ResolveFloatValue, req.Msg.GetFlagKey(), req.Msg.GetContext(), &floatResponse{res},
	)

	return res, err
}

func (s *ConnectService) ResolveObject(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveObjectRequest],
) (*connect.Response[schemaV1.ResolveObjectResponse], error) {
	res := connect.NewResponse(&schemaV1.ResolveObjectResponse{})
	err := resolve[map[string]any](
		s.Logger, s.Eval.ResolveObjectValue, req.Msg.GetFlagKey(), req.Msg.GetContext(), &objectResponse{res},
	)

	return res, err
}

func formatContextKeys(context *structpb.Struct) []string {
	res := []string{}
	for k := range context.AsMap() {
		res = append(res, k)
	}
	return res
}

func (s *ConnectService) newCORS() *cors.Cors {
	return cors.New(cors.Options{
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedOrigins: s.ConnectServiceConfiguration.CORS,
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{
			// Content-Type is in the default safelist.
			"Accept",
			"Accept-Encoding",
			"Accept-Post",
			"Connect-Accept-Encoding",
			"Connect-Content-Encoding",
			"Content-Encoding",
			"Grpc-Accept-Encoding",
			"Grpc-Encoding",
			"Grpc-Message",
			"Grpc-Status",
			"Grpc-Status-Details-Bin",
		},
	})
}

func bindMetrics(s *ConnectService, svcConf iservice.Configuration) {
	s.Logger.Info(fmt.Sprintf("metrics and probes listening at %d", svcConf.MetricsPort))
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", svcConf.MetricsPort),
		ReadHeaderTimeout: 3 * time.Second,
	}
	server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/healthz":
			w.WriteHeader(http.StatusOK)
		case "/readyz":
			if svcConf.ReadinessProbe() {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusPreconditionFailed)
			}
		case "/metrics":
			promhttp.Handler().ServeHTTP(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func errFormat(err error) error {
	switch err.Error() {
	case model.FlagNotFoundErrorCode:
		return connect.NewError(connect.CodeNotFound, fmt.Errorf("%s, %s", ErrorPrefix, err.Error()))
	case model.TypeMismatchErrorCode:
		return connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("%s, %s", ErrorPrefix, err.Error()))
	case model.DisabledReason:
		return connect.NewError(connect.CodeUnavailable, fmt.Errorf("%s, %s", ErrorPrefix, err.Error()))
	case model.ParseErrorCode:
		return connect.NewError(connect.CodeDataLoss, fmt.Errorf("%s, %s", ErrorPrefix, err.Error()))
	}

	return err
}
