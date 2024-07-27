package simulated_time

import (
	"github.com/metamogul/timing"
	"sync"
	"time"
)

type EventScheduler struct {
	*clock

	eventGenerators   *eventCombinator
	eventGeneratorsMu sync.RWMutex

	wg sync.WaitGroup
}

func NewEventScheduler(now time.Time) *EventScheduler {
	return &EventScheduler{
		clock:           newClock(now),
		eventGenerators: newEventCombinator(),
	}
}

func (e *EventScheduler) Forward(interval time.Duration) {
	targetTime := e.clock.Now().Add(interval)

	for e.dispatchNextEvent(targetTime) {
	}

	e.wg.Wait()
}

func (e *EventScheduler) dispatchNextEvent(targetTime time.Time) (shouldContinue bool) {
	e.eventGeneratorsMu.RLock()
	defer e.eventGeneratorsMu.RUnlock()

	if e.eventGenerators.finished() {
		e.clock.Set(targetTime)
		return false
	}

	if e.eventGenerators.peek().After(targetTime) {
		e.clock.Set(targetTime)
		return false
	}

	nextEvent := e.eventGenerators.pop()
	e.clock.Set(nextEvent.Time)

	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		nextEvent.Perform()
	}()

	return true
}

func (e *EventScheduler) PerformAfter(action timing.Action, interval time.Duration) {
	e.eventGeneratorsMu.Lock()
	defer e.eventGeneratorsMu.Unlock()

	e.eventGenerators.addInput(newSingleEventGenerator(action, e.now.Add(interval)))
}

func (e *EventScheduler) PerformRepeatedly(action timing.Action, until *time.Time, interval time.Duration) {
	e.eventGeneratorsMu.Lock()
	defer e.eventGeneratorsMu.Unlock()

	e.eventGenerators.addInput(newPeriodicEventGenerator(action, e.Now(), until, interval))
}
