package service

import (
	"context"
	"fmt"
	"net"

	"github.com/open-feature/flagd/pkg/eval"
	"github.com/open-feature/flagd/pkg/model"
	gen "github.com/open-feature/flagd/schemas/protobuf/proto/go-server/schema/v1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type GRPCServiceConfiguration struct {
	Port           int32
	ServerKeyPath  string
	ServerCertPath string
}

type GRPCService struct {
	GRPCServiceConfiguration *GRPCServiceConfiguration
	eval                     eval.IEvaluator
	gen.UnimplementedServiceServer
	Logger *log.Entry
}

// Serve allows for the use of GRPC only without HTTPP, where as HTTP service enables both
// GRPC and http
func (s *GRPCService) Serve(ctx context.Context, eval eval.IEvaluator) error {
	s.eval = eval

	//TODO: Needs TLS implementation
	var grpcServer = grpc.NewServer()
	gen.RegisterServiceServer(grpcServer, s)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.GRPCServiceConfiguration.Port))
	if err != nil {
		return err
	}
	return grpcServer.Serve(lis)
}

// TODO: might be able to simplify some of this with generics.
func (s *GRPCService) ResolveBoolean(
	ctx context.Context,
	req *gen.ResolveBooleanRequest,
) (*gen.ResolveBooleanResponse, error) {
	res := gen.ResolveBooleanResponse{}
	result, variant, reason, err := s.eval.ResolveBooleanValue(req.GetFlagKey(), req.GetContext())
	if err != nil {
		return &res, s.HandleEvaluationError(err, reason)
	}
	res.Reason = reason
	res.Value = result
	res.Variant = variant
	return &res, nil
}

func (s *GRPCService) ResolveString(
	ctx context.Context,
	req *gen.ResolveStringRequest,
) (*gen.ResolveStringResponse, error) {
	res := gen.ResolveStringResponse{}
	result, variant, reason, err := s.eval.ResolveStringValue(req.GetFlagKey(), req.GetContext())
	if err != nil {
		return &res, s.HandleEvaluationError(err, reason)
	}
	res.Reason = reason
	res.Value = result
	res.Variant = variant
	return &res, nil
}

func (s *GRPCService) ResolveNumber(
	ctx context.Context,
	req *gen.ResolveNumberRequest,
) (*gen.ResolveNumberResponse, error) {
	res := gen.ResolveNumberResponse{}
	result, variant, reason, err := s.eval.ResolveNumberValue(req.GetFlagKey(), req.GetContext())
	if err != nil {
		return &res, s.HandleEvaluationError(err, reason)
	}
	res.Reason = reason
	res.Value = result
	res.Variant = variant
	return &res, nil
}

func (s *GRPCService) ResolveObject(
	ctx context.Context,
	req *gen.ResolveObjectRequest,
) (*gen.ResolveObjectResponse, error) {
	res := gen.ResolveObjectResponse{}
	result, variant, reason, err := s.eval.ResolveObjectValue(req.GetFlagKey(), req.GetContext())
	if err != nil {
		return &res, s.HandleEvaluationError(err, reason)
	}
	val, err := structpb.NewStruct(result)
	if err != nil {
		return &res, s.HandleEvaluationError(err, reason)
	}
	res.Reason = reason
	res.Value = val
	res.Variant = variant
	return &res, nil
}

func (s *GRPCService) HandleEvaluationError(err error, reason string) error {
	statusCode := codes.Internal
	message := err.Error()
	switch message {
	case model.FlagNotFoundErrorCode:
		statusCode = codes.NotFound
	case model.TypeMismatchErrorCode:
		statusCode = codes.InvalidArgument
	}
	st := status.New(statusCode, message)
	stWD, err := st.WithDetails(&gen.ErrorResponse{
		ErrorCode: message,
		Reason:    "ERROR",
	})
	if err != nil {
		s.Logger.Error(err)
		return st.Err()
	}
	return stWD.Err()
}
