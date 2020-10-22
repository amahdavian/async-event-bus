package strategy_test

import (
	eb "github.com/amahdavian/async-event-bus/v2"
	"github.com/amahdavian/async-event-bus/v2/strategy"
	"testing"
	"time"
)

func TestRetryStrategyShouldReturnGivenTimeout(t *testing.T) {
	timeout := 100 * time.Millisecond

	logging := strategy.NewRetry(&mockPublisher{}, timeout, 0)

	assertThat(t, logging.GetTimeout(), equals(timeout))
}

func TestShouldRetryGivenEventOnDeliveryFailureUpToRetryLimit(t *testing.T) {
	timeout := 100 * time.Millisecond
	numberOfRetries := 5
	publisher := newMockPublisher(2 * numberOfRetries)
	retryStrategy := strategy.NewRetry(publisher, timeout, numberOfRetries)

	failedEvent := eb.Event{
		Name:    "SomeEvent",
		Details: "SomeDetails",
	}
	retryStrategy.OnDeliveryFailure(failedEvent)

	assertThat(t, publisher.events, equals([]eb.Event{failedEvent, failedEvent, failedEvent, failedEvent, failedEvent}))
	assertThat(t, len(publisher.strategies), equals(numberOfRetries))
}

func TestShouldRetryGivenEventOnDeliveryFailureOnlyTilSuccessful(t *testing.T) {
	timeout := 100 * time.Millisecond
	numberOfRetries := 5
	publisher := newMockPublisher(1)
	retryStrategy := strategy.NewRetry(publisher, timeout, numberOfRetries)

	failedEvent := eb.Event{
		Name:    "SomeEvent",
		Details: "SomeDetails",
	}
	retryStrategy.OnDeliveryFailure(failedEvent)

	assertThat(t, publisher.events, equals([]eb.Event{failedEvent, failedEvent}))
	assertThat(t, len(publisher.strategies), equals(2))
	assertThat(t, publisher.strategies[1].(*strategy.Retry), equals(strategy.NewRetry(publisher, timeout, 3)))
}

func newMockPublisher(numberOfTimesToFail int) *mockPublisher {
	return &mockPublisher{numberOfTimesToFail: numberOfTimesToFail}
}

type mockPublisher struct {
	events              []eb.Event
	strategies          []eb.BackoffStrategy
	numberOfTimesToFail int
}

func (m *mockPublisher) Publish(event eb.Event, strategy eb.BackoffStrategy) {
	m.events = append(m.events, event)
	m.strategies = append(m.strategies, strategy)
	if m.numberOfTimesToFail > 0 {
		m.numberOfTimesToFail--
		strategy.OnDeliveryFailure(event)
	}
}
