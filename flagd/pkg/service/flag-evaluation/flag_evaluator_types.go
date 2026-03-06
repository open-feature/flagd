package service

import (
	"fmt"

	evalV1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/flagd/evaluation/v1"
	evalV2 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/flagd/evaluation/v2"
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
	//nolint:staticcheck
	schemaV1Resp *connect.Response[schemaV1.ResolveBooleanResponse]
	evalV1Resp   *connect.Response[evalV1.ResolveBooleanResponse]
}

//nolint:staticcheck
func (r *booleanResponse) SetResult(value bool, variant, reason string, metadata map[string]interface{}) error {
	newStruct, err := structpb.NewStruct(metadata)
	if err != nil {
		return fmt.Errorf("failure to wrap metadata %w", err)
	}

	if r.schemaV1Resp != nil {
		r.schemaV1Resp.Msg.Value = value
		r.schemaV1Resp.Msg.Variant = variant
		r.schemaV1Resp.Msg.Reason = reason
		r.schemaV1Resp.Msg.Metadata = newStruct
	}
	if r.evalV1Resp != nil {
		r.evalV1Resp.Msg.Value = value
		r.evalV1Resp.Msg.Variant = variant
		r.evalV1Resp.Msg.Reason = reason
		r.evalV1Resp.Msg.Metadata = newStruct
	}

	return nil
}

type stringResponse struct {
	//nolint:staticcheck
	schemaV1Resp *connect.Response[schemaV1.ResolveStringResponse]
	evalV1Resp   *connect.Response[evalV1.ResolveStringResponse]
}

//nolint:staticcheck
func (r *stringResponse) SetResult(value string, variant, reason string, metadata map[string]interface{}) error {
	newStruct, err := structpb.NewStruct(metadata)
	if err != nil {
		return fmt.Errorf("failure to wrap metadata %w", err)
	}

	if r.schemaV1Resp != nil {
		r.schemaV1Resp.Msg.Value = value
		r.schemaV1Resp.Msg.Variant = variant
		r.schemaV1Resp.Msg.Reason = reason
		r.schemaV1Resp.Msg.Metadata = newStruct
	}
	if r.evalV1Resp != nil {
		r.evalV1Resp.Msg.Value = value
		r.evalV1Resp.Msg.Variant = variant
		r.evalV1Resp.Msg.Reason = reason
		r.evalV1Resp.Msg.Metadata = newStruct
	}

	return nil
}

type floatResponse struct {
	//nolint:staticcheck
	schemaV1Resp *connect.Response[schemaV1.ResolveFloatResponse]
	evalV1Resp   *connect.Response[evalV1.ResolveFloatResponse]
}

//nolint:staticcheck
func (r *floatResponse) SetResult(value float64, variant, reason string, metadata map[string]interface{}) error {
	newStruct, err := structpb.NewStruct(metadata)
	if err != nil {
		return fmt.Errorf("failure to wrap metadata %w", err)
	}

	if r.schemaV1Resp != nil {
		// nolint:staticcheck
		r.schemaV1Resp.Msg.Value = value
		r.schemaV1Resp.Msg.Variant = variant
		r.schemaV1Resp.Msg.Reason = reason
		r.schemaV1Resp.Msg.Metadata = newStruct
	}
	if r.evalV1Resp != nil {
		r.evalV1Resp.Msg.Value = value
		r.evalV1Resp.Msg.Variant = variant
		r.evalV1Resp.Msg.Reason = reason
		r.evalV1Resp.Msg.Metadata = newStruct
	}

	return nil
}

type intResponse struct {
	//nolint:staticcheck
	schemaV1Resp *connect.Response[schemaV1.ResolveIntResponse]
	evalV1Resp   *connect.Response[evalV1.ResolveIntResponse]
}

//nolint:staticcheck
func (r *intResponse) SetResult(value int64, variant, reason string, metadata map[string]interface{}) error {
	newStruct, err := structpb.NewStruct(metadata)
	if err != nil {
		return fmt.Errorf("failure to wrap metadata %w", err)
	}

	if r.schemaV1Resp != nil {
		r.schemaV1Resp.Msg.Value = value
		r.schemaV1Resp.Msg.Variant = variant
		r.schemaV1Resp.Msg.Reason = reason
		r.schemaV1Resp.Msg.Metadata = newStruct
	}
	if r.evalV1Resp != nil {
		r.evalV1Resp.Msg.Value = value
		r.evalV1Resp.Msg.Variant = variant
		r.evalV1Resp.Msg.Reason = reason
		r.evalV1Resp.Msg.Metadata = newStruct
	}
	return nil
}

type objectResponse struct {
	// nolint:staticcheck
	schemaV1Resp *connect.Response[schemaV1.ResolveObjectResponse]
	evalV1Resp   *connect.Response[evalV1.ResolveObjectResponse]
}

//nolint:staticcheck
func (r *objectResponse) SetResult(value map[string]any, variant, reason string,
	metadata map[string]interface{},
) error {
	newStruct, err := structpb.NewStruct(metadata)
	if err != nil {
		return fmt.Errorf("failure to wrap metadata %w", err)
	}
	if r.schemaV1Resp != nil {
		r.schemaV1Resp.Msg.Reason = reason
		val, err := structpb.NewStruct(value)
		if err != nil {
			return fmt.Errorf("struct response construction: %w", err)
		}

		r.schemaV1Resp.Msg.Value = val
		r.schemaV1Resp.Msg.Variant = variant
		r.schemaV1Resp.Msg.Metadata = newStruct
	}
	if r.evalV1Resp != nil {
		r.evalV1Resp.Msg.Reason = reason
		val, err := structpb.NewStruct(value)
		if err != nil {
			return fmt.Errorf("struct response construction: %w", err)
		}

		r.evalV1Resp.Msg.Value = val
		r.evalV1Resp.Msg.Variant = variant
		r.evalV1Resp.Msg.Metadata = newStruct
	}
	return nil
}

// V2 response types with optional value and variant

type responseV2[T constraints] interface {
	SetResult(value T, variant, reason string, metadata map[string]interface{}) error
	SetReasonOnly(reason string, metadata map[string]interface{}) error
}

type booleanResponseV2 struct {
	evalV2Resp *connect.Response[evalV2.ResolveBooleanResponse]
}

func (r *booleanResponseV2) SetResult(value bool, variant, reason string, metadata map[string]interface{}) error {
	newStruct, err := structpb.NewStruct(metadata)
	if err != nil {
		return fmt.Errorf("failure to wrap metadata %w", err)
	}

	if r.evalV2Resp != nil {
		r.evalV2Resp.Msg.Value = &value
		if variant != "" {
			r.evalV2Resp.Msg.Variant = &variant
		}
		r.evalV2Resp.Msg.Reason = reason
		r.evalV2Resp.Msg.Metadata = newStruct
	}

	return nil
}

func (r *booleanResponseV2) SetReasonOnly(reason string, metadata map[string]interface{}) error {
	newStruct, err := structpb.NewStruct(metadata)
	if err != nil {
		return fmt.Errorf("failure to wrap metadata %w", err)
	}

	if r.evalV2Resp != nil {
		// Leave Value and Variant as nil (unset)
		r.evalV2Resp.Msg.Reason = reason
		r.evalV2Resp.Msg.Metadata = newStruct
	}

	return nil
}

type stringResponseV2 struct {
	evalV2Resp *connect.Response[evalV2.ResolveStringResponse]
}

func (r *stringResponseV2) SetResult(value string, variant, reason string, metadata map[string]interface{}) error {
	newStruct, err := structpb.NewStruct(metadata)
	if err != nil {
		return fmt.Errorf("failure to wrap metadata %w", err)
	}

	if r.evalV2Resp != nil {
		r.evalV2Resp.Msg.Value = &value
		if variant != "" {
			r.evalV2Resp.Msg.Variant = &variant
		}
		r.evalV2Resp.Msg.Reason = reason
		r.evalV2Resp.Msg.Metadata = newStruct
	}

	return nil
}

func (r *stringResponseV2) SetReasonOnly(reason string, metadata map[string]interface{}) error {
	newStruct, err := structpb.NewStruct(metadata)
	if err != nil {
		return fmt.Errorf("failure to wrap metadata %w", err)
	}

	if r.evalV2Resp != nil {
		// Leave Value and Variant as nil (unset)
		r.evalV2Resp.Msg.Reason = reason
		r.evalV2Resp.Msg.Metadata = newStruct
	}

	return nil
}

type floatResponseV2 struct {
	evalV2Resp *connect.Response[evalV2.ResolveFloatResponse]
}

func (r *floatResponseV2) SetResult(value float64, variant, reason string, metadata map[string]interface{}) error {
	newStruct, err := structpb.NewStruct(metadata)
	if err != nil {
		return fmt.Errorf("failure to wrap metadata %w", err)
	}

	if r.evalV2Resp != nil {
		r.evalV2Resp.Msg.Value = &value
		if variant != "" {
			r.evalV2Resp.Msg.Variant = &variant
		}
		r.evalV2Resp.Msg.Reason = reason
		r.evalV2Resp.Msg.Metadata = newStruct
	}

	return nil
}

func (r *floatResponseV2) SetReasonOnly(reason string, metadata map[string]interface{}) error {
	newStruct, err := structpb.NewStruct(metadata)
	if err != nil {
		return fmt.Errorf("failure to wrap metadata %w", err)
	}

	if r.evalV2Resp != nil {
		// Leave Value and Variant as nil (unset)
		r.evalV2Resp.Msg.Reason = reason
		r.evalV2Resp.Msg.Metadata = newStruct
	}

	return nil
}

type intResponseV2 struct {
	evalV2Resp *connect.Response[evalV2.ResolveIntResponse]
}

func (r *intResponseV2) SetResult(value int64, variant, reason string, metadata map[string]interface{}) error {
	newStruct, err := structpb.NewStruct(metadata)
	if err != nil {
		return fmt.Errorf("failure to wrap metadata %w", err)
	}

	if r.evalV2Resp != nil {
		r.evalV2Resp.Msg.Value = &value
		if variant != "" {
			r.evalV2Resp.Msg.Variant = &variant
		}
		r.evalV2Resp.Msg.Reason = reason
		r.evalV2Resp.Msg.Metadata = newStruct
	}
	return nil
}

func (r *intResponseV2) SetReasonOnly(reason string, metadata map[string]interface{}) error {
	newStruct, err := structpb.NewStruct(metadata)
	if err != nil {
		return fmt.Errorf("failure to wrap metadata %w", err)
	}

	if r.evalV2Resp != nil {
		// Leave Value and Variant as nil (unset)
		r.evalV2Resp.Msg.Reason = reason
		r.evalV2Resp.Msg.Metadata = newStruct
	}
	return nil
}

type objectResponseV2 struct {
	evalV2Resp *connect.Response[evalV2.ResolveObjectResponse]
}

func (r *objectResponseV2) SetResult(value map[string]any, variant, reason string,
	metadata map[string]interface{},
) error {
	newStruct, err := structpb.NewStruct(metadata)
	if err != nil {
		return fmt.Errorf("failure to wrap metadata %w", err)
	}
	if r.evalV2Resp != nil {
		r.evalV2Resp.Msg.Reason = reason
		val, err := structpb.NewStruct(value)
		if err != nil {
			return fmt.Errorf("struct response construction: %w", err)
		}

		r.evalV2Resp.Msg.Value = val
		if variant != "" {
			r.evalV2Resp.Msg.Variant = &variant
		}
		r.evalV2Resp.Msg.Metadata = newStruct
	}
	return nil
}
