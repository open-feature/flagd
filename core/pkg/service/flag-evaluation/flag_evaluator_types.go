package service

import (
	"fmt"

	schemaV1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/schema/v1"
	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/structpb"
)

type response[T constraints] interface {
	SetResult(value T, variant, reason string, metadata map[string]interface{}) error
}

type constraints interface {
	bool | string | map[string]any | float64 | int64
}

type booleanResponse struct {
	*connect.Response[schemaV1.ResolveBooleanResponse]
}

func (r *booleanResponse) SetResult(value bool, variant, reason string, metadata map[string]interface{}) error {
	r.Msg.Value = value
	r.Msg.Variant = variant
	r.Msg.Reason = reason

	newStruct, err := structpb.NewStruct(metadata)
	if err != nil {
		return fmt.Errorf("failure to wrap metadata %w", err)
	}

	r.Msg.Metadata = newStruct
	return nil
}

type stringResponse struct {
	*connect.Response[schemaV1.ResolveStringResponse]
}

func (r *stringResponse) SetResult(value string, variant, reason string, metadata map[string]interface{}) error {
	r.Msg.Value = value
	r.Msg.Variant = variant
	r.Msg.Reason = reason

	newStruct, err := structpb.NewStruct(metadata)
	if err != nil {
		return fmt.Errorf("failure to wrap metadata %w", err)
	}

	r.Msg.Metadata = newStruct
	return nil
}

type floatResponse struct {
	*connect.Response[schemaV1.ResolveFloatResponse]
}

func (r *floatResponse) SetResult(value float64, variant, reason string, metadata map[string]interface{}) error {
	r.Msg.Value = value
	r.Msg.Variant = variant
	r.Msg.Reason = reason

	newStruct, err := structpb.NewStruct(metadata)
	if err != nil {
		return fmt.Errorf("failure to wrap metadata %w", err)
	}

	r.Msg.Metadata = newStruct
	return nil
}

type intResponse struct {
	*connect.Response[schemaV1.ResolveIntResponse]
}

func (r *intResponse) SetResult(value int64, variant, reason string, metadata map[string]interface{}) error {
	r.Msg.Value = value
	r.Msg.Variant = variant
	r.Msg.Reason = reason

	newStruct, err := structpb.NewStruct(metadata)
	if err != nil {
		return fmt.Errorf("failure to wrap metadata %w", err)
	}

	r.Msg.Metadata = newStruct
	return nil
}

type objectResponse struct {
	*connect.Response[schemaV1.ResolveObjectResponse]
}

func (r *objectResponse) SetResult(value map[string]any, variant, reason string,
	metadata map[string]interface{},
) error {
	r.Msg.Reason = reason
	val, err := structpb.NewStruct(value)
	if err != nil {
		return fmt.Errorf("struct response construction: %w", err)
	}

	r.Msg.Value = val
	r.Msg.Variant = variant

	newStruct, err := structpb.NewStruct(metadata)
	if err != nil {
		return fmt.Errorf("failure to wrap metadata %w", err)
	}

	r.Msg.Metadata = newStruct
	return nil
}
