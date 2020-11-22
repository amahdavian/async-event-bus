package strategy

import (
	eb "github.com/amahdavian/async-event-bus/v2"
	"time"
)

type publisher interface {
	Publish(event eb.Event, strategy eb.BackoffStrategy)
}

func NewRetry(publisher publisher, timeout time.Duration, numberOfRetries int) *retry {
	return &retry{
		publisher:       publisher,
		timeout:         timeout,
		numberOfRetries: numberOfRetries,
	}
}

type retry struct {
	publisher       publisher
	timeout         time.Duration
	numberOfRetries int
}

func (r retry) GetTimeout() time.Duration {
	return r.timeout
}

func (r retry) OnDeliveryFailure(event eb.Event) {
	if r.numberOfRetries > 0 {
		r.publisher.Publish(event, NewRetry(r.publisher, r.timeout, r.numberOfRetries-1))
	}
}
