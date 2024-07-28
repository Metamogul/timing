package timing

import "time"

type Clock interface {
	Now() time.Time
}

type Action interface {
	Perform(Clock)
}

type EventScheduler interface {
	Clock
	PerformAfter(action Action, duration time.Duration)
	PerformRepeatedly(action Action, until *time.Time, interval time.Duration)
}
