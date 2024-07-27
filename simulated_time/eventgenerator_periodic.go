package simulated_time

import "time"

type periodicEventGenerator struct {
	action   action
	from     time.Time
	to       *time.Time
	interval time.Duration

	currentEvent *event
}

func newPeriodicEventGenerator(
	action action,
	from time.Time,
	to *time.Time,
	interval time.Duration,
) *periodicEventGenerator {
	if action == nil {
		panic("action can't be nil")
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
	if p.to == nil {
		return false
	}

	return !p.currentEvent.Add(p.interval).Before(*p.to)
}
