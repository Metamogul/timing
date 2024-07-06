package system

import (
	"time"
)

type Clock struct{}

func (s *Clock) Now() time.Time {
	return time.Now()
}

type EventScheduler struct{}

func (s *EventScheduler) DoAfter(duration time.Duration, f func()) {
	time.AfterFunc(duration, f)
}

func (s *EventScheduler) DoRepeatedly(duration time.Duration, f func()) {
	ticker := time.NewTicker(500 * time.Millisecond)

	go func() {
		for {
			<-ticker.C
			f()
		}
	}()
}
