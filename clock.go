package timing

import "time"

type Clock interface {
	Now() time.Time
}

type EventScheduler interface {
	Clock
	DoAfter(time.Duration, func())
	DoRepeatedly(time.Duration, func())
}
