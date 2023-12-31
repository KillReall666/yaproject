// Code generated by mockery v2.35.4. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// DBStatusChecker is an autogenerated mock type for the DBStatusChecker type
type DBStatusChecker struct {
	mock.Mock
}

// DBStatusCheck provides a mock function with given fields:
func (_m *DBStatusChecker) DBStatusCheck() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewDBStatusChecker creates a new instance of DBStatusChecker. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDBStatusChecker(t interface {
	mock.TestingT
	Cleanup(func())
}) *DBStatusChecker {
	mock := &DBStatusChecker{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
