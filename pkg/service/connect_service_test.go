package service_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/golang/mock/gomock"
	"github.com/open-feature/flagd/pkg/model"
	service "github.com/open-feature/flagd/pkg/service"
	log "github.com/sirupsen/logrus"
	gen "go.buf.build/open-feature/flagd-connect/open-feature/flagd-dev/schema/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
}

func TestConnectService_UnixConnection(t *testing.T) {
	type evalFields struct {
		result  bool
		variant string
		reason  string
		err     error
	}

	tests := []struct {
		name       string
		socketPath string
		evalFields evalFields
		req        *gen.ResolveBooleanRequest
		want       *gen.ResolveBooleanResponse
		wantErr    error
	}{
		{
			name:       "happy path",
			socketPath: "/tmp/flagd.sock",
			evalFields: evalFields{
				result:  true,
				variant: "on",
				reason:  model.DefaultReason,
				err:     nil,
			},
			req: &gen.ResolveBooleanRequest{
				FlagKey: "myBoolFlag",
				Context: &structpb.Struct{},
			},
			want: &gen.ResolveBooleanResponse{
				Value:   true,
				Reason:  model.DefaultReason,
				Variant: "on",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			eval := NewMockIEvaluator(ctrl)
			eval.EXPECT().ResolveBooleanValue(tt.req.FlagKey, gomock.Any()).Return(
				tt.evalFields.result,
				tt.evalFields.variant,
				tt.evalFields.reason,
				tt.evalFields.err,
			).AnyTimes()
			service := service.ConnectService{
				ConnectServiceConfiguration: &service.ConnectServiceConfiguration{
					ServerSocketPath: tt.socketPath,
				},
			}
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			go func() { _ = service.Serve(ctx, eval) }()
			conn, err := grpc.Dial(
				fmt.Sprintf("unix://%s", tt.socketPath),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithBlock(),
				grpc.WithTimeout(2*time.Second),
			)
			if err != nil {
				log.Errorf("grpc - fail to dial: %v", err)
				return
			}
			client := gen.NewServiceClient(
				conn,
			)
			res, err := client.ResolveBoolean(ctx, tt.req)
			if (err != nil) && !errors.Is(err, tt.wantErr) {
				t.Errorf("ConnectService.ResolveBoolean() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(res.Reason, tt.want.Reason) {
				t.Errorf("ConnectService.ResolveBoolean() = %v, want %v", res, tt.want)
			}
			if !reflect.DeepEqual(res.Value, tt.want.Value) {
				t.Errorf("ConnectService.ResolveBoolean() = %v, want %v", res, tt.want)
			}
			if !reflect.DeepEqual(res.Variant, tt.want.Variant) {
				t.Errorf("ConnectService.ResolveBoolean() = %v, want %v", res, tt.want)
			}
		})
	}
}

func TestConnectService_ResolveBoolean(t *testing.T) {
	ctrl := gomock.NewController(t)
	tests := map[string]resolveBooleanArgs{
		"happy path": {
			evalFields: resolveBooleanEvalFields{
				result:  true,
				variant: "on",
				reason:  model.DefaultReason,
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
				Reason:  model.DefaultReason,
				Variant: "on",
			},
			wantErr: nil,
		},
		"eval returns error": {
			evalFields: resolveBooleanEvalFields{
				result:  true,
				variant: ":(",
				reason:  model.ErrorReason,
			},
			functionArgs: resolveBooleanFunctionArgs{
				context.Background(),
				&gen.ResolveBooleanRequest{
					FlagKey: "bool",
					Context: &structpb.Struct{},
				},
			},
			want: &gen.ResolveBooleanResponse{
				Reason: model.ErrorReason,
			},
			wantErr: errors.New("eval interface error"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			eval := NewMockIEvaluator(ctrl)
			eval.EXPECT().ResolveBooleanValue(tt.functionArgs.req.FlagKey, gomock.Any()).Return(
				tt.evalFields.result,
				tt.evalFields.variant,
				tt.evalFields.reason,
				tt.wantErr,
			).AnyTimes()
			s := service.ConnectService{
				Eval: eval,
			}
			got, err := s.ResolveBoolean(tt.functionArgs.ctx, connect.NewRequest(tt.functionArgs.req))
			if err != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("ConnectService.ResolveBoolean() error = %v, wantErr %v", err.Error(), tt.wantErr.Error())
				return
			}
			if !reflect.DeepEqual(got.Msg, tt.want) {
				t.Errorf("ConnectService.ResolveBoolean() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkConnectService_ResolveBoolean(b *testing.B) {
	ctrl := gomock.NewController(b)
	tests := map[string]resolveBooleanArgs{
		"happy path": {
			evalFields: resolveBooleanEvalFields{
				result:  true,
				variant: "on",
				reason:  model.DefaultReason,
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
				Reason:  model.DefaultReason,
				Variant: "on",
			},
			wantErr: nil,
		},
	}
	for name, tt := range tests {
		eval := NewMockIEvaluator(ctrl)
		eval.EXPECT().ResolveBooleanValue(tt.functionArgs.req.FlagKey, gomock.Any()).Return(
			tt.evalFields.result,
			tt.evalFields.variant,
			tt.evalFields.reason,
			tt.wantErr,
		).AnyTimes()
		s := service.ConnectService{
			Eval: eval,
		}
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				got, err := s.ResolveBoolean(tt.functionArgs.ctx, connect.NewRequest(tt.functionArgs.req))
				if (err != nil) && !errors.Is(err, tt.wantErr) {
					b.Errorf("ConnectService.ResolveBoolean() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got.Msg, tt.want) {
					b.Errorf("ConnectService.ResolveBoolean() = %v, want %v", got, tt.want)
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
}

func TestConnectService_ResolveString(t *testing.T) {
	ctrl := gomock.NewController(t)
	tests := map[string]resolveStringArgs{
		"happy path": {
			evalFields: resolveStringEvalFields{
				result:  "true",
				variant: "on",
				reason:  model.DefaultReason,
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
				Reason:  model.DefaultReason,
				Variant: "on",
			},
			wantErr: nil,
		},
		"eval returns error": {
			evalFields: resolveStringEvalFields{
				result:  "true",
				variant: ":(",
				reason:  model.ErrorReason,
			},
			functionArgs: resolveStringFunctionArgs{
				context.Background(),
				&gen.ResolveStringRequest{
					FlagKey: "string",
					Context: &structpb.Struct{},
				},
			},
			want: &gen.ResolveStringResponse{
				Reason: model.ErrorReason,
			},
			wantErr: errors.New("eval interface error"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			eval := NewMockIEvaluator(ctrl)
			eval.EXPECT().ResolveStringValue(tt.functionArgs.req.FlagKey, gomock.Any()).Return(
				tt.evalFields.result,
				tt.evalFields.variant,
				tt.evalFields.reason,
				tt.wantErr,
			)
			s := service.ConnectService{
				Eval: eval,
			}
			got, err := s.ResolveString(tt.functionArgs.ctx, connect.NewRequest(tt.functionArgs.req))
			if (err != nil) && !errors.Is(err, tt.wantErr) {
				t.Errorf("ConnectService.ResolveString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Msg, tt.want) {
				t.Errorf("ConnectService.ResolveString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkConnectService_ResolveString(b *testing.B) {
	ctrl := gomock.NewController(b)
	tests := map[string]resolveStringArgs{
		"happy path": {
			evalFields: resolveStringEvalFields{
				result:  "true",
				variant: "on",
				reason:  model.DefaultReason,
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
				Reason:  model.DefaultReason,
				Variant: "on",
			},
			wantErr: nil,
		},
	}
	for name, tt := range tests {
		eval := NewMockIEvaluator(ctrl)
		eval.EXPECT().ResolveStringValue(tt.functionArgs.req.FlagKey, gomock.Any()).Return(
			tt.evalFields.result,
			tt.evalFields.variant,
			tt.evalFields.reason,
			tt.wantErr,
		).AnyTimes()
		s := service.ConnectService{
			Eval: eval,
		}
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				got, err := s.ResolveString(tt.functionArgs.ctx, connect.NewRequest(tt.functionArgs.req))
				if (err != nil) && !errors.Is(err, tt.wantErr) {
					b.Errorf("ConnectService.ResolveString() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got.Msg, tt.want) {
					b.Errorf("ConnectService.ResolveString() = %v, want %v", got, tt.want)
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
}

func TestConnectService_ResolveFloat(t *testing.T) {
	ctrl := gomock.NewController(t)
	tests := map[string]resolveFloatArgs{
		"happy path": {
			evalFields: resolveFloatEvalFields{
				result:  12,
				variant: "on",
				reason:  model.DefaultReason,
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
				Reason:  model.DefaultReason,
				Variant: "on",
			},
			wantErr: nil,
		},
		"eval returns error": {
			evalFields: resolveFloatEvalFields{
				result:  12,
				variant: ":(",
				reason:  model.ErrorReason,
			},
			functionArgs: resolveFloatFunctionArgs{
				context.Background(),
				&gen.ResolveFloatRequest{
					FlagKey: "float",
					Context: &structpb.Struct{},
				},
			},
			want: &gen.ResolveFloatResponse{
				Reason: model.ErrorReason,
			},
			wantErr: errors.New("eval interface error"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			eval := NewMockIEvaluator(ctrl)
			eval.EXPECT().ResolveFloatValue(tt.functionArgs.req.FlagKey, gomock.Any()).Return(
				tt.evalFields.result,
				tt.evalFields.variant,
				tt.evalFields.reason,
				tt.wantErr,
			).AnyTimes()
			s := service.ConnectService{
				Eval: eval,
			}
			got, err := s.ResolveFloat(tt.functionArgs.ctx, connect.NewRequest(tt.functionArgs.req))
			if (err != nil) && !errors.Is(err, tt.wantErr) {
				t.Errorf("ConnectService.ResolveNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Msg, tt.want) {
				t.Errorf("ConnectService.ResolveNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkConnectService_ResolveFloat(b *testing.B) {
	ctrl := gomock.NewController(b)
	tests := map[string]resolveFloatArgs{
		"happy path": {
			evalFields: resolveFloatEvalFields{
				result:  12,
				variant: "on",
				reason:  model.DefaultReason,
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
				Reason:  model.DefaultReason,
				Variant: "on",
			},
			wantErr: nil,
		},
	}
	for name, tt := range tests {
		eval := NewMockIEvaluator(ctrl)
		eval.EXPECT().ResolveFloatValue(tt.functionArgs.req.FlagKey, gomock.Any()).Return(
			tt.evalFields.result,
			tt.evalFields.variant,
			tt.evalFields.reason,
			tt.wantErr,
		).AnyTimes()
		s := service.ConnectService{
			Eval: eval,
		}
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				got, err := s.ResolveFloat(tt.functionArgs.ctx, connect.NewRequest(tt.functionArgs.req))
				if (err != nil) && !errors.Is(err, tt.wantErr) {
					b.Errorf("ConnectService.ResolveNumber() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got.Msg, tt.want) {
					b.Errorf("ConnectService.ResolveNumber() = %v, want %v", got, tt.want)
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
}

func TestConnectService_ResolveInt(t *testing.T) {
	ctrl := gomock.NewController(t)
	tests := map[string]resolveIntArgs{
		"happy path": {
			evalFields: resolveIntEvalFields{
				result:  12,
				variant: "on",
				reason:  model.DefaultReason,
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
				Reason:  model.DefaultReason,
				Variant: "on",
			},
			wantErr: nil,
		},
		"eval returns error": {
			evalFields: resolveIntEvalFields{
				result:  12,
				variant: ":(",
				reason:  model.ErrorReason,
			},
			functionArgs: resolveIntFunctionArgs{
				context.Background(),
				&gen.ResolveIntRequest{
					FlagKey: "int",
					Context: &structpb.Struct{},
				},
			},
			want: &gen.ResolveIntResponse{
				Reason: model.ErrorReason,
			},
			wantErr: errors.New("eval interface error"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			eval := NewMockIEvaluator(ctrl)
			eval.EXPECT().ResolveIntValue(tt.functionArgs.req.FlagKey, gomock.Any()).Return(
				tt.evalFields.result,
				tt.evalFields.variant,
				tt.evalFields.reason,
				tt.wantErr,
			).AnyTimes()
			s := service.ConnectService{
				Eval: eval,
			}
			got, err := s.ResolveInt(tt.functionArgs.ctx, connect.NewRequest(tt.functionArgs.req))
			if (err != nil) && !errors.Is(err, tt.wantErr) {
				t.Errorf("ConnectService.ResolveNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Msg, tt.want) {
				t.Errorf("ConnectService.ResolveNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkConnectService_ResolveInt(b *testing.B) {
	ctrl := gomock.NewController(b)
	tests := map[string]resolveIntArgs{
		"happy path": {
			evalFields: resolveIntEvalFields{
				result:  12,
				variant: "on",
				reason:  model.DefaultReason,
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
				Reason:  model.DefaultReason,
				Variant: "on",
			},
			wantErr: nil,
		},
	}
	for name, tt := range tests {
		eval := NewMockIEvaluator(ctrl)
		eval.EXPECT().ResolveIntValue(tt.functionArgs.req.FlagKey, gomock.Any()).Return(
			tt.evalFields.result,
			tt.evalFields.variant,
			tt.evalFields.reason,
			tt.wantErr,
		).AnyTimes()
		s := service.ConnectService{
			Eval: eval,
		}
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				got, err := s.ResolveInt(tt.functionArgs.ctx, connect.NewRequest(tt.functionArgs.req))
				if (err != nil) && !errors.Is(err, tt.wantErr) {
					b.Errorf("ConnectService.ResolveNumber() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got.Msg, tt.want) {
					b.Errorf("ConnectService.ResolveNumber() = %v, want %v", got, tt.want)
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
}

func TestConnectService_ResolveObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	tests := map[string]resolveObjectArgs{
		"happy path": {
			evalFields: resolveObjectEvalFields{
				result: map[string]interface{}{
					"food": "bars",
				},
				variant: "on",
				reason:  model.DefaultReason,
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
				Reason:  model.DefaultReason,
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
				reason:  model.ErrorReason,
			},
			functionArgs: resolveObjectFunctionArgs{
				context.Background(),
				&gen.ResolveObjectRequest{
					FlagKey: "object",
					Context: &structpb.Struct{},
				},
			},
			want: &gen.ResolveObjectResponse{
				Reason: model.ErrorReason,
			},
			wantErr: errors.New("eval interface error"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			eval := NewMockIEvaluator(ctrl)
			eval.EXPECT().ResolveObjectValue(tt.functionArgs.req.FlagKey, gomock.Any()).Return(
				tt.evalFields.result,
				tt.evalFields.variant,
				tt.evalFields.reason,
				tt.wantErr,
			).AnyTimes()
			s := service.ConnectService{
				Eval: eval,
			}

			if name != "eval returns error" {
				outParsed, err := structpb.NewStruct(tt.evalFields.result)
				if err != nil {
					t.Error(err)
				}
				tt.want.Value = outParsed
			}
			got, err := s.ResolveObject(tt.functionArgs.ctx, connect.NewRequest(tt.functionArgs.req))
			if (err != nil) && !errors.Is(err, tt.wantErr) {
				t.Errorf("ConnectService.ResolveObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Msg.Value.AsMap(), tt.want.Value.AsMap()) {
				t.Errorf("ConnectService.ResolveObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkConnectService_ResolveObject(b *testing.B) {
	ctrl := gomock.NewController(b)
	tests := map[string]resolveObjectArgs{
		"happy path": {
			evalFields: resolveObjectEvalFields{
				result: map[string]interface{}{
					"food": "bars",
				},
				variant: "on",
				reason:  model.DefaultReason,
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
				Reason:  model.DefaultReason,
				Variant: "on",
			},
			wantErr: nil,
		},
	}
	for name, tt := range tests {
		eval := NewMockIEvaluator(ctrl)
		eval.EXPECT().ResolveObjectValue(tt.functionArgs.req.FlagKey, gomock.Any()).Return(
			tt.evalFields.result,
			tt.evalFields.variant,
			tt.evalFields.reason,
			tt.wantErr,
		).AnyTimes()
		s := service.ConnectService{
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
				got, err := s.ResolveObject(tt.functionArgs.ctx, connect.NewRequest(tt.functionArgs.req))
				if (err != nil) && !errors.Is(err, tt.wantErr) {
					b.Errorf("ConnectService.ResolveObject() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got.Msg.Value.AsMap(), tt.want.Value.AsMap()) {
					b.Errorf("ConnectService.ResolveObject() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
