package simulated_time

import (
	"context"
	"github.com/metamogul/timing"
	"time"
)

type singleEventGenerator struct {
	*Event
	ctx context.Context
}

func newSingleEventGenerator(action timing.Action, time time.Time, ctx context.Context) *singleEventGenerator {
	return &singleEventGenerator{
		Event: NewEvent(action, time),
		ctx:   ctx,
	}
}

func (s *singleEventGenerator) Pop() *Event {
	if s.Finished() {
		panic(ErrEventGeneratorFinished)
	}

	defer func() { s.Event = nil }()

	return s.Event
}

func (s *singleEventGenerator) Peek() Event {
	if s.Finished() {
		panic(ErrEventGeneratorFinished)
	}

	return *s.Event
}

func (s *singleEventGenerator) Finished() bool {
	return s.Event == nil || s.ctx.Err() != nil
}
