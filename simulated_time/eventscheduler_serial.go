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
	if s.eventGenerators.Finished() {
		s.clock.set(targetTime)
		return false
	}

	if s.eventGenerators.Peek().After(targetTime) {
		s.clock.set(targetTime)
		return false
	}

	nextEvent := s.eventGenerators.Pop()
	s.clock.set(nextEvent.Time)

	nextEvent.Perform(newActionContext(nextEvent.Context, s.clock.copy(), nil))

	return true
}

func (s *SerialEventScheduler) ForwardToNextEvent() {
	if s.eventGenerators.Finished() {
		return
	}

	nextEvent := s.eventGenerators.Pop()
	s.clock.set(nextEvent.Time)

	nextEvent.Perform(newActionContext(nextEvent.Context, s.clock.copy(), nil))
}

func (s *SerialEventScheduler) PerformNow(action timing.Action, ctx context.Context) {
	s.AddGenerator(newSingleEventGenerator(action, s.now, ctx))
}

func (s *SerialEventScheduler) PerformAfter(action timing.Action, interval time.Duration, ctx context.Context) {
	s.AddGenerator(newSingleEventGenerator(action, s.now.Add(interval), ctx))
}

func (s *SerialEventScheduler) PerformRepeatedly(action timing.Action, until *time.Time, interval time.Duration, ctx context.Context) {
	s.AddGenerator(newPeriodicEventGenerator(action, s.Now(), until, interval, ctx))
}

func (s *SerialEventScheduler) AddGenerator(generator EventGenerator) {
	s.eventGenerators.add(generator)
}
