package strategy_test

import (
	eb "github.com/amahdavian/async-event-bus"
	"github.com/amahdavian/async-event-bus/strategy"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"testing"
	"time"
)

var assertThat = then.AssertThat
var equals = is.EqualTo

func TestShouldReturnGivenTimeout(t *testing.T) {
	timeout := 100 * time.Millisecond

	logging := strategy.NewLogging(timeout, &mockLogger{})

	assertThat(t, logging.GetTimeout(), equals(timeout))
}

func TestShouldLogGivenEventOnDeliveryFailure(t *testing.T) {
	timeout := 100 * time.Millisecond
	logger := &mockLogger{}
	loggingStrategy := strategy.NewLogging(timeout, logger)

	failedEvent := eb.Event{
		Name:    "SomeEvent",
		Details: "SomeDetails",
	}
	loggingStrategy.OnDeliveryFailure(failedEvent)

	assertThat(t, logger.Message, equals("failed to deliver message: %v"))
	assertThat(t, logger.Parameters, equals([]interface{}{failedEvent}))
}

type mockLogger struct {
	Message string
	Parameters []interface{}
}

func (m *mockLogger) Infof(message string, event ...interface{}) {
	m.Message = message
	m.Parameters = event
}
