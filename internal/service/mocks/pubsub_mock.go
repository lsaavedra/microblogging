package mocks

import (
	"github.com/stretchr/testify/mock"

	"microblogging/internal/model"
)

type PubSubMock struct {
	mock.Mock
}

func (_m *PubSubMock) Publish(_a0 string, _a1 model.Event) {
	//Add logic
}

func (_m *PubSubMock) Subscribe(topic string, handler func(event model.Event)) {
	//Add logic
}

func NewPubSubMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *PubSubMock {
	mock := &PubSubMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
