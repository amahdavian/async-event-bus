package strategy

import (
	eb "github.com/amahdavian/async-event-bus/v2"
	"time"
)

type logger interface {
	Infof(string, ...interface{})
}

func NewLogging(timeout time.Duration, logger logger) *logging {
	return &logging{timeout: timeout, logger: logger}
}

type logging struct {
	timeout time.Duration
	logger  logger
}

func (n logging) GetTimeout() time.Duration {
	return n.timeout
}

func (n logging) OnDeliveryFailure(event eb.Event) {
	n.logger.Infof("failed to deliver message: %v", event)
}
