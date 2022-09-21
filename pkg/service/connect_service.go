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
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	schemaV1 "go.buf.build/open-feature/flagd-connect/open-feature/flagd/schema/v1"
	schemaConnectV1 "go.buf.build/open-feature/flagd-connect/open-feature/flagd/schema/v1/schemav1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/protobuf/types/known/structpb"
)

type ConnectService struct {
	Eval                        eval.IEvaluator
	ConnectServiceConfiguration *ConnectServiceConfiguration
	handler                     http.Handler
	tls                         bool
}

type ConnectServiceConfiguration struct {
	Port             int32
	ServerCertPath   string
	ServerKeyPath    string
	ServerSocketPath string
	CORS             bool
}

func (s *ConnectService) Serve(ctx context.Context, eval eval.IEvaluator) error {
	s.Eval = eval
	mux := http.NewServeMux()

	lis, err := s.buildListener()
	if err != nil {
		return err
	}
	path, handler := schemaConnectV1.NewServiceHandler(s)
	s.handler = handler
	mux.Handle(path, s.handler)
	s.buildHandler(mux)

	srv := http.Server{
		ReadTimeout:       2 * time.Second,
		WriteTimeout:      4 * time.Second,
		ReadHeaderTimeout: time.Second,
		Handler:           handler,
	}
	errChan := make(chan error, 1)
	go func() {
		if s.tls {
			if err := srv.ServeTLS(
				lis,
				s.ConnectServiceConfiguration.ServerCertPath,
				s.ConnectServiceConfiguration.ServerKeyPath,
			); err != nil && !errors.Is(err, http.ErrServerClosed) {
				errChan <- err
			}
		} else {
			if err := srv.Serve(
				lis,
			); err != nil && !errors.Is(err, http.ErrServerClosed) {
				errChan <- err
			}
		}
		close(errChan)
	}()
	<-ctx.Done()
	if err := lis.Close(); err != nil {
		return err
	}
	return <-errChan
}

func (s *ConnectService) buildHandler(mux *http.ServeMux) {
	if s.ConnectServiceConfiguration.ServerCertPath != "" && s.ConnectServiceConfiguration.ServerKeyPath != "" {
		s.tls = true
		if s.ConnectServiceConfiguration.CORS {
			s.handler = newCORS().Handler(mux)
		} else {
			s.handler = mux
		}
	} else {
		if s.ConnectServiceConfiguration.CORS {
			s.handler = h2c.NewHandler(
				newCORS().Handler(mux),
				&http2.Server{},
			)
		} else {
			h2c.NewHandler(
				newCORS().Handler(mux),
				&http2.Server{},
			)
		}
	}
}

func (s *ConnectService) buildListener() (net.Listener, error) {
	if s.ConnectServiceConfiguration.ServerSocketPath != "" {
		return net.Listen("unix", s.ConnectServiceConfiguration.ServerSocketPath)
	} else {
		address := net.JoinHostPort("localhost", fmt.Sprintf("%d", s.ConnectServiceConfiguration.Port))
		return net.Listen("tcp", address)
	}
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

func newCORS() *cors.Cors {
	return cors.New(cors.Options{
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowOriginFunc: func(origin string) bool {
			// Allow all origins, which effectively disables CORS.
			return true
		},
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
		return connect.NewError(connect.CodeNotFound, err)
	case model.TypeMismatchErrorCode:
		return connect.NewError(connect.CodeInvalidArgument, err)
	case model.DisabledReason:
		return connect.NewError(connect.CodeUnavailable, err)
	case model.ParseErrorCode:
		return connect.NewError(connect.CodeDataLoss, err)
	}
	return err
}
