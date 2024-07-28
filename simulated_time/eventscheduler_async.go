package simulated_time

import (
	"github.com/metamogul/timing"
	"sync"
	"time"
)

type AsyncEventScheduler struct {
	*clock

	eventGenerators   *eventCombinator
	eventGeneratorsMu sync.RWMutex

	wg sync.WaitGroup
}

func NewAsyncEventScheduler(now time.Time) *SerialEventScheduler {
	return &SerialEventScheduler{
		clock:           newClock(now),
		eventGenerators: newEventCombinator(),
	}
}

func (e *AsyncEventScheduler) Forward(interval time.Duration) {
	targetTime := e.clock.Now().Add(interval)

	for e.performNextEvent(targetTime) {
	}

	e.wg.Wait()
}

func (e *AsyncEventScheduler) performNextEvent(targetTime time.Time) (shouldContinue bool) {
	e.eventGeneratorsMu.RLock()
	defer e.eventGeneratorsMu.RUnlock()

	if e.eventGenerators.finished() {
		e.clock.set(targetTime)
		return false
	}

	if e.eventGenerators.peek().After(targetTime) {
		e.clock.set(targetTime)
		return false
	}

	nextEvent := e.eventGenerators.pop()
	e.clock.set(nextEvent.Time)

	currentClock := e.clock.copy()
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		nextEvent.Perform(currentClock)
	}()

	return true
}

func (e *AsyncEventScheduler) PerformAfter(action timing.Action, interval time.Duration) {
	e.eventGeneratorsMu.Lock()
	defer e.eventGeneratorsMu.Unlock()

	e.eventGenerators.addInput(newSingleEventGenerator(action, e.now.Add(interval)))
}

func (a *AsyncEventScheduler) PerformRepeatedly(action timing.Action, until *time.Time, interval time.Duration) {
	a.eventGeneratorsMu.Lock()
	defer a.eventGeneratorsMu.Unlock()

	a.eventGenerators.addInput(newPeriodicEventGenerator(action, a.Now(), until, interval))
}
