package strategy

import (
	eb "github.com/amahdavian/async-event-bus"
	"time"
)

type Logger interface {
	Infof(string, ...interface{})
}

func NewLogging(timeout time.Duration, logger Logger) *Logging {
	return &Logging{timeout: timeout, logger: logger}
}

type Logging struct {
	timeout time.Duration
	logger Logger
}

func (n Logging) GetTimeout() time.Duration {
	return n.timeout
}

func (n Logging) OnDeliveryFailure(event eb.Event) {
	n.logger.Infof("failed to deliver message: %v", event)
}
