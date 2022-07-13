package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/open-feature/flagd/pkg/eval"
	"github.com/open-feature/flagd/pkg/model"
	gen "github.com/open-feature/flagd/schemas/protobuf/gen/v1"
	log "github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type HTTPServiceConfiguration struct {
	Port int32
}

type HTTPService struct {
	HTTPServiceConfiguration *HTTPServiceConfiguration
	eval                     eval.IEvaluator
	gen.UnimplementedServiceServer
}

func (s *HTTPService) Serve(ctx context.Context, eval eval.IEvaluator) error {
	s.eval = eval
	grpcServer := grpc.NewServer()
	gen.RegisterServiceServer(grpcServer, s)

	mux := runtime.NewServeMux(
		runtime.WithErrorHandler(s.HTTPErrorHandler),
	)
	err := gen.RegisterServiceHandlerFromEndpoint(
		context.Background(),
		mux,
		fmt.Sprintf("localhost:%d", s.HTTPServiceConfiguration.Port),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		log.Fatal(err)
	}

	server := http.Server{
		Handler: mux,
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.HTTPServiceConfiguration.Port))
	if err != nil {
		log.Fatal(err)
	}
	m := cmux.New(l)

	httpL := m.Match(cmux.HTTP1Fast())
	grpcL := m.Match(cmux.HTTP2())

	go func() { handleServiceError(server.Serve(httpL)) }()
	go func() { handleServiceError(grpcServer.Serve(grpcL)) }()
	go func() { handleServiceError(m.Serve()) }()

	<-ctx.Done()
	return nil
}

func (s HTTPService) HTTPErrorHandler(
	ctx context.Context,
	m *runtime.ServeMux,
	ma runtime.Marshaler,
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	if s, ok := status.FromError(err); ok {
		code := s.Code()
		switch {
		case codes.Unknown == code:
			w.WriteHeader(http.StatusInternalServerError)
		case codes.InvalidArgument == code:
			w.WriteHeader(http.StatusBadRequest)
		case codes.NotFound == code:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		var res []byte
		if res, err = json.Marshal(gen.ErrorResponse{
			ErrorCode: s.Message(),
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
			return
		}
		if _, err = w.Write(res); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
			return
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
	}
}

// TODO: might be able to simplify some of this with generics.
func (s HTTPService) ResolveBoolean(
	ctx context.Context,
	req *gen.ResolveBooleanRequest,
) (*gen.ResolveBooleanResponse, error) {
	res := gen.ResolveBooleanResponse{}
	result, variant, reason, err := s.eval.ResolveBooleanValue(req.GetFlagKey(), req.GetContext())
	if err != nil {
		return &res, handleEvaluationError(err, reason)
	}
	res.Reason = reason
	res.Value = result
	res.Variant = variant
	return &res, nil
}

func (s HTTPService) ResolveString(
	ctx context.Context,
	req *gen.ResolveStringRequest,
) (*gen.ResolveStringResponse, error) {
	res := gen.ResolveStringResponse{}
	result, variant, reason, err := s.eval.ResolveStringValue(req.GetFlagKey(), req.GetContext())
	if err != nil {
		return &res, handleEvaluationError(err, reason)
	}
	res.Reason = reason
	res.Value = result
	res.Variant = variant
	return &res, nil
}

func (s HTTPService) ResolveNumber(
	ctx context.Context,
	req *gen.ResolveNumberRequest,
) (*gen.ResolveNumberResponse, error) {
	res := gen.ResolveNumberResponse{}
	result, variant, reason, err := s.eval.ResolveNumberValue(req.GetFlagKey(), req.GetContext())
	if err != nil {
		return &res, handleEvaluationError(err, reason)
	}
	res.Reason = reason
	res.Value = result
	res.Variant = variant
	return &res, nil
}

func (s HTTPService) ResolveObject(
	ctx context.Context,
	req *gen.ResolveObjectRequest,
) (*gen.ResolveObjectResponse, error) {
	res := gen.ResolveObjectResponse{}
	result, variant, reason, err := s.eval.ResolveObjectValue(req.GetFlagKey(), req.GetContext())
	if err != nil {
		return &res, handleEvaluationError(err, reason)
	}
	val, err := structpb.NewStruct(result)
	if err != nil {
		return &res, handleEvaluationError(err, reason)
	}
	res.Reason = reason
	res.Value = val
	res.Variant = variant
	return &res, nil
}

// some basic mapping of errors from model to HTTP
func handleEvaluationError(err error, reason string) error {
	// TODO: we should consider creating a custom error that includes a code instead of using the message for this.
	statusCode := codes.Internal
	message := err.Error()
	switch message {
	case model.FlagNotFoundErrorCode:
		statusCode = codes.NotFound
	case model.TypeMismatchErrorCode:
		statusCode = codes.InvalidArgument
	}
	log.Error(message)
	return status.Error(statusCode, message)
}

// TODO: could be replaced with a logging client
func handleServiceError(err error) {
	log.Fatal(err)
}
