package timing

import (
	"context"
	"time"
)

type Clock interface {
	Now() time.Time
}

type Action interface {
	Perform(Clock)
}

type EventScheduler interface {
	Clock
	PerformNow(action Action, ctx context.Context)
	PerformAfter(action Action, duration time.Duration, ctx context.Context)
	PerformRepeatedly(action Action, until *time.Time, interval time.Duration, ctx context.Context)
}
