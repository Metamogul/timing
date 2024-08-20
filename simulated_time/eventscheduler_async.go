package simulated_time

import (
	"context"
	"sync"
	"time"

	"github.com/metamogul/timing"
)

type AsyncEventScheduler struct {
	*clock

	eventGenerators   *eventCombinator
	eventGeneratorsMu sync.RWMutex

	wg sync.WaitGroup
}

func NewAsyncEventScheduler(now time.Time) *AsyncEventScheduler {
	return &AsyncEventScheduler{
		clock:           newClock(now),
		eventGenerators: newEventCombinator(),
	}
}

func (a *AsyncEventScheduler) Forward(interval time.Duration) {
	targetTime := a.clock.Now().Add(interval)

	for a.performNextEvent(targetTime) {
	}

	a.wg.Wait()
}

func (a *AsyncEventScheduler) performNextEvent(targetTime time.Time) (shouldContinue bool) {
	a.eventGeneratorsMu.RLock()
	defer a.eventGeneratorsMu.RUnlock()

	if a.eventGenerators.Finished() {
		a.clock.set(targetTime)
		return false
	}

	if a.eventGenerators.Peek().After(targetTime) {
		a.clock.set(targetTime)
		return false
	}

	nextEvent := a.eventGenerators.Pop()
	a.clock.set(nextEvent.Time)

	currentClock := a.clock.copy()
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		nextEvent.Perform(currentClock)
	}()

	return true
}

func (a *AsyncEventScheduler) ForwardToNextEvent() {
	a.eventGeneratorsMu.RLock()
	defer a.eventGeneratorsMu.RUnlock()

	if a.eventGenerators.Finished() {
		return
	}

	nextEvent := a.eventGenerators.Pop()
	a.clock.set(nextEvent.Time)

	currentClock := a.clock.copy()
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		nextEvent.Perform(currentClock)
	}()
	a.wg.Wait()
}

func (a *AsyncEventScheduler) PerformNow(action timing.Action, ctx context.Context) {
	a.AddGenerator(newSingleEventGenerator(action, a.now, ctx))
}

func (a *AsyncEventScheduler) PerformAfter(action timing.Action, interval time.Duration, ctx context.Context) {
	a.AddGenerator(newSingleEventGenerator(action, a.now.Add(interval), ctx))
}

func (a *AsyncEventScheduler) PerformRepeatedly(action timing.Action, until *time.Time, interval time.Duration, ctx context.Context) {
	a.AddGenerator(newPeriodicEventGenerator(action, a.Now(), until, interval, ctx))
}

func (a *AsyncEventScheduler) AddGenerator(generator EventGenerator) {
	a.eventGeneratorsMu.Lock()
	defer a.eventGeneratorsMu.Unlock()

	a.eventGenerators.add(generator)
}
