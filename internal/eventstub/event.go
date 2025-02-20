package eventstub

import "microblogging/internal/model"

type EventProcessor interface {
	Publish(topic string, message model.Event)
	Subscribe(topic string, handler func(event model.Event))
}
