package eventstub

import (
	"sync"

	"microblogging/internal/model"
)

type PubSubManager struct {
	subscribers map[string][]chan model.Event
	mu          sync.RWMutex
}

func NewPubSub() *PubSubManager {
	return &PubSubManager{
		subscribers: make(map[string][]chan model.Event),
	}
}

func (p *PubSubManager) Publish(topic string, message model.Event) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if chans, found := p.subscribers[topic]; found {
		for _, ch := range chans {
			go func(c chan model.Event) { c <- message }(ch)
		}
	}
}

func (p *PubSubManager) Subscribe(topic string, handler func(model.Event)) {
	p.mu.Lock()
	defer p.mu.Unlock()

	ch := make(chan model.Event, 10)
	p.subscribers[topic] = append(p.subscribers[topic], ch)

	// listen for messages
	go func() {
		for msg := range ch {
			handler(msg)
		}
	}()
}
