// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	"github.com/gol4ng/httpware/v3/auth"
	"github.com/stretchr/testify/mock"
)

// Authenticator is an autogenerated mock type for the Authenticator type
type Authenticator struct {
	mock.Mock
}

// Authenticate provides a mock function with given fields: _a0
func (_m *Authenticator) Authenticate(_a0 auth.Credential) (auth.Credential, error) {
	ret := _m.Called(_a0)

	var r0 auth.Credential
	if rf, ok := ret.Get(0).(func(auth.Credential) auth.Credential); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(auth.Credential)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(auth.Credential) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}