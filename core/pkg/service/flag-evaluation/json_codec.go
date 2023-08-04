//nolint:wrapcheck
package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bufbuild/connect-go"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoiface"
)

// WithJSON customizes a connect-go Client or Handler's JSON by exposing MarshalOptions, and UnmarshalOptions
// See: https://github.com/bufbuild/connect-go/blob/main/codec.go
// Heavily inspired by https://github.com/akshayjshah/connectproto
func WithJSON(marshal protojson.MarshalOptions, unmarshal protojson.UnmarshalOptions) connect.Option {
	return connect.WithOptions(
		// mark the codec with the correct content-type
		connect.WithCodec(&jsonCodec{name: "json", marshal: marshal, unmarshal: unmarshal}),
		connect.WithCodec(&jsonCodec{name: "json; charset=utf-8", marshal: marshal, unmarshal: unmarshal}),
	)
}

type jsonCodec struct {
	name      string
	marshal   protojson.MarshalOptions
	unmarshal protojson.UnmarshalOptions
}
var _ connect.Codec = (*jsonCodec)(nil)
func (j *jsonCodec) Name() string {
	return j.name
}

func (j *jsonCodec) IsBinary() bool {
	return false
}

// Marshal marshals the given proto.Message in the JSON format.
// Do not depend on the output being stable.
// It may change over time across different versions of the program.
func (j *jsonCodec) Marshal(message any) ([]byte, error) {
	return j.MarshalAppend(nil, message)
}

// MarshalAppend appends the JSON format encoding of message to dst, returning the result.
func (j *jsonCodec) MarshalAppend(dst []byte, message any) ([]byte, error) {
	protoMessage, ok := message.(proto.Message)
	if !ok {
		return nil, errNotMessage(message)
	}
	return j.marshal.MarshalAppend(dst, protoMessage)
}

// Unmarshal reads the given []byte into the given proto.Message.
// The provided message must be mutable (e.g., a non-nil pointer to a message).
func (j *jsonCodec) Unmarshal(binary []byte, message any) error {
	protoMessage, ok := message.(proto.Message)
	if !ok {
		return errNotMessage(message)
	}
	if len(binary) == 0 {
		return errors.New("zero-length payload is not a valid JSON object")
	}
	return j.unmarshal.Unmarshal(binary, protoMessage)
}

// Marshal marshals the given proto.Message in the JSON format.
// It attempts to offer a determinist output by removing inconsistent whitespace.
func (j *jsonCodec) MarshalStable(message any) ([]byte, error) {
	// protojson doesn't offer deterministic output. It does order fields by
	// number, but it deliberately introduce inconsistent whitespace (see
	// https://github.com/golang/protobuf/issues/1373). To make the output as
	// consistent as possible, we'll need to normalize.
	messageJSON, err := j.Marshal(message)
	if err != nil {
		return nil, err
	}
	compacted := bytes.NewBuffer(messageJSON[:0])
	if err = json.Compact(compacted, messageJSON); err != nil {
		return nil, err
	}
	return compacted.Bytes(), nil
}

func errNotMessage(msg any) error {
	if _, ok := msg.(protoiface.MessageV1); ok {
		return fmt.Errorf(
			"%T uses github.com/golang/protobuf, not google.golang.org/protobuf: see https://go.dev/blog/protobuf-apiv2",
			msg)
	}
	return fmt.Errorf("%T doesn't implement proto.Message", msg)
}
