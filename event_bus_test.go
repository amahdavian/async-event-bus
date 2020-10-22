package messaging_test

import (
	"testing"
	"time"

	eb "github.com/amahdavian/async-event-bus/v2"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
)

var assertThat = then.AssertThat
var equals = is.EqualTo

func TestShouldAllowConcurrentSubscriptionsAndPublishEventToThem(t *testing.T) {
	strategy := NewMockBackoffStrategy(100 * time.Millisecond)

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

	eventBus.Publish(expectedEvent, strategy)

	for i := 0; i < 100; i++ {
		select {
		case <-assertionCompletedChannel:
		case <-time.After(2 * time.Second):
			t.Errorf("did not receive assertion result within 2 second")
		}
	}
	assertThat(t, len(strategy.failedDeliveries), equals(0))
}

func TestShouldAllowUnsubscribingFromEvents(t *testing.T) {
	strategy := NewMockBackoffStrategy(100 * time.Millisecond)

	expectedEvent := eb.Event{Name: "SomeEvent", Details: "SomeDetails"}

	eventBus := eb.NewEventBus()

	channel := eventBus.Subscribe("SomeEvent")

	eventBus.Publish(expectedEvent, strategy)

	select {
	case actualEvent := <-channel:
		assertThat(t, actualEvent, equals(expectedEvent))
	case <-time.After(1 * time.Second):
		t.Errorf("did not receive the event within 1 second")
	}

	eventBus.UnSubscribe("SomeEvent", channel)

	eventBus.Publish(expectedEvent, strategy)

	select {
	case <-channel:
		t.Errorf("received the event after unsubscribing")
	case <-time.After(2 * time.Second):
	}
	assertThat(t, len(strategy.failedDeliveries), equals(0))
}

func TestShouldBackoffUsingGivenBackoffStrategyWhenSubscriberChannelIsFull(t *testing.T) {
	strategy := NewMockBackoffStrategy(100 * time.Millisecond)

	expectedEvent := eb.Event{Name: "SomeEvent", Details: "SomeDetails"}

	eventBus := eb.NewEventBus()

	_ = eventBus.Subscribe("SomeEvent")

	eventBus.Publish(expectedEvent, strategy)

	select {
	case <-time.After(200 * time.Millisecond):
		assertThat(t, len(strategy.failedDeliveries), equals(1))
		assertThat(t, strategy.failedDeliveries[0], equals(expectedEvent))
	}
}

func NewMockBackoffStrategy(timeout time.Duration) *mockBackoffStrategy {
	return &mockBackoffStrategy{
		failedDeliveries: []eb.Event{},
		timeout:          timeout,
	}
}

type mockBackoffStrategy struct {
	failedDeliveries []eb.Event
	timeout          time.Duration
}

func (m mockBackoffStrategy) GetTimeout() time.Duration {
	return m.timeout
}

func (m *mockBackoffStrategy) OnDeliveryFailure(event eb.Event) {
	m.failedDeliveries = append(m.failedDeliveries, event)
}
