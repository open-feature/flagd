package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/open-feature/flagd/pkg/eval"
	gen "github.com/open-feature/flagd/schemas/protobuf/gen/v1"
	log "github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type HTTPServiceConfiguration struct {
	Port int32
}

type HTTPService struct {
	HTTPServiceConfiguration *HTTPServiceConfiguration
	GRPCService              *GRPCService
}

func (s *HTTPService) Serve(ctx context.Context, eval eval.IEvaluator) error {
	s.GRPCService.eval = eval
	grpcServer := grpc.NewServer()
	gen.RegisterServiceServer(grpcServer, s.GRPCService)

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
	st := status.Convert(err)
	switch {
	case st.Code() == codes.Unknown:
		w.WriteHeader(http.StatusInternalServerError)
	case st.Code() == codes.InvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	case st.Code() == codes.NotFound:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	details := st.Details()
	if len(details) != 1 {
		log.Errorf("malformed error recieved by error handler, details recieved: %d - %v", len(details), details)
		return
	}
	var res []byte
	if res, err = json.Marshal(details[0]); err != nil {
		log.Error(err)
		return
	}
	if _, err = w.Write(res); err != nil {
		log.Error(err)
		return
	}
}

// TODO: could be replaced with a logging client
func handleServiceError(err error) {
	log.Fatal(err)
}
