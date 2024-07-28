package system

import (
	"github.com/metamogul/timing"
	"time"
)

type Clock struct{}

func (s Clock) Now() time.Time {
	return time.Now()
}

type EventScheduler struct {
	Clock
}

func (s *EventScheduler) PerformAfter(duration time.Duration, action timing.Action) {
	time.AfterFunc(duration, func() {
		action.Perform(Clock{})
	})
}

func (s *EventScheduler) PerformRepeatedly(duration time.Duration, action timing.Action) {
	ticker := time.NewTicker(duration)

	go func() {
		for {
			<-ticker.C
			action.Perform(Clock{})
		}
	}()
}
