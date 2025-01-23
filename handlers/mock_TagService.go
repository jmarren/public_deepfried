// Code generated by mockery v2.51.1. DO NOT EDIT.

package handlers

import mock "github.com/stretchr/testify/mock"

// MockTagService is an autogenerated mock type for the TagService type
type MockTagService struct {
	mock.Mock
}

type MockTagService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockTagService) EXPECT() *MockTagService_Expecter {
	return &MockTagService_Expecter{mock: &_m.Mock}
}

// GetMostPopular provides a mock function with no fields
func (_m *MockTagService) GetMostPopular() []string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetMostPopular")
	}

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// MockTagService_GetMostPopular_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetMostPopular'
type MockTagService_GetMostPopular_Call struct {
	*mock.Call
}

// GetMostPopular is a helper method to define mock.On call
func (_e *MockTagService_Expecter) GetMostPopular() *MockTagService_GetMostPopular_Call {
	return &MockTagService_GetMostPopular_Call{Call: _e.mock.On("GetMostPopular")}
}

func (_c *MockTagService_GetMostPopular_Call) Run(run func()) *MockTagService_GetMostPopular_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockTagService_GetMostPopular_Call) Return(_a0 []string) *MockTagService_GetMostPopular_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTagService_GetMostPopular_Call) RunAndReturn(run func() []string) *MockTagService_GetMostPopular_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockTagService creates a new instance of MockTagService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockTagService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockTagService {
	mock := &MockTagService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
