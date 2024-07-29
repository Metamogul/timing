package simulated_time

import (
	"context"
	"github.com/metamogul/timing"
	"time"
)

type singleEventGenerator struct {
	*event
	ctx context.Context
}

func newSingleEventGenerator(action timing.Action, time time.Time, ctx context.Context) *singleEventGenerator {
	return &singleEventGenerator{
		event: newEvent(action, time),
		ctx:   ctx,
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
	return s.event == nil || s.ctx.Err() != nil
}
