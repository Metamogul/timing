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

	currentEvent *event

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

	firstEvent := newEvent(action, from.Add(interval))

	return &periodicEventGenerator{
		action:   action,
		from:     from,
		to:       to,
		interval: interval,

		currentEvent: firstEvent,

		ctx: ctx,
	}
}

func (p *periodicEventGenerator) pop() *event {
	if p.finished() {
		panic(ErrEventGeneratorFinished)
	}

	defer func() { p.currentEvent = newEvent(p.action, p.currentEvent.Time.Add(p.interval)) }()

	return p.currentEvent
}

func (p *periodicEventGenerator) peek() event {
	if p.finished() {
		panic(ErrEventGeneratorFinished)
	}

	return *p.currentEvent
}

func (p *periodicEventGenerator) finished() bool {
	if p.ctx.Err() != nil {
		return true
	}

	if p.to == nil {
		return false
	}

	return p.currentEvent.Add(p.interval).After(*p.to)
}
