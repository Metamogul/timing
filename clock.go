package timing

import "time"

type Clock interface {
	Now() time.Time
}

type Action func()

func (a Action) perform() { a() }

type EventScheduler interface {
	Clock
	PerformAfter(time.Duration, Action)
	PerformRepeatedly(time.Duration, Action)
}
