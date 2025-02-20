package mock

import (
	"fmt"

	"microblogging/internal/model"
)

type PubSubMock struct{}

func (p *PubSubMock) Publish(topic string, message model.Event) {
	fmt.Sprintf("publishing to topic:%v - message:%v", topic, message)
}

func (p *PubSubMock) Subscribe(topic string, handler func(model.Event)) {
	go func() {
		fmt.Sprintf("suscribing to topic:%v", topic)
	}()
}
