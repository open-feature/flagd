// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/sync/grpc/credentials/builder.go

// Package credendialsmock is a generated GoMock package.
package credendialsmock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	credentials "google.golang.org/grpc/credentials"
)

// MockBuilder is a mock of Builder interface.
type MockBuilder struct {
	ctrl     *gomock.Controller
	recorder *MockBuilderMockRecorder
}

// MockBuilderMockRecorder is the mock recorder for MockBuilder.
type MockBuilderMockRecorder struct {
	mock *MockBuilder
}

// NewMockBuilder creates a new mock instance.
func NewMockBuilder(ctrl *gomock.Controller) *MockBuilder {
	mock := &MockBuilder{ctrl: ctrl}
	mock.recorder = &MockBuilderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBuilder) EXPECT() *MockBuilderMockRecorder {
	return m.recorder
}

// Build mocks base method.
func (m *MockBuilder) Build(secure bool, certPath string) (credentials.TransportCredentials, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Build", secure, certPath)
	ret0, _ := ret[0].(credentials.TransportCredentials)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Build indicates an expected call of Build.
func (mr *MockBuilderMockRecorder) Build(secure, certPath interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Build", reflect.TypeOf((*MockBuilder)(nil).Build), secure, certPath)
}
