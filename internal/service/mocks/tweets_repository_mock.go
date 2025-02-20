package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"microblogging/internal/model"
)

type TweetsRepositoryMock struct {
	mock.Mock
}

func (_m *TweetsRepositoryMock) CreateTweet(_a0 context.Context, _a1 model.Tweet) (model.Tweet, error) {
	ret := _m.Called(_a0, _a1)

	var r0 model.Tweet
	if rf, ok := ret.Get(0).(func(context.Context, model.Tweet) model.Tweet); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(model.Tweet)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.Tweet) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *TweetsRepositoryMock) GetByFilters(_a0 map[string]interface{}, _a1 int, _a2 int, _a3 string) ([]model.Tweet, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 []model.Tweet
	if rf, ok := ret.Get(0).(func(map[string]interface{}, int, int, string) []model.Tweet); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Get(0).([]model.Tweet)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(map[string]interface{}, int, int, string) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func NewTweetsRepositoryMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *TweetsRepositoryMock {
	mock := &TweetsRepositoryMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
