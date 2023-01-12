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
	"github.com/open-feature/flagd/pkg/eval"
	"github.com/open-feature/flagd/pkg/logger"
	"github.com/open-feature/flagd/pkg/model"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"github.com/rs/xid"
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
	Port             int32
	MetricsPort      int32
	ServerCertPath   string
	ServerKeyPath    string
	ServerSocketPath string
	CORS             []string
}

type eventingConfiguration struct {
	mu   *sync.RWMutex
	subs map[interface{}]chan Notification
}

func (s *ConnectService) Serve(ctx context.Context, eval eval.IEvaluator) error {
	s.Eval = eval
	s.eventingConfiguration = &eventingConfiguration{
		subs: make(map[interface{}]chan Notification),
		mu:   &sync.RWMutex{},
	}
	lis, err := s.setupServer()
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
	<-ctx.Done()
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	return <-errChan
}

func (s *ConnectService) setupServer() (net.Listener, error) {
	var lis net.Listener
	var err error
	mux := http.NewServeMux()
	if s.ConnectServiceConfiguration.ServerSocketPath != "" {
		lis, err = net.Listen("unix", s.ConnectServiceConfiguration.ServerSocketPath)
	} else {
		address := fmt.Sprintf(":%d", s.ConnectServiceConfiguration.Port)
		lis, err = net.Listen("tcp", address)
	}
	if err != nil {
		return nil, err
	}
	path, handler := schemaConnectV1.NewServiceHandler(s)
	mux.Handle(path, handler)
	mdlw := New(middlewareConfig{
		Recorder: NewRecorder(prometheusConfig{}),
	})
	h := Handler("", mdlw, mux)
	go func() {
		s.Logger.Info(fmt.Sprintf("metrics listening at %d", s.ConnectServiceConfiguration.MetricsPort))
		server := &http.Server{
			Addr:              fmt.Sprintf(":%d", s.ConnectServiceConfiguration.MetricsPort),
			ReadHeaderTimeout: 3 * time.Second,
		}
		server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/metrics" {
				promhttp.Handler().ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		})
		err := server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

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
	requestNotificationChan := make(chan Notification, 1)
	s.eventingConfiguration.mu.Lock()
	s.eventingConfiguration.subs[req] = requestNotificationChan
	s.eventingConfiguration.mu.Unlock()
	defer func() {
		s.eventingConfiguration.mu.Lock()
		delete(s.eventingConfiguration.subs, req)
		s.eventingConfiguration.mu.Unlock()
	}()
	requestNotificationChan <- Notification{
		Type: ProviderReady,
	}
	for {
		select {
		case <-time.After(20 * time.Second):
			err := stream.Send(&schemaV1.EventStreamResponse{
				Type: string(KeepAlive),
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

func (s *ConnectService) Notify(n Notification) {
	s.eventingConfiguration.mu.RLock()
	defer s.eventingConfiguration.mu.RUnlock()
	for _, send := range s.eventingConfiguration.subs {
		send <- n
	}
}

func (s *ConnectService) ResolveBoolean(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveBooleanRequest],
) (*connect.Response[schemaV1.ResolveBooleanResponse], error) {
	reqID := xid.New().String()
	defer s.Logger.ClearFields(reqID)
	s.Logger.WriteFields(
		reqID,
		zap.String("flag-key", req.Msg.GetFlagKey()),
		zap.Strings("context-keys", formatContextKeys(req.Msg.GetContext())),
	)
	s.Logger.WarnWithID(reqID, "test")
	s.Logger.DebugWithID(reqID, "boolean flag value requested")

	res := connect.NewResponse(&schemaV1.ResolveBooleanResponse{})
	result, variant, reason, err := s.Eval.ResolveBooleanValue(reqID, req.Msg.GetFlagKey(), req.Msg.GetContext())
	if err != nil {
		s.Logger.WarnWithID(reqID, fmt.Sprintf("returning error response, reason: %s", err.Error()))
		res.Msg.Reason = model.ErrorReason
		return res, errFormat(err)
	}

	s.Logger.DebugWithID(reqID, fmt.Sprintf("flag evaluation response: %t, %s, %s", result, variant, reason))
	res.Msg.Reason = reason
	res.Msg.Value = result
	res.Msg.Variant = variant
	return res, nil
}

func (s *ConnectService) ResolveString(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveStringRequest],
) (*connect.Response[schemaV1.ResolveStringResponse], error) {
	reqID := xid.New().String()
	defer s.Logger.ClearFields(reqID)
	s.Logger.WriteFields(
		reqID,
		zap.String("flag-key", req.Msg.GetFlagKey()),
		zap.Strings("context-keys", formatContextKeys(req.Msg.GetContext())),
	)
	s.Logger.DebugWithID(reqID, "string flag value requested")

	res := connect.NewResponse(&schemaV1.ResolveStringResponse{})
	result, variant, reason, err := s.Eval.ResolveStringValue(reqID, req.Msg.GetFlagKey(), req.Msg.GetContext())
	if err != nil {
		s.Logger.WarnWithID(reqID, fmt.Sprintf("returning error response, reason: %s", err.Error()))
		res.Msg.Reason = model.ErrorReason
		return res, errFormat(err)
	}

	s.Logger.DebugWithID(reqID, fmt.Sprintf("flag evaluation response: %s, %s, %s", result, variant, reason))
	res.Msg.Reason = reason
	res.Msg.Value = result
	res.Msg.Variant = variant
	return res, nil
}

func (s *ConnectService) ResolveInt(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveIntRequest],
) (*connect.Response[schemaV1.ResolveIntResponse], error) {
	reqID := xid.New().String()
	defer s.Logger.ClearFields(reqID)
	s.Logger.WriteFields(
		reqID,
		zap.String("flag-key", req.Msg.GetFlagKey()),
		zap.Strings("context-keys", formatContextKeys(req.Msg.GetContext())),
	)
	s.Logger.DebugWithID(reqID, "int flag value requested")

	res := connect.NewResponse(&schemaV1.ResolveIntResponse{})
	result, variant, reason, err := s.Eval.ResolveIntValue(reqID, req.Msg.GetFlagKey(), req.Msg.GetContext())
	if err != nil {
		s.Logger.WarnWithID(reqID, fmt.Sprintf("returning error response, reason: %s", err.Error()))
		res.Msg.Reason = model.ErrorReason
		return res, errFormat(err)
	}

	s.Logger.DebugWithID(reqID, fmt.Sprintf("flag evaluation response: %d, %s, %s", result, variant, reason))
	res.Msg.Reason = reason
	res.Msg.Value = result
	res.Msg.Variant = variant
	return res, nil
}

func (s *ConnectService) ResolveFloat(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveFloatRequest],
) (*connect.Response[schemaV1.ResolveFloatResponse], error) {
	reqID := xid.New().String()
	defer s.Logger.ClearFields(reqID)
	s.Logger.WriteFields(
		reqID,
		zap.String("flag-key", req.Msg.GetFlagKey()),
		zap.Strings("context-keys", formatContextKeys(req.Msg.GetContext())),
	)
	s.Logger.DebugWithID(reqID, "float flag value requested")

	res := connect.NewResponse(&schemaV1.ResolveFloatResponse{})
	result, variant, reason, err := s.Eval.ResolveFloatValue(reqID, req.Msg.GetFlagKey(), req.Msg.GetContext())
	if err != nil {
		s.Logger.WarnWithID(reqID, fmt.Sprintf("returning error response, reason: %s", err.Error()))
		res.Msg.Reason = model.ErrorReason
		return res, errFormat(err)
	}

	s.Logger.DebugWithID(reqID, fmt.Sprintf("flag evaluation complete: %64f, %s, %s", result, variant, reason))
	res.Msg.Reason = reason
	res.Msg.Value = result
	res.Msg.Variant = variant
	return res, nil
}

func (s *ConnectService) ResolveObject(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveObjectRequest],
) (*connect.Response[schemaV1.ResolveObjectResponse], error) {
	reqID := xid.New().String()
	defer s.Logger.ClearFields(reqID)
	s.Logger.WriteFields(
		reqID,
		zap.String("flag-key", req.Msg.GetFlagKey()),
		zap.Strings("context-keys", formatContextKeys(req.Msg.GetContext())),
	)
	s.Logger.DebugWithID(reqID, "object flag value requested")

	res := connect.NewResponse(&schemaV1.ResolveObjectResponse{})
	result, variant, reason, err := s.Eval.ResolveObjectValue(reqID, req.Msg.GetFlagKey(), req.Msg.GetContext())
	if err != nil {
		s.Logger.WarnWithID(reqID, fmt.Sprintf("returning error response, reason: %s", err.Error()))
		res.Msg.Reason = model.ErrorReason
		return res, errFormat(err)
	}
	val, err := structpb.NewStruct(result)
	if err != nil {
		s.Logger.ErrorWithID(reqID, fmt.Sprintf("struct response construction: %v", err))
		return res, err
	}

	s.Logger.DebugWithID(reqID, fmt.Sprintf("flag evaluation response: %v, %s, %s", result, variant, reason))
	res.Msg.Reason = reason
	res.Msg.Value = val
	res.Msg.Variant = variant
	return res, nil
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
