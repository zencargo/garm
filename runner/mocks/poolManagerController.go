// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"
	common "garm/runner/common"

	databasecommon "garm/database/common"

	mock "github.com/stretchr/testify/mock"

	params "garm/params"
)

// PoolManagerController is an autogenerated mock type for the PoolManagerController type
type PoolManagerController struct {
	mock.Mock
}

// CreateOrgPoolManager provides a mock function with given fields: ctx, org, providers, store
func (_m *PoolManagerController) CreateOrgPoolManager(ctx context.Context, org params.Organization, providers map[string]common.Provider, store databasecommon.Store) (common.PoolManager, error) {
	ret := _m.Called(ctx, org, providers, store)

	var r0 common.PoolManager
	if rf, ok := ret.Get(0).(func(context.Context, params.Organization, map[string]common.Provider, databasecommon.Store) common.PoolManager); ok {
		r0 = rf(ctx, org, providers, store)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.PoolManager)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, params.Organization, map[string]common.Provider, databasecommon.Store) error); ok {
		r1 = rf(ctx, org, providers, store)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateRepoPoolManager provides a mock function with given fields: ctx, repo, providers, store
func (_m *PoolManagerController) CreateRepoPoolManager(ctx context.Context, repo params.Repository, providers map[string]common.Provider, store databasecommon.Store) (common.PoolManager, error) {
	ret := _m.Called(ctx, repo, providers, store)

	var r0 common.PoolManager
	if rf, ok := ret.Get(0).(func(context.Context, params.Repository, map[string]common.Provider, databasecommon.Store) common.PoolManager); ok {
		r0 = rf(ctx, repo, providers, store)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.PoolManager)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, params.Repository, map[string]common.Provider, databasecommon.Store) error); ok {
		r1 = rf(ctx, repo, providers, store)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteOrgPoolManager provides a mock function with given fields: org
func (_m *PoolManagerController) DeleteOrgPoolManager(org params.Organization) error {
	ret := _m.Called(org)

	var r0 error
	if rf, ok := ret.Get(0).(func(params.Organization) error); ok {
		r0 = rf(org)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteRepoPoolManager provides a mock function with given fields: repo
func (_m *PoolManagerController) DeleteRepoPoolManager(repo params.Repository) error {
	ret := _m.Called(repo)

	var r0 error
	if rf, ok := ret.Get(0).(func(params.Repository) error); ok {
		r0 = rf(repo)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetOrgPoolManager provides a mock function with given fields: org
func (_m *PoolManagerController) GetOrgPoolManager(org params.Organization) (common.PoolManager, error) {
	ret := _m.Called(org)

	var r0 common.PoolManager
	if rf, ok := ret.Get(0).(func(params.Organization) common.PoolManager); ok {
		r0 = rf(org)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.PoolManager)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(params.Organization) error); ok {
		r1 = rf(org)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOrgPoolManagers provides a mock function with given fields:
func (_m *PoolManagerController) GetOrgPoolManagers() (map[string]common.PoolManager, error) {
	ret := _m.Called()

	var r0 map[string]common.PoolManager
	if rf, ok := ret.Get(0).(func() map[string]common.PoolManager); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]common.PoolManager)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRepoPoolManager provides a mock function with given fields: repo
func (_m *PoolManagerController) GetRepoPoolManager(repo params.Repository) (common.PoolManager, error) {
	ret := _m.Called(repo)

	var r0 common.PoolManager
	if rf, ok := ret.Get(0).(func(params.Repository) common.PoolManager); ok {
		r0 = rf(repo)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.PoolManager)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(params.Repository) error); ok {
		r1 = rf(repo)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRepoPoolManagers provides a mock function with given fields:
func (_m *PoolManagerController) GetRepoPoolManagers() (map[string]common.PoolManager, error) {
	ret := _m.Called()

	var r0 map[string]common.PoolManager
	if rf, ok := ret.Get(0).(func() map[string]common.PoolManager); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]common.PoolManager)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewPoolManagerController interface {
	mock.TestingT
	Cleanup(func())
}

// NewPoolManagerController creates a new instance of PoolManagerController. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPoolManagerController(t mockConstructorTestingTNewPoolManagerController) *PoolManagerController {
	mock := &PoolManagerController{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
