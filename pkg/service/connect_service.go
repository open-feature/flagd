package service

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/open-feature/flagd/pkg/eval"
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
}

type ConnectServiceConfiguration struct {
	Port             int32
	ServerCertPath   string
	ServerKeyPath    string
	ServerSocketPath string
}

func (s *ConnectService) Serve(ctx context.Context, eval eval.IEvaluator) error {
	var handler http.Handler
	var lis net.Listener
	var err error

	tls := false

	s.Eval = eval
	mux := http.NewServeMux()
	// sockets

	if s.ConnectServiceConfiguration.ServerSocketPath != "" {
		lis, err = net.Listen("unix", s.ConnectServiceConfiguration.ServerSocketPath)
	} else {
		address := net.JoinHostPort("localhost", fmt.Sprintf("%d", s.ConnectServiceConfiguration.Port))
		lis, err = net.Listen("tcp", address)
	}
	if err != nil {
		return err
	}
	// TLS
	path, handler := schemaConnectV1.NewServiceHandler(s)
	mux.Handle(path, handler)

	if s.ConnectServiceConfiguration.ServerCertPath != "" && s.ConnectServiceConfiguration.ServerKeyPath != "" {
		tls = true
		handler = newCORS().Handler(mux)
	} else {
		handler = h2c.NewHandler(
			newCORS().Handler(mux),
			&http2.Server{},
		)
	}

	errChan := make(chan error, 1)
	go func() {
		if tls {
			if err := http.ServeTLS(
				lis,
				handler,
				s.ConnectServiceConfiguration.ServerCertPath,
				s.ConnectServiceConfiguration.ServerKeyPath,
			); err != nil && !errors.Is(err, http.ErrServerClosed) {
				errChan <- err
			}
		} else {
			if err := http.Serve(
				lis,
				handler,
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

func (s *ConnectService) ResolveBoolean(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveBooleanRequest],
) (*connect.Response[schemaV1.ResolveBooleanResponse], error) {
	res := connect.NewResponse(&schemaV1.ResolveBooleanResponse{})
	result, variant, reason, err := s.Eval.ResolveBooleanValue(req.Msg.GetFlagKey(), req.Msg.GetContext())
	if err != nil {
		log.Error(err)
		return res, err
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
		return res, err
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
		return res, err
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
		return res, err
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
		return res, err
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
	// To let web developers play with the demo ConnectService from browsers, we need a
	// very permissive CORS setup.
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
