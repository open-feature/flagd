package service

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/open-feature/flagd/pkg/eval"
	log "github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
	gen "go.buf.build/open-feature/flagd-server/open-feature/flagd/schema/v1"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type HTTPServiceConfiguration struct {
	Port             int32
	ServerCertPath   string
	ServerKeyPath    string
	ServerSocketPath string
}

type HTTPService struct {
	HTTPServiceConfiguration *HTTPServiceConfiguration
	GRPCService              *GRPCService
	Logger                   *log.Entry
}

func (s *HTTPService) tlsListener(l net.Listener) net.Listener {
	// Load TLS config
	config, err := loadTLSConfig(s.HTTPServiceConfiguration.ServerCertPath,
		s.HTTPServiceConfiguration.ServerKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	tlsl := tls.NewListener(l, config)
	return tlsl
}

func (s *HTTPService) ServerGRPC(ctx context.Context, mux *runtime.ServeMux) *grpc.Server {
	var address string
	var dialOpts []grpc.DialOption
	var err error
	// handle cert
	if s.HTTPServiceConfiguration.ServerCertPath != "" && s.HTTPServiceConfiguration.ServerKeyPath != "" {
		tlsCreds, err := loadTLSCredentials(s.HTTPServiceConfiguration.ServerCertPath,
			s.HTTPServiceConfiguration.ServerKeyPath)
		if err != nil {
			log.Fatal(err)
		}
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(tlsCreds))
	} else {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	// handle unix socket
	if s.HTTPServiceConfiguration.ServerSocketPath != "" {
		address = s.HTTPServiceConfiguration.ServerSocketPath
		dialOpts = append(dialOpts, grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}))
	} else {
		address = net.JoinHostPort("localhost", fmt.Sprintf("%d", s.HTTPServiceConfiguration.Port))
	}
	grpcServer := grpc.NewServer()
	gen.RegisterServiceServer(grpcServer, s.GRPCService)
	err = gen.RegisterServiceHandlerFromEndpoint(
		ctx,
		mux,
		address,
		dialOpts,
	)
	if err != nil {
		log.Fatal(err)
	}
	return grpcServer
}

func (s *HTTPService) ServeHTTP(mux *runtime.ServeMux) *http.Server {
	server := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 60 * time.Second,
	}

	return server
}

func (s *HTTPService) Serve(ctx context.Context, eval eval.IEvaluator) error {
	s.GRPCService.Eval = eval

	g, gCtx := errgroup.WithContext(ctx)

	// Mux Setup
	mux := runtime.NewServeMux(
		runtime.WithErrorHandler(s.HTTPErrorHandler),
	)

	// GRPC Setup
	grpcServer := s.ServerGRPC(ctx, mux)
	// HTTP Setup
	httpServer := s.ServeHTTP(mux)
	// Net listener
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.HTTPServiceConfiguration.Port))
	if err != nil {
		log.Fatal(err)
	}

	tcpm := cmux.New(l)
	// We first match on HTTP 1.1 methods.
	httpl := tcpm.Match(cmux.HTTP1Fast())
	// If not matched, we assume that its TLS.
	tlsl := tcpm.Match(cmux.Any())
	if s.HTTPServiceConfiguration.ServerCertPath != "" && s.HTTPServiceConfiguration.ServerKeyPath != "" {
		tlsl = s.tlsListener(tlsl)
	}
	// Now, we build another mux recursively to match HTTPS and GoRPC.
	// You can use the same trick for SSH.
	tlsm := cmux.New(tlsl)
	httpsl := tlsm.Match(cmux.HTTP1Fast())
	// If a socket path has been provided, create a new listener for the grpc service to listen on
	var gorpcl net.Listener
	if s.HTTPServiceConfiguration.ServerSocketPath != "" {
		gorpcl, err = net.Listen("unix", s.HTTPServiceConfiguration.ServerSocketPath)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		gorpcl = tlsm.Match(cmux.Any())
	}

	g.Go(func() error {
		return httpServer.Serve(httpl) // HTTP
	})
	g.Go(func() error {
		return httpServer.Serve(httpsl) // HTTPS
	})
	g.Go(func() error {
		return grpcServer.Serve(gorpcl) // GRPC
	})
	g.Go(func() error {
		return tlsm.Serve() // GRPC
	})
	g.Go(func() error {
		return tcpm.Serve() // GRPC
	})

	<-gCtx.Done()
	grpcServer.GracefulStop()
	if err = httpServer.Shutdown(context.Background()); err != nil {
		return err
	}
	err = g.Wait()
	if err != nil && !errors.Is(err, grpc.ErrServerStopped) && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
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
