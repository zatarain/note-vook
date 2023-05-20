// Code generated by mockery v2.27.1. DO NOT EDIT.

package mocks

import (
	http "net/http"

	gin "github.com/gin-gonic/gin"

	mock "github.com/stretchr/testify/mock"
)

// MockedEngine is an autogenerated mock type for the MockedEngine type
type MockedEngine struct {
	mock.Mock
}

// Any provides a mock function with given fields: _a0, _a1
func (_m *MockedEngine) Any(_a0 string, _a1 ...gin.HandlerFunc) gin.IRoutes {
	_va := make([]interface{}, len(_a1))
	for _i := range _a1 {
		_va[_i] = _a1[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 gin.IRoutes
	if rf, ok := ret.Get(0).(func(string, ...gin.HandlerFunc) gin.IRoutes); ok {
		r0 = rf(_a0, _a1...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gin.IRoutes)
		}
	}

	return r0
}

// DELETE provides a mock function with given fields: _a0, _a1
func (_m *MockedEngine) DELETE(_a0 string, _a1 ...gin.HandlerFunc) gin.IRoutes {
	_va := make([]interface{}, len(_a1))
	for _i := range _a1 {
		_va[_i] = _a1[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 gin.IRoutes
	if rf, ok := ret.Get(0).(func(string, ...gin.HandlerFunc) gin.IRoutes); ok {
		r0 = rf(_a0, _a1...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gin.IRoutes)
		}
	}

	return r0
}

// GET provides a mock function with given fields: _a0, _a1
func (_m *MockedEngine) GET(_a0 string, _a1 ...gin.HandlerFunc) gin.IRoutes {
	_va := make([]interface{}, len(_a1))
	for _i := range _a1 {
		_va[_i] = _a1[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 gin.IRoutes
	if rf, ok := ret.Get(0).(func(string, ...gin.HandlerFunc) gin.IRoutes); ok {
		r0 = rf(_a0, _a1...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gin.IRoutes)
		}
	}

	return r0
}

// Group provides a mock function with given fields: _a0, _a1
func (_m *MockedEngine) Group(_a0 string, _a1 ...gin.HandlerFunc) *gin.RouterGroup {
	_va := make([]interface{}, len(_a1))
	for _i := range _a1 {
		_va[_i] = _a1[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *gin.RouterGroup
	if rf, ok := ret.Get(0).(func(string, ...gin.HandlerFunc) *gin.RouterGroup); ok {
		r0 = rf(_a0, _a1...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gin.RouterGroup)
		}
	}

	return r0
}

// HEAD provides a mock function with given fields: _a0, _a1
func (_m *MockedEngine) HEAD(_a0 string, _a1 ...gin.HandlerFunc) gin.IRoutes {
	_va := make([]interface{}, len(_a1))
	for _i := range _a1 {
		_va[_i] = _a1[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 gin.IRoutes
	if rf, ok := ret.Get(0).(func(string, ...gin.HandlerFunc) gin.IRoutes); ok {
		r0 = rf(_a0, _a1...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gin.IRoutes)
		}
	}

	return r0
}

// Handle provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockedEngine) Handle(_a0 string, _a1 string, _a2 ...gin.HandlerFunc) gin.IRoutes {
	_va := make([]interface{}, len(_a2))
	for _i := range _a2 {
		_va[_i] = _a2[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0, _a1)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 gin.IRoutes
	if rf, ok := ret.Get(0).(func(string, string, ...gin.HandlerFunc) gin.IRoutes); ok {
		r0 = rf(_a0, _a1, _a2...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gin.IRoutes)
		}
	}

	return r0
}

// Match provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockedEngine) Match(_a0 []string, _a1 string, _a2 ...gin.HandlerFunc) gin.IRoutes {
	_va := make([]interface{}, len(_a2))
	for _i := range _a2 {
		_va[_i] = _a2[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0, _a1)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 gin.IRoutes
	if rf, ok := ret.Get(0).(func([]string, string, ...gin.HandlerFunc) gin.IRoutes); ok {
		r0 = rf(_a0, _a1, _a2...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gin.IRoutes)
		}
	}

	return r0
}

// OPTIONS provides a mock function with given fields: _a0, _a1
func (_m *MockedEngine) OPTIONS(_a0 string, _a1 ...gin.HandlerFunc) gin.IRoutes {
	_va := make([]interface{}, len(_a1))
	for _i := range _a1 {
		_va[_i] = _a1[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 gin.IRoutes
	if rf, ok := ret.Get(0).(func(string, ...gin.HandlerFunc) gin.IRoutes); ok {
		r0 = rf(_a0, _a1...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gin.IRoutes)
		}
	}

	return r0
}

// PATCH provides a mock function with given fields: _a0, _a1
func (_m *MockedEngine) PATCH(_a0 string, _a1 ...gin.HandlerFunc) gin.IRoutes {
	_va := make([]interface{}, len(_a1))
	for _i := range _a1 {
		_va[_i] = _a1[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 gin.IRoutes
	if rf, ok := ret.Get(0).(func(string, ...gin.HandlerFunc) gin.IRoutes); ok {
		r0 = rf(_a0, _a1...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gin.IRoutes)
		}
	}

	return r0
}

// POST provides a mock function with given fields: _a0, _a1
func (_m *MockedEngine) POST(_a0 string, _a1 ...gin.HandlerFunc) gin.IRoutes {
	_va := make([]interface{}, len(_a1))
	for _i := range _a1 {
		_va[_i] = _a1[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 gin.IRoutes
	if rf, ok := ret.Get(0).(func(string, ...gin.HandlerFunc) gin.IRoutes); ok {
		r0 = rf(_a0, _a1...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gin.IRoutes)
		}
	}

	return r0
}

// PUT provides a mock function with given fields: _a0, _a1
func (_m *MockedEngine) PUT(_a0 string, _a1 ...gin.HandlerFunc) gin.IRoutes {
	_va := make([]interface{}, len(_a1))
	for _i := range _a1 {
		_va[_i] = _a1[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 gin.IRoutes
	if rf, ok := ret.Get(0).(func(string, ...gin.HandlerFunc) gin.IRoutes); ok {
		r0 = rf(_a0, _a1...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gin.IRoutes)
		}
	}

	return r0
}

// Static provides a mock function with given fields: _a0, _a1
func (_m *MockedEngine) Static(_a0 string, _a1 string) gin.IRoutes {
	ret := _m.Called(_a0, _a1)

	var r0 gin.IRoutes
	if rf, ok := ret.Get(0).(func(string, string) gin.IRoutes); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gin.IRoutes)
		}
	}

	return r0
}

// StaticFS provides a mock function with given fields: _a0, _a1
func (_m *MockedEngine) StaticFS(_a0 string, _a1 http.FileSystem) gin.IRoutes {
	ret := _m.Called(_a0, _a1)

	var r0 gin.IRoutes
	if rf, ok := ret.Get(0).(func(string, http.FileSystem) gin.IRoutes); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gin.IRoutes)
		}
	}

	return r0
}

// StaticFile provides a mock function with given fields: _a0, _a1
func (_m *MockedEngine) StaticFile(_a0 string, _a1 string) gin.IRoutes {
	ret := _m.Called(_a0, _a1)

	var r0 gin.IRoutes
	if rf, ok := ret.Get(0).(func(string, string) gin.IRoutes); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gin.IRoutes)
		}
	}

	return r0
}

// StaticFileFS provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockedEngine) StaticFileFS(_a0 string, _a1 string, _a2 http.FileSystem) gin.IRoutes {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 gin.IRoutes
	if rf, ok := ret.Get(0).(func(string, string, http.FileSystem) gin.IRoutes); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gin.IRoutes)
		}
	}

	return r0
}

// Use provides a mock function with given fields: _a0
func (_m *MockedEngine) Use(_a0 ...gin.HandlerFunc) gin.IRoutes {
	_va := make([]interface{}, len(_a0))
	for _i := range _a0 {
		_va[_i] = _a0[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 gin.IRoutes
	if rf, ok := ret.Get(0).(func(...gin.HandlerFunc) gin.IRoutes); ok {
		r0 = rf(_a0...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gin.IRoutes)
		}
	}

	return r0
}

type mockConstructorTestingTNewMockedEngine interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockedEngine creates a new instance of MockedEngine. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockedEngine(t mockConstructorTestingTNewMockedEngine) *MockedEngine {
	mock := &MockedEngine{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}