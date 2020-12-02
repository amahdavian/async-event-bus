package messaging

import (
	"github.com/amahdavian/async-event-bus/v2/internal/concurrency"
	"time"
)

type Event struct {
	Name    string
	Details interface{}
}

// EventChannel is a channel which can accept an Event
type EventChannel chan Event

// BackoffStrategy specifies what should happen to the message if the intended subscriber/s is not able to receive the message within the timeout period.
type BackoffStrategy interface {
	// GetTimeout specifies the timeout period.
	// A timeout of 0 will lead to non-blocking channel operations as in the message will only be sent if the subscriber is ready to receive, otherwise, onDeliveryFailure will be called.
	GetTimeout() time.Duration
	// OnDeliveryFailure is called when the subscriber has been unable to receive the message within the timeout period.
	OnDeliveryFailure(Event)
}

type eventBus struct {
	subscribers *concurrency.MapOfSlice
}

// NewEventBus creates a new instance of eventBus
func NewEventBus() *eventBus {
	return &eventBus{subscribers: concurrency.NewMapOfSlice()}
}

// Publish publishes the given event to all subscribers and upon delivery failure, uses backoff strategy to recover.
func (eventBus *eventBus) Publish(event Event, backoffStrategy BackoffStrategy) {
	if eventSubscribers, found := eventBus.subscribers.Get(event.Name); found {
		for _, subscriber := range eventSubscribers {
			go func(event Event, eventChannel interface{}) {
				select {
				case eventChannel.(EventChannel) <- event:
				case <-time.After(backoffStrategy.GetTimeout()):
					backoffStrategy.OnDeliveryFailure(event)
				}
			}(event, subscriber)
		}
	}
}

// Subscribe creates a subscription to the specific topic/event name.
func (eventBus *eventBus) Subscribe(eventName string) EventChannel {
	channel := make(EventChannel)
	eventBus.subscribers.AppendAt(eventName, channel)
	return channel
}

// UnSubscribe terminates the subscription of the given channel to the given topic/event name.
func (eventBus *eventBus) UnSubscribe(eventName string, channel EventChannel) {
	eventBus.subscribers.RemoveAt(eventName, channel)
}
