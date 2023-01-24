package grpc

import "testing"

func TestUrlToGRPCTarget(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UrlToGRPCTarget(tt.args.url); got != tt.want {
				t.Errorf("UrlToGRPCTarget() = %v, want %v", got, tt.want)
			}
		})
	}
}
