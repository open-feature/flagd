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
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	schemaV1 "go.buf.build/bufbuild/connect-go/james-milligan/flagd/schema/v1"
	schemav1 "go.buf.build/bufbuild/connect-go/james-milligan/flagd/schema/v1"
	schemaConnectV1 "go.buf.build/bufbuild/connect-go/james-milligan/flagd/schema/v1/schemav1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
)

type Service struct {
	Eval                 eval.IEvaluator
	ServiceConfiguration *ServiceConfiguration
	subs                 map[interface{}]chan NotificationType
}

type ServiceConfiguration struct {
	Port             int32
	ServerCertPath   string
	ServerKeyPath    string
	ServerSocketPath string
}

func (s *Service) Serve(ctx context.Context, eval eval.IEvaluator) error {
	var address string
	var handler http.Handler
	tls := false
	s.subs = make(map[interface{}]chan NotificationType)
	s.Eval = eval
	mux := http.NewServeMux()
	// sockets
	if s.ServiceConfiguration.ServerSocketPath != "" {
		address = s.ServiceConfiguration.ServerSocketPath
	} else {
		address = net.JoinHostPort("localhost", fmt.Sprintf("%d", s.ServiceConfiguration.Port))
	}
	// TLS
	path, handler := schemaConnectV1.NewServiceHandler(s)
	mux.Handle(path, handler)
	if s.ServiceConfiguration.ServerCertPath != "" && s.ServiceConfiguration.ServerKeyPath != "" {
		tls = true
		handler = newCORS().Handler(mux)
	} else {
		handler = h2c.NewHandler(
			newCORS().Handler(mux),
			&http2.Server{},
		)
	}

	srv := &http.Server{
		Addr:              address,
		Handler:           handler,
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
		MaxHeaderBytes:    8 * 1024, // 8KiB
	}

	errChan := make(chan error, 1)
	go func() {
		if tls {
			if err := srv.ListenAndServeTLS(s.ServiceConfiguration.ServerCertPath, s.ServiceConfiguration.ServerKeyPath); err != nil && !errors.Is(err, http.ErrServerClosed) {
				errChan <- err
			}
		} else {
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				errChan <- err
			}
		}
		close(errChan)
	}()

	<-ctx.Done()
	if err := srv.Shutdown(ctx); err != nil {
		return err
	}
	return <-errChan
}

func (s *Service) EventStream(
	ctx context.Context,
	req *connect.Request[emptypb.Empty],
	stream *connect.ServerStream[schemav1.EventStreamResponse],
) error {
	s.subs[req] = make(chan NotificationType, 1)
	defer func() {
		delete(s.subs, req)
	}()
	s.subs[req] <- PROVIDER_READY
	for {
		select {
		case notification := <-s.subs[req]:
			err := stream.Send(&schemav1.EventStreamResponse{
				Type: fmt.Sprintf("%s", notification),
			})
			if err != nil {
				log.Error(err)
			}
		case <-ctx.Done():
			log.Info("client connection closed")
			return nil
		}
	}
}

func (s *Service) ResolveBoolean(
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

func (s *Service) Notify(n NotificationType) {
	for _, send := range s.subs {
		send <- n
	}
}

func (s *Service) ResolveString(
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

func (s *Service) ResolveInt(
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

func (s *Service) ResolveFloat(
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

func (s *Service) ResolveObject(
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
	// To let web developers play with the demo service from browsers, we need a
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
