// Code generated by mockery v1.0.1. DO NOT EDIT.

package mocks

import (
	context "context"

	admin "github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"

	mock "github.com/stretchr/testify/mock"

	service "github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/service"
)

// AdminFetcherExtInterface is an autogenerated mock type for the AdminFetcherExtInterface type
type AdminFetcherExtInterface struct {
	mock.Mock
}

type AdminFetcherExtInterface_AdminServiceClient struct {
	*mock.Call
}

func (_m AdminFetcherExtInterface_AdminServiceClient) Return(_a0 service.AdminServiceClient) *AdminFetcherExtInterface_AdminServiceClient {
	return &AdminFetcherExtInterface_AdminServiceClient{Call: _m.Call.Return(_a0)}
}

func (_m *AdminFetcherExtInterface) OnAdminServiceClient() *AdminFetcherExtInterface_AdminServiceClient {
	c := _m.On("AdminServiceClient")
	return &AdminFetcherExtInterface_AdminServiceClient{Call: c}
}

func (_m *AdminFetcherExtInterface) OnAdminServiceClientMatch(matchers ...interface{}) *AdminFetcherExtInterface_AdminServiceClient {
	c := _m.On("AdminServiceClient", matchers...)
	return &AdminFetcherExtInterface_AdminServiceClient{Call: c}
}

// AdminServiceClient provides a mock function with given fields:
func (_m *AdminFetcherExtInterface) AdminServiceClient() service.AdminServiceClient {
	ret := _m.Called()

	var r0 service.AdminServiceClient
	if rf, ok := ret.Get(0).(func() service.AdminServiceClient); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.AdminServiceClient)
		}
	}

	return r0
}

type AdminFetcherExtInterface_FetchAllVerOfLP struct {
	*mock.Call
}

func (_m AdminFetcherExtInterface_FetchAllVerOfLP) Return(_a0 []*admin.LaunchPlan, _a1 error) *AdminFetcherExtInterface_FetchAllVerOfLP {
	return &AdminFetcherExtInterface_FetchAllVerOfLP{Call: _m.Call.Return(_a0, _a1)}
}

func (_m *AdminFetcherExtInterface) OnFetchAllVerOfLP(ctx context.Context, lpName string, project string, domain string) *AdminFetcherExtInterface_FetchAllVerOfLP {
	c := _m.On("FetchAllVerOfLP", ctx, lpName, project, domain)
	return &AdminFetcherExtInterface_FetchAllVerOfLP{Call: c}
}

func (_m *AdminFetcherExtInterface) OnFetchAllVerOfLPMatch(matchers ...interface{}) *AdminFetcherExtInterface_FetchAllVerOfLP {
	c := _m.On("FetchAllVerOfLP", matchers...)
	return &AdminFetcherExtInterface_FetchAllVerOfLP{Call: c}
}

// FetchAllVerOfLP provides a mock function with given fields: ctx, lpName, project, domain
func (_m *AdminFetcherExtInterface) FetchAllVerOfLP(ctx context.Context, lpName string, project string, domain string) ([]*admin.LaunchPlan, error) {
	ret := _m.Called(ctx, lpName, project, domain)

	var r0 []*admin.LaunchPlan
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) []*admin.LaunchPlan); ok {
		r0 = rf(ctx, lpName, project, domain)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*admin.LaunchPlan)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string) error); ok {
		r1 = rf(ctx, lpName, project, domain)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type AdminFetcherExtInterface_FetchAllVerOfTask struct {
	*mock.Call
}

func (_m AdminFetcherExtInterface_FetchAllVerOfTask) Return(_a0 []*admin.Task, _a1 error) *AdminFetcherExtInterface_FetchAllVerOfTask {
	return &AdminFetcherExtInterface_FetchAllVerOfTask{Call: _m.Call.Return(_a0, _a1)}
}

func (_m *AdminFetcherExtInterface) OnFetchAllVerOfTask(ctx context.Context, name string, project string, domain string) *AdminFetcherExtInterface_FetchAllVerOfTask {
	c := _m.On("FetchAllVerOfTask", ctx, name, project, domain)
	return &AdminFetcherExtInterface_FetchAllVerOfTask{Call: c}
}

func (_m *AdminFetcherExtInterface) OnFetchAllVerOfTaskMatch(matchers ...interface{}) *AdminFetcherExtInterface_FetchAllVerOfTask {
	c := _m.On("FetchAllVerOfTask", matchers...)
	return &AdminFetcherExtInterface_FetchAllVerOfTask{Call: c}
}

// FetchAllVerOfTask provides a mock function with given fields: ctx, name, project, domain
func (_m *AdminFetcherExtInterface) FetchAllVerOfTask(ctx context.Context, name string, project string, domain string) ([]*admin.Task, error) {
	ret := _m.Called(ctx, name, project, domain)

	var r0 []*admin.Task
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) []*admin.Task); ok {
		r0 = rf(ctx, name, project, domain)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*admin.Task)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string) error); ok {
		r1 = rf(ctx, name, project, domain)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type AdminFetcherExtInterface_FetchExecution struct {
	*mock.Call
}

func (_m AdminFetcherExtInterface_FetchExecution) Return(_a0 *admin.Execution, _a1 error) *AdminFetcherExtInterface_FetchExecution {
	return &AdminFetcherExtInterface_FetchExecution{Call: _m.Call.Return(_a0, _a1)}
}

func (_m *AdminFetcherExtInterface) OnFetchExecution(ctx context.Context, name string, project string, domain string) *AdminFetcherExtInterface_FetchExecution {
	c := _m.On("FetchExecution", ctx, name, project, domain)
	return &AdminFetcherExtInterface_FetchExecution{Call: c}
}

func (_m *AdminFetcherExtInterface) OnFetchExecutionMatch(matchers ...interface{}) *AdminFetcherExtInterface_FetchExecution {
	c := _m.On("FetchExecution", matchers...)
	return &AdminFetcherExtInterface_FetchExecution{Call: c}
}

// FetchExecution provides a mock function with given fields: ctx, name, project, domain
func (_m *AdminFetcherExtInterface) FetchExecution(ctx context.Context, name string, project string, domain string) (*admin.Execution, error) {
	ret := _m.Called(ctx, name, project, domain)

	var r0 *admin.Execution
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) *admin.Execution); ok {
		r0 = rf(ctx, name, project, domain)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*admin.Execution)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string) error); ok {
		r1 = rf(ctx, name, project, domain)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type AdminFetcherExtInterface_FetchLPLatestVersion struct {
	*mock.Call
}

func (_m AdminFetcherExtInterface_FetchLPLatestVersion) Return(_a0 *admin.LaunchPlan, _a1 error) *AdminFetcherExtInterface_FetchLPLatestVersion {
	return &AdminFetcherExtInterface_FetchLPLatestVersion{Call: _m.Call.Return(_a0, _a1)}
}

func (_m *AdminFetcherExtInterface) OnFetchLPLatestVersion(ctx context.Context, name string, project string, domain string) *AdminFetcherExtInterface_FetchLPLatestVersion {
	c := _m.On("FetchLPLatestVersion", ctx, name, project, domain)
	return &AdminFetcherExtInterface_FetchLPLatestVersion{Call: c}
}

func (_m *AdminFetcherExtInterface) OnFetchLPLatestVersionMatch(matchers ...interface{}) *AdminFetcherExtInterface_FetchLPLatestVersion {
	c := _m.On("FetchLPLatestVersion", matchers...)
	return &AdminFetcherExtInterface_FetchLPLatestVersion{Call: c}
}

// FetchLPLatestVersion provides a mock function with given fields: ctx, name, project, domain
func (_m *AdminFetcherExtInterface) FetchLPLatestVersion(ctx context.Context, name string, project string, domain string) (*admin.LaunchPlan, error) {
	ret := _m.Called(ctx, name, project, domain)

	var r0 *admin.LaunchPlan
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) *admin.LaunchPlan); ok {
		r0 = rf(ctx, name, project, domain)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*admin.LaunchPlan)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string) error); ok {
		r1 = rf(ctx, name, project, domain)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type AdminFetcherExtInterface_FetchLPVersion struct {
	*mock.Call
}

func (_m AdminFetcherExtInterface_FetchLPVersion) Return(_a0 *admin.LaunchPlan, _a1 error) *AdminFetcherExtInterface_FetchLPVersion {
	return &AdminFetcherExtInterface_FetchLPVersion{Call: _m.Call.Return(_a0, _a1)}
}

func (_m *AdminFetcherExtInterface) OnFetchLPVersion(ctx context.Context, name string, version string, project string, domain string) *AdminFetcherExtInterface_FetchLPVersion {
	c := _m.On("FetchLPVersion", ctx, name, version, project, domain)
	return &AdminFetcherExtInterface_FetchLPVersion{Call: c}
}

func (_m *AdminFetcherExtInterface) OnFetchLPVersionMatch(matchers ...interface{}) *AdminFetcherExtInterface_FetchLPVersion {
	c := _m.On("FetchLPVersion", matchers...)
	return &AdminFetcherExtInterface_FetchLPVersion{Call: c}
}

// FetchLPVersion provides a mock function with given fields: ctx, name, version, project, domain
func (_m *AdminFetcherExtInterface) FetchLPVersion(ctx context.Context, name string, version string, project string, domain string) (*admin.LaunchPlan, error) {
	ret := _m.Called(ctx, name, version, project, domain)

	var r0 *admin.LaunchPlan
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, string) *admin.LaunchPlan); ok {
		r0 = rf(ctx, name, version, project, domain)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*admin.LaunchPlan)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string, string) error); ok {
		r1 = rf(ctx, name, version, project, domain)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type AdminFetcherExtInterface_FetchTaskLatestVersion struct {
	*mock.Call
}

func (_m AdminFetcherExtInterface_FetchTaskLatestVersion) Return(_a0 *admin.Task, _a1 error) *AdminFetcherExtInterface_FetchTaskLatestVersion {
	return &AdminFetcherExtInterface_FetchTaskLatestVersion{Call: _m.Call.Return(_a0, _a1)}
}

func (_m *AdminFetcherExtInterface) OnFetchTaskLatestVersion(ctx context.Context, name string, project string, domain string) *AdminFetcherExtInterface_FetchTaskLatestVersion {
	c := _m.On("FetchTaskLatestVersion", ctx, name, project, domain)
	return &AdminFetcherExtInterface_FetchTaskLatestVersion{Call: c}
}

func (_m *AdminFetcherExtInterface) OnFetchTaskLatestVersionMatch(matchers ...interface{}) *AdminFetcherExtInterface_FetchTaskLatestVersion {
	c := _m.On("FetchTaskLatestVersion", matchers...)
	return &AdminFetcherExtInterface_FetchTaskLatestVersion{Call: c}
}

// FetchTaskLatestVersion provides a mock function with given fields: ctx, name, project, domain
func (_m *AdminFetcherExtInterface) FetchTaskLatestVersion(ctx context.Context, name string, project string, domain string) (*admin.Task, error) {
	ret := _m.Called(ctx, name, project, domain)

	var r0 *admin.Task
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) *admin.Task); ok {
		r0 = rf(ctx, name, project, domain)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*admin.Task)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string) error); ok {
		r1 = rf(ctx, name, project, domain)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type AdminFetcherExtInterface_FetchTaskVersion struct {
	*mock.Call
}

func (_m AdminFetcherExtInterface_FetchTaskVersion) Return(_a0 *admin.Task, _a1 error) *AdminFetcherExtInterface_FetchTaskVersion {
	return &AdminFetcherExtInterface_FetchTaskVersion{Call: _m.Call.Return(_a0, _a1)}
}

func (_m *AdminFetcherExtInterface) OnFetchTaskVersion(ctx context.Context, name string, version string, project string, domain string) *AdminFetcherExtInterface_FetchTaskVersion {
	c := _m.On("FetchTaskVersion", ctx, name, version, project, domain)
	return &AdminFetcherExtInterface_FetchTaskVersion{Call: c}
}

func (_m *AdminFetcherExtInterface) OnFetchTaskVersionMatch(matchers ...interface{}) *AdminFetcherExtInterface_FetchTaskVersion {
	c := _m.On("FetchTaskVersion", matchers...)
	return &AdminFetcherExtInterface_FetchTaskVersion{Call: c}
}

// FetchTaskVersion provides a mock function with given fields: ctx, name, version, project, domain
func (_m *AdminFetcherExtInterface) FetchTaskVersion(ctx context.Context, name string, version string, project string, domain string) (*admin.Task, error) {
	ret := _m.Called(ctx, name, version, project, domain)

	var r0 *admin.Task
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, string) *admin.Task); ok {
		r0 = rf(ctx, name, version, project, domain)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*admin.Task)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string, string) error); ok {
		r1 = rf(ctx, name, version, project, domain)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
