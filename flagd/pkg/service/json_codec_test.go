//nolint:wrapcheck
package service

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestJSONCodec_Marshaling(t *testing.T) {
	originalMap := map[string]interface{}{
		"val1": false,
		"val2": 0.0,
		"val3": "",
		"val4": true,
		"val5": 1.0,
		"val6": "hi!",
	}
	tests := map[string]struct {
		message    func() (any, error)
		marshaller func(any) ([]byte, error)
		wantErr    bool
		wantMap    map[string]interface{}
	}{
		"Marshal no error": {
			message: func() (any, error) {
				return structpb.NewStruct(originalMap)
			},
			marshaller: func(message any) ([]byte, error) {
				jsonCodec := jsonCodec{}
				return jsonCodec.Marshal(message)
			},
			wantErr: false,
			wantMap: originalMap,
		},
		"MarshalStable no error": {
			message: func() (any, error) {
				return structpb.NewStruct(originalMap)
			},
			marshaller: func(message any) ([]byte, error) {
				jsonCodec := jsonCodec{}
				return jsonCodec.MarshalStable(message)
			},
			wantErr: false,
			wantMap: originalMap,
		},
		"MarshalAppend no error": {
			message: func() (any, error) {
				return structpb.NewStruct(originalMap)
			},
			marshaller: func(message any) ([]byte, error) {
				jsonCodec := jsonCodec{}
				return jsonCodec.MarshalAppend(nil, message)
			},
			wantErr: false,
			wantMap: originalMap,
		},
		"Marshal not valid message": {
			message: func() (any, error) {
				return map[string]interface{}{}, nil
			},
			marshaller: func(message any) ([]byte, error) {
				jsonCodec := jsonCodec{}
				return jsonCodec.Marshal(nil)
			},
			wantErr: true,
			wantMap: map[string]interface{}{},
		},
		"MarshalStable not valid message": {
			message: func() (any, error) {
				return map[string]interface{}{}, nil
			},
			marshaller: func(message any) ([]byte, error) {
				jsonCodec := jsonCodec{}
				return jsonCodec.MarshalStable(nil)
			},
			wantErr: true,
			wantMap: map[string]interface{}{},
		},
		"MarshalAppend not valid message": {
			message: func() (any, error) {
				return map[string]interface{}{}, nil
			},
			marshaller: func(message any) ([]byte, error) {
				jsonCodec := jsonCodec{}
				return jsonCodec.MarshalAppend(nil, nil)
			},
			wantErr: true,
			wantMap: map[string]interface{}{},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			jsonMap := map[string]interface{}{}
			message, err := tt.message()
			if err != nil {
				t.Errorf("Got error creating message: %v", err.Error())
			}
			bytes, err := tt.marshaller(message)
			if tt.wantErr {
				require.NotNilf(t, err, "Expected error but got none")
			}
			//nolint:errcheck
			json.Unmarshal(bytes, &jsonMap)
			require.Equal(t, tt.wantMap, jsonMap)
		})
	}
}

func TestJSONCodec_Unmarshal(t *testing.T) {
	originalMap := map[string]interface{}{
		"val1": false,
		"val2": 0.0,
		"val3": "",
		"val4": true,
		"val5": 1.0,
		"val6": "hi!",
	}

	bytes, err := json.Marshal(originalMap)
	if err != nil {
		t.Errorf("Got error marshalling json: %v", err.Error())
	}

	jsonCodec := jsonCodec{}

	message, err := structpb.NewStruct(map[string]interface{}{})
	if err != nil {
		t.Errorf("Got error creating struct: %v", err.Error())
	}

	//nolint:errcheck
	jsonCodec.Unmarshal(bytes, message)

	require.Equal(t, originalMap, message.AsMap())
}
