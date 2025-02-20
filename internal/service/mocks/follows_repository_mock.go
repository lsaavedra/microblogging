package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"microblogging/internal/model"
)

type FollowsRepositoryMock struct {
	mock.Mock
}

func (_m *FollowsRepositoryMock) CreateFollowing(_a0 context.Context, _a1 model.Follow) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.Follow) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *FollowsRepositoryMock) GetFollowers(_a0 context.Context, _a1 uuid.UUID) ([]model.Follow, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []model.Follow
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) []model.Follow); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).([]model.Follow)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func NewFollowsRepositoryMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *FollowsRepositoryMock {
	mock := &FollowsRepositoryMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
