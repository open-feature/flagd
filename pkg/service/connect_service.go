package service

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/open-feature/flagd/pkg/eval"
	"github.com/open-feature/flagd/pkg/model"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	schemaV1 "go.buf.build/open-feature/flagd-connect/open-feature/flagd/schema/v1"
	schemaConnectV1 "go.buf.build/open-feature/flagd-connect/open-feature/flagd/schema/v1/schemav1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/protobuf/types/known/structpb"
)

const ErrorPrefix = "FlagdError:"

type ConnectService struct {
	Eval                        eval.IEvaluator
	ConnectServiceConfiguration *ConnectServiceConfiguration
	tls                         bool
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

func (s *ConnectService) Serve(ctx context.Context, eval eval.IEvaluator) error {
	s.Eval = eval

	lis, err := s.setupServer()
	if err != nil {
		return err
	}

	errChan := make(chan error, 1)
	go func() {
		if s.tls {
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
		log.Printf("metrics listening at %d", s.ConnectServiceConfiguration.MetricsPort)
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
		ReadTimeout:       2 * time.Second,
		WriteTimeout:      4 * time.Second,
		ReadHeaderTimeout: time.Second,
		Handler:           handler,
	}
	return lis, nil
}

func (s *ConnectService) ResolveBoolean(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveBooleanRequest],
) (*connect.Response[schemaV1.ResolveBooleanResponse], error) {
	res := connect.NewResponse(&schemaV1.ResolveBooleanResponse{})
	result, variant, reason, err := s.Eval.ResolveBooleanValue(req.Msg.GetFlagKey(), req.Msg.GetContext())
	if err != nil {
		log.Error(err)
		res.Msg.Reason = model.ErrorReason
		return res, errFormat(err)
	}
	res.Msg.Reason = reason
	res.Msg.Value = result
	res.Msg.Variant = variant
	return res, nil
}

func (s *ConnectService) ResolveString(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveStringRequest],
) (*connect.Response[schemaV1.ResolveStringResponse], error) {
	res := connect.NewResponse(&schemaV1.ResolveStringResponse{})
	result, variant, reason, err := s.Eval.ResolveStringValue(req.Msg.GetFlagKey(), req.Msg.GetContext())
	if err != nil {
		log.Error(err)
		res.Msg.Reason = model.ErrorReason
		return res, errFormat(err)
	}
	res.Msg.Reason = reason
	res.Msg.Value = result
	res.Msg.Variant = variant
	return res, nil
}

func (s *ConnectService) ResolveInt(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveIntRequest],
) (*connect.Response[schemaV1.ResolveIntResponse], error) {
	res := connect.NewResponse(&schemaV1.ResolveIntResponse{})
	result, variant, reason, err := s.Eval.ResolveIntValue(req.Msg.GetFlagKey(), req.Msg.GetContext())
	if err != nil {
		log.Error(err)
		res.Msg.Reason = model.ErrorReason
		return res, errFormat(err)
	}
	res.Msg.Reason = reason
	res.Msg.Value = result
	res.Msg.Variant = variant
	return res, nil
}

func (s *ConnectService) ResolveFloat(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveFloatRequest],
) (*connect.Response[schemaV1.ResolveFloatResponse], error) {
	res := connect.NewResponse(&schemaV1.ResolveFloatResponse{})
	result, variant, reason, err := s.Eval.ResolveFloatValue(req.Msg.GetFlagKey(), req.Msg.GetContext())
	if err != nil {
		log.Error(err)
		res.Msg.Reason = model.ErrorReason
		return res, errFormat(err)
	}
	res.Msg.Reason = reason
	res.Msg.Value = result
	res.Msg.Variant = variant
	return res, nil
}

func (s *ConnectService) ResolveObject(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveObjectRequest],
) (*connect.Response[schemaV1.ResolveObjectResponse], error) {
	res := connect.NewResponse(&schemaV1.ResolveObjectResponse{})
	result, variant, reason, err := s.Eval.ResolveObjectValue(req.Msg.GetFlagKey(), req.Msg.GetContext())
	if err != nil {
		log.Error(err)
		res.Msg.Reason = model.ErrorReason
		return res, errFormat(err)
	}
	val, err := structpb.NewStruct(result)
	if err != nil {
		return res, err
	}
	res.Msg.Reason = reason
	res.Msg.Value = val
	res.Msg.Variant = variant
	return res, nil
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
