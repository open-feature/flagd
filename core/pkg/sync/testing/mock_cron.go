package testing

import (
	"reflect"

	"go.uber.org/mock/gomock"
)

// MockCron is a mock of Cron interface.
type MockCron struct {
	ctrl     *gomock.Controller
	recorder *MockCronMockRecorder
	cmd      func()
}

// MockCronMockRecorder is the mock recorder for MockCron.
type MockCronMockRecorder struct {
	mock *MockCron
}

// NewMockCron creates a new mock instance.
func NewMockCron(ctrl *gomock.Controller) *MockCron {
	mock := &MockCron{ctrl: ctrl}
	mock.recorder = &MockCronMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCron) EXPECT() *MockCronMockRecorder {
	return m.recorder
}

// AddFunc mocks base method.
func (m *MockCron) AddFunc(spec string, cmd func()) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddFunc", spec, cmd)
	ret0, _ := ret[0].(error)
	m.cmd = cmd
	return ret0
}

func (m *MockCron) Tick() {
	m.cmd()
}

// AddFunc indicates an expected call of AddFunc.
func (mr *MockCronMockRecorder) AddFunc(spec, cmd any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddFunc", reflect.TypeOf((*MockCron)(nil).AddFunc), spec, cmd)
}

// Start mocks base method.
func (m *MockCron) Start() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Start")
}

// Start indicates an expected call of Start.
func (mr *MockCronMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockCron)(nil).Start))
}

// Stop mocks base method.
func (m *MockCron) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop.
func (mr *MockCronMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockCron)(nil).Stop))
}
