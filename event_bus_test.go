package messaging_test

import (
	"testing"
	"time"

	eb "github.com/amahdavian/async-event-bus"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
)

var assertThat = then.AssertThat
var equals = is.EqualTo
var contains = is.ValueContaining

func TestShouldAllowConcurrentSubscriptionsAndPublishEventToThem(t *testing.T) {
	expectedEvent := eb.Event{Name: "SomeEvent", Details: "SomeDetails"}

	eventBus := eb.NewEventBus()

	subscriptionCompletedChannel := make(chan bool)
	assertionCompletedChannel := make(chan bool)
	for i := 0; i < 100; i++ {
		go func() {
			channel := eventBus.Subscribe("SomeEvent")
			subscriptionCompletedChannel <- true

			go func() {
				select {
				case actualEvent := <-channel:
					assertThat(t, actualEvent, equals(expectedEvent))
					assertionCompletedChannel <- true
				case <-time.After(1 * time.Second):
					t.Errorf("did not receive the event within 1 second")
					assertionCompletedChannel <- true
				}
			}()
		}()
	}

	for i := 0; i < 100; i++ {
		select {
		case <-subscriptionCompletedChannel:
		case <-time.After(1 * time.Second):
			t.Errorf("did not complete subscribing within 1 second")
		}
	}

	eventBus.Publish(expectedEvent)

	for i := 0; i < 100; i++ {
		select {
		case <-assertionCompletedChannel:
		case <-time.After(2 * time.Second):
			t.Errorf("did not receive assertion result within 2 second")
		}
	}
}

func TestShouldAllowUnsubscribingFromEvents(t *testing.T) {
	expectedEvent := eb.Event{Name: "SomeEvent", Details: "SomeDetails"}

	eventBus := eb.NewEventBus()

	channel := eventBus.Subscribe("SomeEvent")

	eventBus.Publish(expectedEvent)

	select {
	case actualEvent := <-channel:
		assertThat(t, actualEvent, equals(expectedEvent))
	case <-time.After(1 * time.Second):
		t.Errorf("did not receive the event within 1 second")
	}

	eventBus.UnSubscribe("SomeEvent", channel)

	eventBus.Publish(expectedEvent)

	select {
	case <-channel:
		t.Errorf("received the event after unsubscribing")
	case <-time.After(2 * time.Second):
	}
}
