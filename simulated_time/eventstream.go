package simulated_time

import (
	"errors"
	"time"
)

var ErrEventStreamClosed = errors.New("event stream is closed")

type eventStream interface {
	popNextEvent() *event
	peakNextTime() time.Time

	closed() bool
}
