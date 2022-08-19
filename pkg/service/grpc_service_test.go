package service_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	service "github.com/open-feature/flagd/pkg/service"
	gen "go.buf.build/open-feature/flagd-server/open-feature/flagd/schema/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

type resolveBooleanArgs struct {
	evalFields   resolveBooleanEvalFields
	functionArgs resolveBooleanFunctionArgs
	want         *gen.ResolveBooleanResponse
	wantErr      error
}
type resolveBooleanFunctionArgs struct {
	ctx context.Context
	req *gen.ResolveBooleanRequest
}
type resolveBooleanEvalFields struct {
	result  bool
	variant string
	reason  string
	err     error
}

func TestGRPCService_ResolveBoolean(t *testing.T) {
	ctrl := gomock.NewController(t)
	grpcS := service.GRPCService{}
	tests := map[string]resolveBooleanArgs{
		"happy path": {
			evalFields: resolveBooleanEvalFields{
				result:  true,
				variant: "on",
				reason:  "STATIC",
				err:     nil,
			},
			functionArgs: resolveBooleanFunctionArgs{
				context.Background(),
				&gen.ResolveBooleanRequest{
					FlagKey: "bool",
					Context: &structpb.Struct{},
				},
			},
			want: &gen.ResolveBooleanResponse{
				Value:   true,
				Reason:  "STATIC",
				Variant: "on",
			},
			wantErr: nil,
		},
		"eval returns error": {
			evalFields: resolveBooleanEvalFields{
				result:  true,
				variant: ":(",
				reason:  "ERROR",
				err:     errors.New("eval interface error"),
			},
			functionArgs: resolveBooleanFunctionArgs{
				context.Background(),
				&gen.ResolveBooleanRequest{
					FlagKey: "bool",
					Context: &structpb.Struct{},
				},
			},
			want:    &gen.ResolveBooleanResponse{},
			wantErr: grpcS.HandleEvaluationError(errors.New("eval interface error"), "ERROR"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			eval := NewMockIEvaluator(ctrl)
			eval.EXPECT().ResolveBooleanValue(tt.functionArgs.req.FlagKey, gomock.Any()).Return(
				tt.evalFields.result,
				tt.evalFields.variant,
				tt.evalFields.reason,
				tt.evalFields.err,
			).AnyTimes()
			s := service.GRPCService{
				Eval: eval,
			}
			got, err := s.ResolveBoolean(tt.functionArgs.ctx, tt.functionArgs.req)
			if (err != nil) && !errors.Is(err, tt.wantErr) {
				t.Errorf("GRPCService.ResolveBoolean() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCService.ResolveBoolean() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkGRPCService_ResolveBoolean(b *testing.B) {
	ctrl := gomock.NewController(b)
	grpcS := service.GRPCService{}
	tests := map[string]resolveBooleanArgs{
		"happy path": {
			evalFields: resolveBooleanEvalFields{
				result:  true,
				variant: "on",
				reason:  "STATIC",
				err:     nil,
			},
			functionArgs: resolveBooleanFunctionArgs{
				context.Background(),
				&gen.ResolveBooleanRequest{
					FlagKey: "bool",
					Context: &structpb.Struct{},
				},
			},
			want: &gen.ResolveBooleanResponse{
				Value:   true,
				Reason:  "STATIC",
				Variant: "on",
			},
			wantErr: nil,
		},
		"eval returns error": {
			evalFields: resolveBooleanEvalFields{
				result:  true,
				variant: ":(",
				reason:  "ERROR",
				err:     errors.New("eval interface error"),
			},
			functionArgs: resolveBooleanFunctionArgs{
				context.Background(),
				&gen.ResolveBooleanRequest{
					FlagKey: "bool",
					Context: &structpb.Struct{},
				},
			},
			want:    &gen.ResolveBooleanResponse{},
			wantErr: grpcS.HandleEvaluationError(errors.New("eval interface error"), "ERROR"),
		},
	}
	for name, tt := range tests {
		eval := NewMockIEvaluator(ctrl)
		eval.EXPECT().ResolveBooleanValue(tt.functionArgs.req.FlagKey, gomock.Any()).Return(
			tt.evalFields.result,
			tt.evalFields.variant,
			tt.evalFields.reason,
			tt.evalFields.err,
		).AnyTimes()
		s := service.GRPCService{
			Eval: eval,
		}
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				got, err := s.ResolveBoolean(tt.functionArgs.ctx, tt.functionArgs.req)
				if (err != nil) && !errors.Is(err, tt.wantErr) {
					b.Errorf("GRPCService.ResolveBoolean() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					b.Errorf("GRPCService.ResolveBoolean() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

type resolveStringArgs struct {
	evalFields   resolveStringEvalFields
	functionArgs resolveStringFunctionArgs
	want         *gen.ResolveStringResponse
	wantErr      error
}
type resolveStringFunctionArgs struct {
	ctx context.Context
	req *gen.ResolveStringRequest
}
type resolveStringEvalFields struct {
	result  string
	variant string
	reason  string
	err     error
}

func TestGRPCService_ResolveString(t *testing.T) {
	ctrl := gomock.NewController(t)
	grpcS := service.GRPCService{}
	tests := map[string]resolveStringArgs{
		"happy path": {
			evalFields: resolveStringEvalFields{
				result:  "true",
				variant: "on",
				reason:  "STATIC",
				err:     nil,
			},
			functionArgs: resolveStringFunctionArgs{
				context.Background(),
				&gen.ResolveStringRequest{
					FlagKey: "string",
					Context: &structpb.Struct{},
				},
			},
			want: &gen.ResolveStringResponse{
				Value:   "true",
				Reason:  "STATIC",
				Variant: "on",
			},
			wantErr: nil,
		},
		"eval returns error": {
			evalFields: resolveStringEvalFields{
				result:  "true",
				variant: ":(",
				reason:  "ERROR",
				err:     errors.New("eval interface error"),
			},
			functionArgs: resolveStringFunctionArgs{
				context.Background(),
				&gen.ResolveStringRequest{
					FlagKey: "string",
					Context: &structpb.Struct{},
				},
			},
			want:    &gen.ResolveStringResponse{},
			wantErr: grpcS.HandleEvaluationError(errors.New("eval interface error"), "ERROR"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			eval := NewMockIEvaluator(ctrl)
			eval.EXPECT().ResolveStringValue(tt.functionArgs.req.FlagKey, gomock.Any()).Return(
				tt.evalFields.result,
				tt.evalFields.variant,
				tt.evalFields.reason,
				tt.evalFields.err,
			)
			s := service.GRPCService{
				Eval: eval,
			}
			got, err := s.ResolveString(tt.functionArgs.ctx, tt.functionArgs.req)
			if (err != nil) && !errors.Is(err, tt.wantErr) {
				t.Errorf("GRPCService.ResolveString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCService.ResolveString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkGRPCService_ResolveString(b *testing.B) {
	ctrl := gomock.NewController(b)
	grpcS := service.GRPCService{}
	tests := map[string]resolveStringArgs{
		"happy path": {
			evalFields: resolveStringEvalFields{
				result:  "true",
				variant: "on",
				reason:  "STATIC",
				err:     nil,
			},
			functionArgs: resolveStringFunctionArgs{
				context.Background(),
				&gen.ResolveStringRequest{
					FlagKey: "string",
					Context: &structpb.Struct{},
				},
			},
			want: &gen.ResolveStringResponse{
				Value:   "true",
				Reason:  "STATIC",
				Variant: "on",
			},
			wantErr: nil,
		},
		"eval returns error": {
			evalFields: resolveStringEvalFields{
				result:  "true",
				variant: ":(",
				reason:  "ERROR",
				err:     errors.New("eval interface error"),
			},
			functionArgs: resolveStringFunctionArgs{
				context.Background(),
				&gen.ResolveStringRequest{
					FlagKey: "string",
					Context: &structpb.Struct{},
				},
			},
			want:    &gen.ResolveStringResponse{},
			wantErr: grpcS.HandleEvaluationError(errors.New("eval interface error"), "ERROR"),
		},
	}
	for name, tt := range tests {
		eval := NewMockIEvaluator(ctrl)
		eval.EXPECT().ResolveStringValue(tt.functionArgs.req.FlagKey, gomock.Any()).Return(
			tt.evalFields.result,
			tt.evalFields.variant,
			tt.evalFields.reason,
			tt.evalFields.err,
		).AnyTimes()
		s := service.GRPCService{
			Eval: eval,
		}
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				got, err := s.ResolveString(tt.functionArgs.ctx, tt.functionArgs.req)
				if (err != nil) && !errors.Is(err, tt.wantErr) {
					b.Errorf("GRPCService.ResolveString() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					b.Errorf("GRPCService.ResolveString() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

type resolveFloatArgs struct {
	evalFields   resolveFloatEvalFields
	functionArgs resolveFloatFunctionArgs
	want         *gen.ResolveFloatResponse
	wantErr      error
}
type resolveFloatFunctionArgs struct {
	ctx context.Context
	req *gen.ResolveFloatRequest
}
type resolveFloatEvalFields struct {
	result  float64
	variant string
	reason  string
	err     error
}

func TestGRPCService_ResolveFloat(t *testing.T) {
	ctrl := gomock.NewController(t)
	grpcS := service.GRPCService{}
	tests := map[string]resolveFloatArgs{
		"happy path": {
			evalFields: resolveFloatEvalFields{
				result:  12,
				variant: "on",
				reason:  "STATIC",
				err:     nil,
			},
			functionArgs: resolveFloatFunctionArgs{
				context.Background(),
				&gen.ResolveFloatRequest{
					FlagKey: "float",
					Context: &structpb.Struct{},
				},
			},
			want: &gen.ResolveFloatResponse{
				Value:   12,
				Reason:  "STATIC",
				Variant: "on",
			},
			wantErr: nil,
		},
		"eval returns error": {
			evalFields: resolveFloatEvalFields{
				result:  12,
				variant: ":(",
				reason:  "ERROR",
				err:     errors.New("eval interface error"),
			},
			functionArgs: resolveFloatFunctionArgs{
				context.Background(),
				&gen.ResolveFloatRequest{
					FlagKey: "float",
					Context: &structpb.Struct{},
				},
			},
			want:    &gen.ResolveFloatResponse{},
			wantErr: grpcS.HandleEvaluationError(errors.New("eval interface error"), "ERROR"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			eval := NewMockIEvaluator(ctrl)
			eval.EXPECT().ResolveFloatValue(tt.functionArgs.req.FlagKey, gomock.Any()).Return(
				tt.evalFields.result,
				tt.evalFields.variant,
				tt.evalFields.reason,
				tt.evalFields.err,
			).AnyTimes()
			s := service.GRPCService{
				Eval: eval,
			}
			got, err := s.ResolveFloat(tt.functionArgs.ctx, tt.functionArgs.req)
			if (err != nil) && !errors.Is(err, tt.wantErr) {
				t.Errorf("GRPCService.ResolveNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCService.ResolveNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkGRPCService_ResolveFloat(b *testing.B) {
	ctrl := gomock.NewController(b)
	grpcS := service.GRPCService{}
	tests := map[string]resolveFloatArgs{
		"happy path": {
			evalFields: resolveFloatEvalFields{
				result:  12,
				variant: "on",
				reason:  "STATIC",
				err:     nil,
			},
			functionArgs: resolveFloatFunctionArgs{
				context.Background(),
				&gen.ResolveFloatRequest{
					FlagKey: "float",
					Context: &structpb.Struct{},
				},
			},
			want: &gen.ResolveFloatResponse{
				Value:   12,
				Reason:  "STATIC",
				Variant: "on",
			},
			wantErr: nil,
		},
		"eval returns error": {
			evalFields: resolveFloatEvalFields{
				result:  12,
				variant: ":(",
				reason:  "ERROR",
				err:     errors.New("eval interface error"),
			},
			functionArgs: resolveFloatFunctionArgs{
				context.Background(),
				&gen.ResolveFloatRequest{
					FlagKey: "float",
					Context: &structpb.Struct{},
				},
			},
			want:    &gen.ResolveFloatResponse{},
			wantErr: grpcS.HandleEvaluationError(errors.New("eval interface error"), "ERROR"),
		},
	}
	for name, tt := range tests {
		eval := NewMockIEvaluator(ctrl)
		eval.EXPECT().ResolveFloatValue(tt.functionArgs.req.FlagKey, gomock.Any()).Return(
			tt.evalFields.result,
			tt.evalFields.variant,
			tt.evalFields.reason,
			tt.evalFields.err,
		).AnyTimes()
		s := service.GRPCService{
			Eval: eval,
		}
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				got, err := s.ResolveFloat(tt.functionArgs.ctx, tt.functionArgs.req)
				if (err != nil) && !errors.Is(err, tt.wantErr) {
					b.Errorf("GRPCService.ResolveNumber() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					b.Errorf("GRPCService.ResolveNumber() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

type resolveIntArgs struct {
	evalFields   resolveIntEvalFields
	functionArgs resolveIntFunctionArgs
	want         *gen.ResolveIntResponse
	wantErr      error
}
type resolveIntFunctionArgs struct {
	ctx context.Context
	req *gen.ResolveIntRequest
}
type resolveIntEvalFields struct {
	result  int64
	variant string
	reason  string
	err     error
}

func TestGRPCService_ResolveInt(t *testing.T) {
	ctrl := gomock.NewController(t)
	grpcS := service.GRPCService{}
	tests := map[string]resolveIntArgs{
		"happy path": {
			evalFields: resolveIntEvalFields{
				result:  12,
				variant: "on",
				reason:  "STATIC",
				err:     nil,
			},
			functionArgs: resolveIntFunctionArgs{
				context.Background(),
				&gen.ResolveIntRequest{
					FlagKey: "int",
					Context: &structpb.Struct{},
				},
			},
			want: &gen.ResolveIntResponse{
				Value:   12,
				Reason:  "STATIC",
				Variant: "on",
			},
			wantErr: nil,
		},
		"eval returns error": {
			evalFields: resolveIntEvalFields{
				result:  12,
				variant: ":(",
				reason:  "ERROR",
				err:     errors.New("eval interface error"),
			},
			functionArgs: resolveIntFunctionArgs{
				context.Background(),
				&gen.ResolveIntRequest{
					FlagKey: "int",
					Context: &structpb.Struct{},
				},
			},
			want:    &gen.ResolveIntResponse{},
			wantErr: grpcS.HandleEvaluationError(errors.New("eval interface error"), "ERROR"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			eval := NewMockIEvaluator(ctrl)
			eval.EXPECT().ResolveIntValue(tt.functionArgs.req.FlagKey, gomock.Any()).Return(
				tt.evalFields.result,
				tt.evalFields.variant,
				tt.evalFields.reason,
				tt.evalFields.err,
			).AnyTimes()
			s := service.GRPCService{
				Eval: eval,
			}
			got, err := s.ResolveInt(tt.functionArgs.ctx, tt.functionArgs.req)
			if (err != nil) && !errors.Is(err, tt.wantErr) {
				t.Errorf("GRPCService.ResolveNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCService.ResolveNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkGRPCService_ResolveInt(b *testing.B) {
	ctrl := gomock.NewController(b)
	grpcS := service.GRPCService{}
	tests := map[string]resolveIntArgs{
		"happy path": {
			evalFields: resolveIntEvalFields{
				result:  12,
				variant: "on",
				reason:  "STATIC",
				err:     nil,
			},
			functionArgs: resolveIntFunctionArgs{
				context.Background(),
				&gen.ResolveIntRequest{
					FlagKey: "int",
					Context: &structpb.Struct{},
				},
			},
			want: &gen.ResolveIntResponse{
				Value:   12,
				Reason:  "STATIC",
				Variant: "on",
			},
			wantErr: nil,
		},
		"eval returns error": {
			evalFields: resolveIntEvalFields{
				result:  12,
				variant: ":(",
				reason:  "ERROR",
				err:     errors.New("eval interface error"),
			},
			functionArgs: resolveIntFunctionArgs{
				context.Background(),
				&gen.ResolveIntRequest{
					FlagKey: "int",
					Context: &structpb.Struct{},
				},
			},
			want:    &gen.ResolveIntResponse{},
			wantErr: grpcS.HandleEvaluationError(errors.New("eval interface error"), "ERROR"),
		},
	}
	for name, tt := range tests {
		eval := NewMockIEvaluator(ctrl)
		eval.EXPECT().ResolveIntValue(tt.functionArgs.req.FlagKey, gomock.Any()).Return(
			tt.evalFields.result,
			tt.evalFields.variant,
			tt.evalFields.reason,
			tt.evalFields.err,
		).AnyTimes()
		s := service.GRPCService{
			Eval: eval,
		}
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				got, err := s.ResolveInt(tt.functionArgs.ctx, tt.functionArgs.req)
				if (err != nil) && !errors.Is(err, tt.wantErr) {
					b.Errorf("GRPCService.ResolveNumber() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					b.Errorf("GRPCService.ResolveNumber() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

type resolveObjectArgs struct {
	evalFields   resolveObjectEvalFields
	functionArgs resolveObjectFunctionArgs
	want         *gen.ResolveObjectResponse
	wantErr      error
}
type resolveObjectFunctionArgs struct {
	ctx context.Context
	req *gen.ResolveObjectRequest
}
type resolveObjectEvalFields struct {
	result  map[string]interface{}
	variant string
	reason  string
	err     error
}

func TestGRPCService_ResolveObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	grpcS := service.GRPCService{}
	tests := map[string]resolveObjectArgs{
		"happy path": {
			evalFields: resolveObjectEvalFields{
				result: map[string]interface{}{
					"food": "bars",
				},
				variant: "on",
				reason:  "STATIC",
				err:     nil,
			},
			functionArgs: resolveObjectFunctionArgs{
				context.Background(),
				&gen.ResolveObjectRequest{
					FlagKey: "object",
					Context: &structpb.Struct{},
				},
			},
			want: &gen.ResolveObjectResponse{
				Value:   nil,
				Reason:  "STATIC",
				Variant: "on",
			},
			wantErr: nil,
		},
		"eval returns error": {
			evalFields: resolveObjectEvalFields{
				result: map[string]interface{}{
					"food": "bars",
				},
				variant: ":(",
				reason:  "ERROR",
				err:     errors.New("eval interface error"),
			},
			functionArgs: resolveObjectFunctionArgs{
				context.Background(),
				&gen.ResolveObjectRequest{
					FlagKey: "object",
					Context: &structpb.Struct{},
				},
			},
			want:    &gen.ResolveObjectResponse{},
			wantErr: grpcS.HandleEvaluationError(errors.New("eval interface error"), "ERROR"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			eval := NewMockIEvaluator(ctrl)
			eval.EXPECT().ResolveObjectValue(tt.functionArgs.req.FlagKey, gomock.Any()).Return(
				tt.evalFields.result,
				tt.evalFields.variant,
				tt.evalFields.reason,
				tt.evalFields.err,
			).AnyTimes()
			s := service.GRPCService{
				Eval: eval,
			}

			if name != "eval returns error" {
				outParsed, err := structpb.NewStruct(tt.evalFields.result)
				if err != nil {
					t.Error(err)
				}
				tt.want.Value = outParsed
			}
			got, err := s.ResolveObject(tt.functionArgs.ctx, tt.functionArgs.req)
			if (err != nil) && !errors.Is(err, tt.wantErr) {
				t.Errorf("GRPCService.ResolveObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Value.AsMap(), tt.want.Value.AsMap()) {
				t.Errorf("GRPCService.ResolveObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkGRPCService_ResolveObject(b *testing.B) {
	ctrl := gomock.NewController(b)
	grpcS := service.GRPCService{}
	tests := map[string]resolveObjectArgs{
		"happy path": {
			evalFields: resolveObjectEvalFields{
				result: map[string]interface{}{
					"food": "bars",
				},
				variant: "on",
				reason:  "STATIC",
				err:     nil,
			},
			functionArgs: resolveObjectFunctionArgs{
				context.Background(),
				&gen.ResolveObjectRequest{
					FlagKey: "object",
					Context: &structpb.Struct{},
				},
			},
			want: &gen.ResolveObjectResponse{
				Value:   nil,
				Reason:  "STATIC",
				Variant: "on",
			},
			wantErr: nil,
		},
		"eval returns error": {
			evalFields: resolveObjectEvalFields{
				result: map[string]interface{}{
					"food": "bars",
				},
				variant: ":(",
				reason:  "ERROR",
				err:     errors.New("eval interface error"),
			},
			functionArgs: resolveObjectFunctionArgs{
				context.Background(),
				&gen.ResolveObjectRequest{
					FlagKey: "object",
					Context: &structpb.Struct{},
				},
			},
			want:    &gen.ResolveObjectResponse{},
			wantErr: grpcS.HandleEvaluationError(errors.New("eval interface error"), "ERROR"),
		},
	}
	for name, tt := range tests {
		eval := NewMockIEvaluator(ctrl)
		eval.EXPECT().ResolveObjectValue(tt.functionArgs.req.FlagKey, gomock.Any()).Return(
			tt.evalFields.result,
			tt.evalFields.variant,
			tt.evalFields.reason,
			tt.evalFields.err,
		).AnyTimes()
		s := service.GRPCService{
			Eval: eval,
		}
		if name != "eval returns error" {
			outParsed, err := structpb.NewStruct(tt.evalFields.result)
			if err != nil {
				b.Error(err)
			}
			tt.want.Value = outParsed
		}
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				got, err := s.ResolveObject(tt.functionArgs.ctx, tt.functionArgs.req)
				if (err != nil) && !errors.Is(err, tt.wantErr) {
					b.Errorf("GRPCService.ResolveObject() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got.Value.AsMap(), tt.want.Value.AsMap()) {
					b.Errorf("GRPCService.ResolveObject() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
