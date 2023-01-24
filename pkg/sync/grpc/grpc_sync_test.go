package grpc

import "testing"

func TestUrlToGRPCTarget(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "With Prefix",
			url:  "grpc://test.com/endpoint",
			want: "test.com/endpoint",
		},
		{
			name: "Without Prefix",
			url:  "test.com/endpoint",
			want: "test.com/endpoint",
		},
		{
			name: "Empty is empty",
			url:  "",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := URLToGRPCTarget(tt.url); got != tt.want {
				t.Errorf("URLToGRPCTarget() = %v, want %v", got, tt.want)
			}
		})
	}
}
