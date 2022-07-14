package service

import (
	"context"

	"github.com/open-feature/flagd/pkg/eval"
	"github.com/open-feature/flagd/pkg/model"
	gen "github.com/open-feature/flagd/schemas/protobuf/gen/v1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type GRPCServiceConfiguration struct {
	Port int32
}

type GRPCService struct {
	GRPCServiceConfiguration *GRPCServiceConfiguration
	eval                     eval.IEvaluator
	gen.UnimplementedServiceServer
}

func (s GRPCService) Serve(ctx context.Context, eval eval.IEvaluator) error {
	return nil
}

// TODO: might be able to simplify some of this with generics.
func (s GRPCService) ResolveBoolean(
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

func (s GRPCService) ResolveString(
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

func (s GRPCService) ResolveNumber(
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

func (s GRPCService) ResolveObject(
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

func handleEvaluationError(err error, reason string) error {
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
	})
	if err != nil {
		log.Error(err)
		return st.Err()
	}
	return stWD.Err()
}
