// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "context"

import mock "github.com/stretchr/testify/mock"
import time "time"

// Recorder is an autogenerated mock type for the Recorder type
type Recorder struct {
	mock.Mock
}

// AddInflightRequests provides a mock function with given fields: ctx, id, quantity
func (_m *Recorder) AddInflightRequests(ctx context.Context, id string, quantity int) {
	_m.Called(ctx, id, quantity)
}

// ObserveHTTPRequestDuration provides a mock function with given fields: ctx, id, duration, method, code
func (_m *Recorder) ObserveHTTPRequestDuration(ctx context.Context, id string, duration time.Duration, method string, code string) {
	_m.Called(ctx, id, duration, method, code)
}

// ObserveHTTPResponseSize provides a mock function with given fields: ctx, id, responseSize, method, code
func (_m *Recorder) ObserveHTTPResponseSize(ctx context.Context, id string, responseSize int64, method string, code string) {
	_m.Called(ctx, id, responseSize, method, code)
}
