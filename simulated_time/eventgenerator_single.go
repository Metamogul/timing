package simulated_time

import (
	"github.com/metamogul/timing"
	"time"
)

type singleEventGenerator struct {
	*event
}

func newSingleEventGenerator(action timing.Action, time time.Time) *singleEventGenerator {
	return &singleEventGenerator{
		event: newEvent(action, time),
	}
}

func (s *singleEventGenerator) pop() *event {
	if s.finished() {
		panic(ErrEventGeneratorFinished)
	}

	defer func() { s.event = nil }()

	return s.event
}

func (s *singleEventGenerator) peek() event {
	if s.finished() {
		panic(ErrEventGeneratorFinished)
	}

	return *s.event
}

func (s *singleEventGenerator) finished() bool {
	return s.event == nil
}
