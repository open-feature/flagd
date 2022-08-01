package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	mock "github.com/open-feature/flagd/pkg/eval/tests/mocks"
	gen "github.com/open-feature/flagd/schemas/protobuf/proto/go-server/schema/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestGRPCService_ResolveBoolean(t *testing.T) {
	type evalFields struct {
		result  bool
		variant string
		reason  string
		err     error
	}
	grpcS := GRPCService{}
	type args struct {
		ctx context.Context
		req *gen.ResolveBooleanRequest
	}
	tests := []struct {
		name       string
		evalFields evalFields
		args       args
		want       *gen.ResolveBooleanResponse
		wantErr    error
	}{
		{
			name: "happy path",
			evalFields: evalFields{
				result:  true,
				variant: "on",
				reason:  "STATIC",
				err:     nil,
			},
			args: args{
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
		{
			name: "eval returns error",
			evalFields: evalFields{
				result:  true,
				variant: ":(",
				reason:  "ERROR",
				err:     errors.New("eval interface error"),
			},
			args: args{
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			eval := mock.NewMockIEvaluator(ctrl)
			eval.EXPECT().ResolveBooleanValue(tt.args.req.FlagKey, gomock.Any()).Return(
				tt.evalFields.result,
				tt.evalFields.variant,
				tt.evalFields.reason,
				tt.evalFields.err,
			).AnyTimes()
			s := GRPCService{
				eval: eval,
			}
			got, err := s.ResolveBoolean(tt.args.ctx, tt.args.req)
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

func TestGRPCService_ResolveString(t *testing.T) {
	type evalFields struct {
		result  string
		variant string
		reason  string
		err     error
	}
	type args struct {
		ctx context.Context
		req *gen.ResolveStringRequest
	}
	grpcS := GRPCService{}
	tests := []struct {
		name       string
		evalFields evalFields
		args       args
		want       *gen.ResolveStringResponse
		wantErr    error
	}{
		{
			name: "happy path",
			evalFields: evalFields{
				result:  "true",
				variant: "on",
				reason:  "STATIC",
				err:     nil,
			},
			args: args{
				context.Background(),
				&gen.ResolveStringRequest{
					FlagKey: "bool",
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
		{
			name: "eval returns error",
			evalFields: evalFields{
				result:  "true",
				variant: ":(",
				reason:  "ERROR",
				err:     errors.New("eval interface error"),
			},
			args: args{
				context.Background(),
				&gen.ResolveStringRequest{
					FlagKey: "bool",
					Context: &structpb.Struct{},
				},
			},
			want:    &gen.ResolveStringResponse{},
			wantErr: grpcS.HandleEvaluationError(errors.New("eval interface error"), "ERROR"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			eval := mock.NewMockIEvaluator(ctrl)
			eval.EXPECT().ResolveStringValue(tt.args.req.FlagKey, gomock.Any()).Return(
				tt.evalFields.result,
				tt.evalFields.variant,
				tt.evalFields.reason,
				tt.evalFields.err,
			).AnyTimes()
			s := GRPCService{
				eval: eval,
			}
			got, err := s.ResolveString(tt.args.ctx, tt.args.req)
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

func TestGRPCService_ResolveNumber(t *testing.T) {
	type evalFields struct {
		result  float32
		variant string
		reason  string
		err     error
	}
	type args struct {
		ctx context.Context
		req *gen.ResolveNumberRequest
	}
	grpcs := GRPCService{}
	tests := []struct {
		name       string
		evalFields evalFields
		args       args
		want       *gen.ResolveNumberResponse
		wantErr    error
	}{
		{
			name: "happy path",
			evalFields: evalFields{
				result:  float32(32),
				variant: "on",
				reason:  "STATIC",
				err:     nil,
			},
			args: args{
				context.Background(),
				&gen.ResolveNumberRequest{
					FlagKey: "bool",
					Context: &structpb.Struct{},
				},
			},
			want: &gen.ResolveNumberResponse{
				Value:   float32(32),
				Reason:  "STATIC",
				Variant: "on",
			},
			wantErr: nil,
		},
		{
			name: "eval returns error",
			evalFields: evalFields{
				result:  float32(32),
				variant: ":(",
				reason:  "ERROR",
				err:     errors.New("eval interface error"),
			},
			args: args{
				context.Background(),
				&gen.ResolveNumberRequest{
					FlagKey: "bool",
					Context: &structpb.Struct{},
				},
			},
			want:    &gen.ResolveNumberResponse{},
			wantErr: grpcs.HandleEvaluationError(errors.New("eval interface error"), "ERROR"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			eval := mock.NewMockIEvaluator(ctrl)
			eval.EXPECT().ResolveNumberValue(tt.args.req.FlagKey, gomock.Any()).Return(
				tt.evalFields.result,
				tt.evalFields.variant,
				tt.evalFields.reason,
				tt.evalFields.err,
			).AnyTimes()
			s := GRPCService{
				eval: eval,
			}
			got, err := s.ResolveNumber(tt.args.ctx, tt.args.req)
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

func TestGRPCService_ResolveObject(t *testing.T) {
	type evalFields struct {
		result  map[string]interface{}
		variant string
		reason  string
		err     error
	}
	type args struct {
		ctx context.Context
		req *gen.ResolveObjectRequest
	}
	grpcs := GRPCService{}
	tests := []struct {
		name       string
		evalFields evalFields
		args       args
		want       *gen.ResolveObjectResponse
		wantErr    error
	}{
		{
			name: "happy path",
			evalFields: evalFields{
				result: map[string]interface{}{
					"food": "bars",
				},
				variant: "on",
				reason:  "STATIC",
				err:     nil,
			},
			args: args{
				context.Background(),
				&gen.ResolveObjectRequest{
					FlagKey: "bool",
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
		{
			name: "eval returns error",
			evalFields: evalFields{
				result: map[string]interface{}{
					"food": "bars",
				},
				variant: ":(",
				reason:  "ERROR",
				err:     errors.New("eval interface error"),
			},
			args: args{
				context.Background(),
				&gen.ResolveObjectRequest{
					FlagKey: "bool",
					Context: &structpb.Struct{},
				},
			},
			want:    &gen.ResolveObjectResponse{},
			wantErr: grpcs.HandleEvaluationError(errors.New("eval interface error"), "ERROR"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			eval := mock.NewMockIEvaluator(ctrl)
			eval.EXPECT().ResolveObjectValue(tt.args.req.FlagKey, gomock.Any()).Return(
				tt.evalFields.result,
				tt.evalFields.variant,
				tt.evalFields.reason,
				tt.evalFields.err,
			).AnyTimes()
			s := GRPCService{
				eval: eval,
			}

			if tt.name != "eval returns error" {
				outParsed, err := structpb.NewStruct(tt.evalFields.result)
				if err != nil {
					t.Error(err)
				}
				tt.want.Value = outParsed
			}

			got, err := s.ResolveObject(tt.args.ctx, tt.args.req)
			if (err != nil) && !errors.Is(err, tt.wantErr) {
				t.Errorf("GRPCService.ResolveObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCService.ResolveObject() = %v, want %v", got, tt.want)
			}
		})
	}
}
