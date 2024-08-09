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

func (s *SerialEventScheduler) Forward(interval time.Duration) {
	targetTime := s.clock.Now().Add(interval)

	for s.performNextEvent(targetTime) {
	}
}

func (s *SerialEventScheduler) performNextEvent(targetTime time.Time) (shouldContinue bool) {
	if s.eventGenerators.finished() {
		s.clock.set(targetTime)
		return false
	}

	if s.eventGenerators.peek().After(targetTime) {
		s.clock.set(targetTime)
		return false
	}

	nextEvent := s.eventGenerators.pop()
	s.clock.set(nextEvent.Time)

	nextEvent.Perform(s.clock.copy())

	return true
}

func (s *SerialEventScheduler) ForwardToNextEvent() {
	if s.eventGenerators.finished() {
		return
	}

	nextEvent := s.eventGenerators.pop()
	s.clock.set(nextEvent.Time)

	nextEvent.Perform(s.clock.copy())
}

func (a *SerialEventScheduler) PerformNow(action timing.Action, ctx context.Context) {
	a.eventGenerators.add(newSingleEventGenerator(action, a.now, ctx))
}

func (s *SerialEventScheduler) PerformAfter(action timing.Action, interval time.Duration, ctx context.Context) {
	s.eventGenerators.add(newSingleEventGenerator(action, s.now.Add(interval), ctx))
}

func (s *SerialEventScheduler) PerformRepeatedly(action timing.Action, until *time.Time, interval time.Duration, ctx context.Context) {
	s.eventGenerators.add(newPeriodicEventGenerator(action, s.Now(), until, interval, ctx))
}
