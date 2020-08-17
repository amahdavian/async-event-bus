package messaging

import (
	"github.com/amahdavian/async-event-bus/internal/concurrency"
)

type Event struct {
	Name    string
	Details interface{}
}

// EventChannel is a channel which can accept an Event
type EventChannel chan Event

type EventBus struct {
	subscribers *concurrency.MapOfSlice
}

func NewEventBus() *EventBus {
	return &EventBus{subscribers: concurrency.NewMapOfSlice()}
}

func (eventBus *EventBus) Publish(event Event) {
	if eventSubscribers, found := eventBus.subscribers.Get(event.Name); found {
		go func(event Event, eventChannels []interface{}) {
			for _, ch := range eventChannels {
				ch.(EventChannel) <- event
			}
		}(event, eventSubscribers)
	}
}

func (eventBus *EventBus) Subscribe(eventName string) EventChannel {
	channel := make(EventChannel)
	eventBus.subscribers.AppendAt(eventName, channel)
	return channel
}

func (eventBus *EventBus) UnSubscribe(eventName string, channel EventChannel) {
	eventBus.subscribers.RemoveAt(eventName, channel)
}
