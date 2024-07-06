package timing

import "time"

type Clock interface {
	Now() time.Time
}

type Action func()

func (s Action) perform() { s() }

type EventScheduler interface {
	Clock
	DoAfter(time.Duration, Action)
	DoRepeatedly(time.Duration, Action)
}
