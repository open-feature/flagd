package service_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	schemaGrpcV1 "buf.build/gen/go/open-feature/flagd/grpc/go/schema/v1/schemav1grpc"
	schemaV1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/schema/v1"
	"github.com/bufbuild/connect-go"
	"github.com/golang/mock/gomock"
	"github.com/open-feature/flagd/core/pkg/eval"
	mock "github.com/open-feature/flagd/core/pkg/eval/mock"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	iservice "github.com/open-feature/flagd/core/pkg/service"
	service "github.com/open-feature/flagd/core/pkg/service/flag-evaluation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"
)

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
		req        *schemaV1.ResolveBooleanRequest
		want       *schemaV1.ResolveBooleanResponse
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
			req: &schemaV1.ResolveBooleanRequest{
				FlagKey: "myBoolFlag",
				Context: &structpb.Struct{},
			},
			want: &schemaV1.ResolveBooleanResponse{
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
			eval := mock.NewMockIEvaluator(ctrl)
			eval.EXPECT().ResolveBooleanValue(gomock.Any(), tt.req.FlagKey, gomock.Any()).Return(
				tt.evalFields.result,
				tt.evalFields.variant,
				tt.evalFields.reason,
				tt.evalFields.err,
			).AnyTimes()
			svc := service.ConnectService{
				ConnectServiceConfiguration: &service.ConnectServiceConfiguration{
					ServerSocketPath: tt.socketPath,
				},
				Logger: logger.NewLogger(nil, false),
			}
			serveConf := iservice.Configuration{
				ReadinessProbe: func() bool {
					return true
				},
			}
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			go func() {
				err := svc.Serve(ctx, eval, serveConf)
				fmt.Println(err)
			}()
			conn, err := grpc.Dial(
				fmt.Sprintf("unix://%s", tt.socketPath),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithBlock(),
				grpc.WithTimeout(2*time.Second),
			)
			if err != nil {
				t.Errorf("grpc - fail to dial: %v", err)
				return
			}
			client := schemaGrpcV1.NewServiceClient(
				conn,
			)

			res, err := client.ResolveBoolean(ctx, tt.req)
			if (err != nil) && !errors.Is(err, tt.wantErr) {
				t.Errorf("ConnectService.ResolveBoolean() error = %v, wantErr %v", err, tt.wantErr)
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

func TestConnectService_ResolveAll(t *testing.T) {
	tests := map[string]struct {
		req     *schemaV1.ResolveAllRequest
		evalRes []eval.AnyValue
		wantErr error
		wantRes *schemaV1.ResolveAllResponse
	}{
		"happy-path": {
			req: &schemaV1.ResolveAllRequest{},
			evalRes: []eval.AnyValue{
				{
					Value:   true,
					Variant: "bool-true",
					Reason:  "true",
					FlagKey: "bool",
				},
				{
					Value:   float64(12.12),
					Variant: "float",
					Reason:  "float",
					FlagKey: "float",
				},
				{
					Value:   "hello",
					Variant: "string",
					Reason:  "string",
					FlagKey: "string",
				},
			},
			wantErr: nil,
			wantRes: &schemaV1.ResolveAllResponse{
				Flags: map[string]*schemaV1.AnyFlag{
					"bool": {
						Value: &schemaV1.AnyFlag_BoolValue{
							BoolValue: true,
						},
						Reason: "STATIC",
					},
					"float": {
						Value: &schemaV1.AnyFlag_DoubleValue{
							DoubleValue: float64(12.12),
						},
						Reason: "STATIC",
					},
					"string": {
						Value: &schemaV1.AnyFlag_StringValue{
							StringValue: "hello",
						},
						Reason: "STATIC",
					},
				},
			},
		},
	}
	ctrl := gomock.NewController(t)
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			eval := mock.NewMockIEvaluator(ctrl)
			eval.EXPECT().ResolveAllValues(gomock.Any(), gomock.Any()).Return(
				tt.evalRes,
			).AnyTimes()
			s := service.ConnectService{
				Eval:   eval,
				Logger: logger.NewLogger(nil, false),
			}
			got, err := s.ResolveAll(context.Background(), connect.NewRequest(tt.req))
			if err != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("ConnectService.ResolveAll() error = %v, wantErr %v", err.Error(), tt.wantErr.Error())
				return
			}
			for _, flag := range tt.evalRes {
				switch v := flag.Value.(type) {
				case bool:
					val := got.Msg.Flags[flag.FlagKey].Value.(*schemaV1.AnyFlag_BoolValue)
					if v != val.BoolValue {
						t.Errorf("ConnectService.ResolveAll(), key %s = %v, want %v", flag.FlagKey, val.BoolValue, v)
					}
				case string:
					val := got.Msg.Flags[flag.FlagKey].Value.(*schemaV1.AnyFlag_StringValue)
					if v != val.StringValue {
						t.Errorf("ConnectService.ResolveAll(), key %s = %s, want %s", flag.FlagKey, val.StringValue, v)
					}
				case float64:
					val := got.Msg.Flags[flag.FlagKey].Value.(*schemaV1.AnyFlag_DoubleValue)
					if v != val.DoubleValue {
						t.Errorf("ConnectService.ResolveAll(), key %s = %f, want %f", flag.FlagKey, val.DoubleValue, v)
					}
				}
			}
		})
	}
}

type resolveBooleanArgs struct {
	evalFields   resolveBooleanEvalFields
	functionArgs resolveBooleanFunctionArgs
	want         *schemaV1.ResolveBooleanResponse
	wantErr      error
}
type resolveBooleanFunctionArgs struct {
	ctx context.Context
	req *schemaV1.ResolveBooleanRequest
}
type resolveBooleanEvalFields struct {
	result  bool
	variant string
	reason  string
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
				&schemaV1.ResolveBooleanRequest{
					FlagKey: "bool",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveBooleanResponse{
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
				&schemaV1.ResolveBooleanRequest{
					FlagKey: "bool",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveBooleanResponse{
				Value:   true,
				Variant: ":(",
				Reason:  model.ErrorReason,
			},
			wantErr: errors.New("eval interface error"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			eval := mock.NewMockIEvaluator(ctrl)
			eval.EXPECT().ResolveBooleanValue(gomock.Any(), tt.functionArgs.req.FlagKey, gomock.Any()).Return(
				tt.evalFields.result,
				tt.evalFields.variant,
				tt.evalFields.reason,
				tt.wantErr,
			).AnyTimes()
			s := service.ConnectService{
				Eval:   eval,
				Logger: logger.NewLogger(nil, false),
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
				&schemaV1.ResolveBooleanRequest{
					FlagKey: "bool",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveBooleanResponse{
				Value:   true,
				Reason:  model.DefaultReason,
				Variant: "on",
			},
			wantErr: nil,
		},
	}
	for name, tt := range tests {
		eval := mock.NewMockIEvaluator(ctrl)
		eval.EXPECT().ResolveBooleanValue(gomock.Any(), tt.functionArgs.req.FlagKey, gomock.Any()).Return(
			tt.evalFields.result,
			tt.evalFields.variant,
			tt.evalFields.reason,
			tt.wantErr,
		).AnyTimes()
		s := service.ConnectService{
			Eval:   eval,
			Logger: logger.NewLogger(nil, false),
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
	want         *schemaV1.ResolveStringResponse
	wantErr      error
}
type resolveStringFunctionArgs struct {
	ctx context.Context
	req *schemaV1.ResolveStringRequest
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
				&schemaV1.ResolveStringRequest{
					FlagKey: "string",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveStringResponse{
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
				&schemaV1.ResolveStringRequest{
					FlagKey: "string",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveStringResponse{
				Value:   "true",
				Variant: ":(",
				Reason:  model.ErrorReason,
			},
			wantErr: errors.New("eval interface error"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			eval := mock.NewMockIEvaluator(ctrl)
			eval.EXPECT().ResolveStringValue(gomock.Any(), tt.functionArgs.req.FlagKey, gomock.Any()).Return(
				tt.evalFields.result,
				tt.evalFields.variant,
				tt.evalFields.reason,
				tt.wantErr,
			)
			s := service.ConnectService{
				Eval:   eval,
				Logger: logger.NewLogger(nil, false),
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
				&schemaV1.ResolveStringRequest{
					FlagKey: "string",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveStringResponse{
				Value:   "true",
				Reason:  model.DefaultReason,
				Variant: "on",
			},
			wantErr: nil,
		},
	}
	for name, tt := range tests {
		eval := mock.NewMockIEvaluator(ctrl)
		eval.EXPECT().ResolveStringValue(gomock.Any(), tt.functionArgs.req.FlagKey, gomock.Any()).Return(
			tt.evalFields.result,
			tt.evalFields.variant,
			tt.evalFields.reason,
			tt.wantErr,
		).AnyTimes()

		s := service.ConnectService{
			Eval:   eval,
			Logger: logger.NewLogger(nil, false),
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
	want         *schemaV1.ResolveFloatResponse
	wantErr      error
}
type resolveFloatFunctionArgs struct {
	ctx context.Context
	req *schemaV1.ResolveFloatRequest
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
				&schemaV1.ResolveFloatRequest{
					FlagKey: "float",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveFloatResponse{
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
				&schemaV1.ResolveFloatRequest{
					FlagKey: "float",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveFloatResponse{
				Value:   12,
				Variant: ":(",
				Reason:  model.ErrorReason,
			},
			wantErr: errors.New("eval interface error"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			eval := mock.NewMockIEvaluator(ctrl)
			eval.EXPECT().ResolveFloatValue(gomock.Any(), tt.functionArgs.req.FlagKey, gomock.Any()).Return(
				tt.evalFields.result,
				tt.evalFields.variant,
				tt.evalFields.reason,
				tt.wantErr,
			).AnyTimes()
			s := service.ConnectService{
				Eval:   eval,
				Logger: logger.NewLogger(nil, false),
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
				&schemaV1.ResolveFloatRequest{
					FlagKey: "float",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveFloatResponse{
				Value:   12,
				Reason:  model.DefaultReason,
				Variant: "on",
			},
			wantErr: nil,
		},
	}
	for name, tt := range tests {
		eval := mock.NewMockIEvaluator(ctrl)
		eval.EXPECT().ResolveFloatValue(gomock.Any(), tt.functionArgs.req.FlagKey, gomock.Any()).Return(
			tt.evalFields.result,
			tt.evalFields.variant,
			tt.evalFields.reason,
			tt.wantErr,
		).AnyTimes()

		s := service.ConnectService{
			Eval:   eval,
			Logger: logger.NewLogger(nil, false),
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
	want         *schemaV1.ResolveIntResponse
	wantErr      error
}
type resolveIntFunctionArgs struct {
	ctx context.Context
	req *schemaV1.ResolveIntRequest
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
				&schemaV1.ResolveIntRequest{
					FlagKey: "int",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveIntResponse{
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
				&schemaV1.ResolveIntRequest{
					FlagKey: "int",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveIntResponse{
				Value:   12,
				Variant: ":(",
				Reason:  model.ErrorReason,
			},
			wantErr: errors.New("eval interface error"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			eval := mock.NewMockIEvaluator(ctrl)
			eval.EXPECT().ResolveIntValue(gomock.Any(), tt.functionArgs.req.FlagKey, gomock.Any()).Return(
				tt.evalFields.result,
				tt.evalFields.variant,
				tt.evalFields.reason,
				tt.wantErr,
			).AnyTimes()
			s := service.ConnectService{
				Eval:   eval,
				Logger: logger.NewLogger(nil, false),
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
				&schemaV1.ResolveIntRequest{
					FlagKey: "int",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveIntResponse{
				Value:   12,
				Reason:  model.DefaultReason,
				Variant: "on",
			},
			wantErr: nil,
		},
	}
	for name, tt := range tests {
		eval := mock.NewMockIEvaluator(ctrl)
		eval.EXPECT().ResolveIntValue(gomock.Any(), tt.functionArgs.req.FlagKey, gomock.Any()).Return(
			tt.evalFields.result,
			tt.evalFields.variant,
			tt.evalFields.reason,
			tt.wantErr,
		).AnyTimes()

		s := service.ConnectService{
			Eval:   eval,
			Logger: logger.NewLogger(nil, false),
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
	want         *schemaV1.ResolveObjectResponse
	wantErr      error
}
type resolveObjectFunctionArgs struct {
	ctx context.Context
	req *schemaV1.ResolveObjectRequest
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
				&schemaV1.ResolveObjectRequest{
					FlagKey: "object",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveObjectResponse{
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
				&schemaV1.ResolveObjectRequest{
					FlagKey: "object",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveObjectResponse{
				Variant: ":(",
				Reason:  model.ErrorReason,
			},
			wantErr: errors.New("eval interface error"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			eval := mock.NewMockIEvaluator(ctrl)
			eval.EXPECT().ResolveObjectValue(gomock.Any(), tt.functionArgs.req.FlagKey, gomock.Any()).Return(
				tt.evalFields.result,
				tt.evalFields.variant,
				tt.evalFields.reason,
				tt.wantErr,
			).AnyTimes()
			s := service.ConnectService{
				Eval:   eval,
				Logger: logger.NewLogger(nil, false),
			}

			outParsed, err := structpb.NewStruct(tt.evalFields.result)
			if err != nil {
				t.Error(err)
			}
			tt.want.Value = outParsed
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
				&schemaV1.ResolveObjectRequest{
					FlagKey: "object",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveObjectResponse{
				Value:   nil,
				Reason:  model.DefaultReason,
				Variant: "on",
			},
			wantErr: nil,
		},
	}
	for name, tt := range tests {
		eval := mock.NewMockIEvaluator(ctrl)
		eval.EXPECT().ResolveObjectValue(gomock.Any(), tt.functionArgs.req.FlagKey, gomock.Any()).Return(
			tt.evalFields.result,
			tt.evalFields.variant,
			tt.evalFields.reason,
			tt.wantErr,
		).AnyTimes()

		s := service.ConnectService{
			Eval:   eval,
			Logger: logger.NewLogger(nil, false),
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
