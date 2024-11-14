// Code generated by MockGen. DO NOT EDIT.
// Source: .\internal\api\publisher.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIPublisher is a mock of IPublisher interface.
type MockIPublisher struct {
	ctrl     *gomock.Controller
	recorder *MockIPublisherMockRecorder
}

// MockIPublisherMockRecorder is the mock recorder for MockIPublisher.
type MockIPublisherMockRecorder struct {
	mock *MockIPublisher
}

// NewMockIPublisher creates a new mock instance.
func NewMockIPublisher(ctrl *gomock.Controller) *MockIPublisher {
	mock := &MockIPublisher{ctrl: ctrl}
	mock.recorder = &MockIPublisherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIPublisher) EXPECT() *MockIPublisherMockRecorder {
	return m.recorder
}

// Publish mocks base method.
func (m *MockIPublisher) Publish(topic string, event interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", topic, event)
	ret0, _ := ret[0].(error)
	return ret0
}

// Publish indicates an expected call of Publish.
func (mr *MockIPublisherMockRecorder) Publish(topic, event interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockIPublisher)(nil).Publish), topic, event)
}