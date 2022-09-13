package service

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/open-feature/flagd/pkg/eval"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	schemaV1 "go.buf.build/bufbuild/connect-go/james-milligan/flagd/schema/v1"
	schemaConnectV1 "go.buf.build/bufbuild/connect-go/james-milligan/flagd/schema/v1/schemav1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/protobuf/types/known/structpb"
)

type Service struct {
	Eval                 eval.IEvaluator
	ServiceConfiguration *ServiceConfiguration
	subs                 map[interface{}]chan int
	mx                   *sync.Mutex
}

type sub struct {
	send chan int
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
	s.mx = &sync.Mutex{}
	s.subs = map[interface{}]chan int{}
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

func (s *Service) Notify() {
	s.mx.Lock()
	for _, sub := range s.subs {
		sub <- 1
	}
	s.mx.Unlock()
}

func (s *Service) StreamBoolean(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveBooleanRequest],
	stream *connect.ServerStream[schemaV1.ResolveBooleanResponse],
) error {
	res := &schemaV1.ResolveBooleanResponse{}
	send := make(chan int, 1)
	s.mx.Lock()
	s.subs[req] = send
	s.mx.Unlock()
	defer func() {
		s.mx.Lock()
		delete(s.subs, req)
		s.mx.Unlock()
	}()

	send <- 1
	for {
		select {
		case <-send:

			value, variant, reason, _ := s.Eval.ResolveBooleanValue(req.Msg.GetFlagKey(), req.Msg.GetContext())
			if res.Value != value || res.Variant != variant || res.Reason != reason {
				res.Reason = reason
				res.Value = value
				res.Variant = variant
				err := stream.Send(res)
				if err != nil {
					log.Error(err)
				}
			}
		case <-ctx.Done():
			log.Info("client connection closed")
			return nil
		}
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

func (s *Service) StreamString(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveStringRequest],
	stream *connect.ServerStream[schemaV1.ResolveStringResponse],
) error {
	res := &schemaV1.ResolveStringResponse{}
	send := make(chan int, 1)
	s.mx.Lock()
	s.subs[req] = send
	s.mx.Unlock()
	defer func() {
		s.mx.Lock()
		delete(s.subs, req)
		s.mx.Unlock()
	}()

	send <- 1
	for {
		select {
		case <-send:
			value, variant, reason, err := s.Eval.ResolveStringValue(req.Msg.GetFlagKey(), req.Msg.GetContext())
			if err != nil {
				log.Error(err)
				return err
			}
			if res.Value != value || res.Variant != variant || res.Reason != reason {
				res.Reason = reason
				res.Value = value
				res.Variant = variant
				err := stream.Send(res)
				if err != nil {
					return err
				}
			}
		case <-ctx.Done():
			log.Info("client connection closed")
			return nil
		}
	}
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

func (s *Service) StreamInt(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveIntRequest],
	stream *connect.ServerStream[schemaV1.ResolveIntResponse],
) error {
	return nil
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

func (s *Service) StreamFloat(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveFloatRequest],
	stream *connect.ServerStream[schemaV1.ResolveFloatResponse],
) error {
	return nil
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

func (s *Service) StreamObject(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveObjectRequest],
	stream *connect.ServerStream[schemaV1.ResolveObjectResponse],
) error {
	return nil
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
