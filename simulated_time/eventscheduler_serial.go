package simulated_time

import (
	"github.com/metamogul/timing"
	"sync"
	"time"
)

type SerialEventScheduler struct {
	*clock

	eventGenerators   *eventCombinator
	eventGeneratorsMu sync.RWMutex
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

func (a *SerialEventScheduler) PerformAfter(action timing.Action, interval time.Duration) {
	a.eventGenerators.addInput(newSingleEventGenerator(action, a.now.Add(interval)))
}

func (a *SerialEventScheduler) PerformRepeatedly(action timing.Action, until *time.Time, interval time.Duration) {
	a.eventGenerators.addInput(newPeriodicEventGenerator(action, a.Now(), until, interval))
}
