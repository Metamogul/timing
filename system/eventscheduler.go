package system

import (
	"context"
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

func (e *EventScheduler) PerformAfter(action timing.Action, duration time.Duration, ctx context.Context) {
	go func() {
		select {
		case <-time.After(duration):
			action.Perform(e.Clock)
		case <-ctx.Done():
			return
		}
	}()
}

func (e *EventScheduler) PerformRepeatedly(action timing.Action, until *time.Time, interval time.Duration, ctx context.Context) {
	ticker := time.NewTicker(interval)

	var timer *time.Timer
	if until != nil {
		timer = time.NewTimer(until.Sub(e.Now()))
	} else {
		timer = &time.Timer{}
	}

	go func() {
		for {
			select {
			case <-ticker.C:
				action.Perform(Clock{})
			case <-timer.C:
				return
			case <-ctx.Done():
				return
			}
		}
	}()
}
