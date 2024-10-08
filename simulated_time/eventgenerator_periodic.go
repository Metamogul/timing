package simulated_time

import (
	"context"
	"github.com/metamogul/timing"
	"time"
)

type periodicEventGenerator struct {
	action   timing.Action
	from     time.Time
	to       *time.Time
	interval time.Duration

	currentEvent *Event

	ctx context.Context
}

func newPeriodicEventGenerator(
	action timing.Action,
	from time.Time,
	to *time.Time,
	interval time.Duration,
	ctx context.Context,
) *periodicEventGenerator {
	if action == nil {
		panic("Action can't be nil")
	}

	if to != nil && !to.After(from) {
		panic("to must be after from")
	}

	if interval == 0 {
		panic("interval must be greater than zero")
	}

	if to != nil && interval >= to.Sub(from) {
		panic("interval must be shorter than timespan given by from and to")
	}

	firstEvent := NewEvent(action, from.Add(interval), ctx)

	return &periodicEventGenerator{
		action:   action,
		from:     from,
		to:       to,
		interval: interval,

		currentEvent: firstEvent,

		ctx: ctx,
	}
}

func (p *periodicEventGenerator) Pop() *Event {
	if p.Finished() {
		panic(ErrEventGeneratorFinished)
	}

	defer func() { p.currentEvent = NewEvent(p.action, p.currentEvent.Time.Add(p.interval), p.ctx) }()

	return p.currentEvent
}

func (p *periodicEventGenerator) Peek() Event {
	if p.Finished() {
		panic(ErrEventGeneratorFinished)
	}

	return *p.currentEvent
}

func (p *periodicEventGenerator) Finished() bool {
	if p.ctx.Err() != nil {
		return true
	}

	if p.to == nil {
		return false
	}

	return p.currentEvent.Add(p.interval).After(*p.to)
}
