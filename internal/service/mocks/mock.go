// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockGetLink is a mock of GetLink interface.
type MockGetLink struct {
	ctrl     *gomock.Controller
	recorder *MockGetLinkMockRecorder
}

// MockGetLinkMockRecorder is the mock recorder for MockGetLink.
type MockGetLinkMockRecorder struct {
	mock *MockGetLink
}

// NewMockGetLink creates a new mock instance.
func NewMockGetLink(ctrl *gomock.Controller) *MockGetLink {
	mock := &MockGetLink{ctrl: ctrl}
	mock.recorder = &MockGetLinkMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGetLink) EXPECT() *MockGetLinkMockRecorder {
	return m.recorder
}

// GetLink mocks base method.
func (m *MockGetLink) GetLink(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLink", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLink indicates an expected call of GetLink.
func (mr *MockGetLinkMockRecorder) GetLink(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLink", reflect.TypeOf((*MockGetLink)(nil).GetLink), arg0)
}

// MockCreateLink is a mock of CreateLink interface.
type MockCreateLink struct {
	ctrl     *gomock.Controller
	recorder *MockCreateLinkMockRecorder
}

// MockCreateLinkMockRecorder is the mock recorder for MockCreateLink.
type MockCreateLinkMockRecorder struct {
	mock *MockCreateLink
}

// NewMockCreateLink creates a new mock instance.
func NewMockCreateLink(ctrl *gomock.Controller) *MockCreateLink {
	mock := &MockCreateLink{ctrl: ctrl}
	mock.recorder = &MockCreateLinkMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCreateLink) EXPECT() *MockCreateLinkMockRecorder {
	return m.recorder
}

// CreateLink mocks base method.
func (m *MockCreateLink) CreateLink(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateLink", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateLink indicates an expected call of CreateLink.
func (mr *MockCreateLinkMockRecorder) CreateLink(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateLink", reflect.TypeOf((*MockCreateLink)(nil).CreateLink), arg0)
}