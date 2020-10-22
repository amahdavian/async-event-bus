package strategy

import (
	eb "github.com/amahdavian/async-event-bus/v2"
	"time"
)

type Publisher interface {
	Publish(event eb.Event, strategy eb.BackoffStrategy)
}

func NewRetry(publisher Publisher, timeout time.Duration, numberOfRetries int) *Retry {
	return &Retry{
		publisher:       publisher,
		timeout:         timeout,
		numberOfRetries: numberOfRetries,
	}
}

type Retry struct {
	publisher       Publisher
	timeout         time.Duration
	numberOfRetries int
}

func (r Retry) GetTimeout() time.Duration {
	return r.timeout
}

func (r Retry) OnDeliveryFailure(event eb.Event) {
	if r.numberOfRetries > 0 {
		r.publisher.Publish(event, NewRetry(r.publisher, r.timeout, r.numberOfRetries-1))
	}
}
