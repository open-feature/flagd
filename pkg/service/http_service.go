package service

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/open-feature/flagd/pkg/eval"
	gen "github.com/open-feature/flagd/schemas/protobuf/proto/go-server/schema/v1"
	log "github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type HTTPServiceConfiguration struct {
	Port           int32
	ServerCertPath string
	ServerKeyPath  string
}

type HTTPService struct {
	HTTPServiceConfiguration *HTTPServiceConfiguration
	GRPCService              *GRPCService
	Logger                   *log.Entry
}

func (s *HTTPService) ServeHTTPS() {

}

func (s *HTTPService) Serve(ctx context.Context, eval eval.IEvaluator) error {
	s.GRPCService.eval = eval

	// Mux Setup
	mux := runtime.NewServeMux(
		runtime.WithErrorHandler(s.HTTPErrorHandler),
	)
	// GRPC Setup
	grpcServer := grpc.NewServer()
	gen.RegisterServiceServer(grpcServer, s.GRPCService)
	err := gen.RegisterServiceHandlerFromEndpoint(
		context.Background(),
		mux,
		fmt.Sprintf("localhost:%d", s.HTTPServiceConfiguration.Port),
		// TODO: Add TLS here when we have a certificate
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		log.Fatal(err)
	}
	// Net listener
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.HTTPServiceConfiguration.Port))
	if err != nil {
		log.Fatal(err)
	}
	// Multiplexer listeners
	m := cmux.New(l)

	var server http.Server
	if s.HTTPServiceConfiguration.ServerCertPath != "" && s.HTTPServiceConfiguration.ServerKeyPath != "" {
		creds, err := tls.LoadX509KeyPair(s.HTTPServiceConfiguration.ServerCertPath, s.HTTPServiceConfiguration.ServerKeyPath)
		if err != nil {
			return err
		}
		server = http.Server{
			Handler:           mux,
			ReadHeaderTimeout: 60 * time.Second,
			TLSConfig: &tls.Config{
				Certificates: []tls.Certificate{creds},
				NextProtos:   []string{"h2"},
			},
		}

	} else {
		server = http.Server{
			Handler:           mux,
			ReadHeaderTimeout: 60 * time.Second,
		}
	}

	httpL := m.Match(cmux.HTTP1Fast())
	httpsL := m.Match(cmux.Any())
	grpcL := m.Match(cmux.HTTP2())

	go func() { handleServiceError(server.Serve(httpL)) }()     // HTTP
	go func() { handleServiceError(server.Serve(httpsL)) }()    // HTTP
	go func() { handleServiceError(grpcServer.Serve(grpcL)) }() // GRPC
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
		log.Error(err)
		log.Errorf("malformed error received by error handler, details received: %d - %v", len(details), details)
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
