package simulated_time

import (
	"context"
	"github.com/metamogul/timing"
	"time"
)

type SerialEventScheduler struct {
	*clock

	eventGenerators *eventCombinator
}

func NewSerialEventScheduler(now time.Time) *SerialEventScheduler {
	return &SerialEventScheduler{
		clock:           newClock(now),
		eventGenerators: newEventCombinator(),
	}
}

func (a *SerialEventScheduler) Forward(interval time.Duration) {
	targetTime := a.clock.Now().Add(interval)

	for a.performNextEvent(targetTime) {
	}
}

func (a *SerialEventScheduler) performNextEvent(targetTime time.Time) (shouldContinue bool) {
	if a.eventGenerators.finished() {
		a.clock.set(targetTime)
		return false
	}

	if a.eventGenerators.peek().After(targetTime) {
		a.clock.set(targetTime)
		return false
	}

	nextEvent := a.eventGenerators.pop()
	a.clock.set(nextEvent.Time)

	nextEvent.Perform(a.clock.copy())

	return true
}

func (a *SerialEventScheduler) ForwardToNextEvent() {
	if a.eventGenerators.finished() {
		return
	}

	nextEvent := a.eventGenerators.pop()
	a.clock.set(nextEvent.Time)

	nextEvent.Perform(a.clock.copy())
}

func (a *SerialEventScheduler) PerformAfter(action timing.Action, interval time.Duration, ctx context.Context) {
	a.eventGenerators.add(newSingleEventGenerator(action, a.now.Add(interval), ctx))
}

func (a *SerialEventScheduler) PerformRepeatedly(action timing.Action, until *time.Time, interval time.Duration, ctx context.Context) {
	a.eventGenerators.add(newPeriodicEventGenerator(action, a.Now(), until, interval, ctx))
}
