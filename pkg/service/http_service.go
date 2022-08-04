package service

import (
	"context"
	"crypto/rand"
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
	"google.golang.org/grpc/credentials"
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

func (s *HTTPService) tlsListener(l net.Listener) net.Listener {
	// Load certificates.
	certificate, err := tls.LoadX509KeyPair(s.HTTPServiceConfiguration.ServerCertPath,
		s.HTTPServiceConfiguration.ServerKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		Rand:         rand.Reader,
	}

	tlsl := tls.NewListener(l, config)
	return tlsl
}
func (s *HTTPService) Serve(ctx context.Context, eval eval.IEvaluator) error {
	s.GRPCService.eval = eval

	// Mux Setup
	mux := runtime.NewServeMux(
		runtime.WithErrorHandler(s.HTTPErrorHandler),
	)
	var tlsCreds credentials.TransportCredentials
	var dialOpts []grpc.DialOption
	var err error
	if s.HTTPServiceConfiguration.ServerCertPath != "" && s.HTTPServiceConfiguration.ServerKeyPath != "" {
		tlsCreds, err = loadTLSCredentials(s.HTTPServiceConfiguration.ServerCertPath,
			s.HTTPServiceConfiguration.ServerKeyPath)
		if err != nil {
			return err
		}
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(tlsCreds))
	} else {
		dialOpts = append(dialOpts, grpc.WithInsecure())
	}
	// GRPC Setup
	grpcServer := grpc.NewServer()
	gen.RegisterServiceServer(grpcServer, s.GRPCService)
	err = gen.RegisterServiceHandlerFromEndpoint(
		context.Background(),
		mux,
		fmt.Sprintf("localhost:%d", s.HTTPServiceConfiguration.Port),
		// TODO: Add TLS here when we have a certificate
		dialOpts,
	)
	if err != nil {
		log.Fatal(err)
	}
	// Net listener
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.HTTPServiceConfiguration.Port))
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 60 * time.Second,
	}

	tcpm := cmux.New(l)
	// We first match on HTTP 1.1 methods.
	httpl := tcpm.Match(cmux.HTTP1Fast())

	// If not matched, we assume that its TLS.
	tlsl := tcpm.Match(cmux.Any())
	tlsl = s.tlsListener(tlsl)

	// Now, we build another mux recursively to match HTTPS and GoRPC.
	// You can use the same trick for SSH.
	tlsm := cmux.New(tlsl)
	httpsl := tlsm.Match(cmux.HTTP1Fast())
	gorpcl := tlsm.Match(cmux.Any())

	go func() { handleServiceError(server.Serve(httpl)) }()      // HTTP
	go func() { handleServiceError(server.Serve(httpsl)) }()     // HTTP
	go func() { handleServiceError(grpcServer.Serve(gorpcl)) }() // GRPC
	go func() { handleServiceError(tlsm.Serve()) }()
	go func() { handleServiceError(tcpm.Serve()) }()
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
