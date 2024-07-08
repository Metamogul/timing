package simulated_time

import "time"

type periodicEventStream struct {
	action    action
	from      time.Time
	to        *time.Time
	interval  time.Duration
	scheduler eventScheduler

	currentEvent *event
}

func newPeriodicEventStream(
	action action,
	from time.Time,
	to *time.Time,
	interval time.Duration,
	scheduler eventScheduler,
) *periodicEventStream {
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

	if scheduler == nil {
		panic("scheduler can't be nil")
	}

	firstEvent := newEvent(action, from.Add(interval), scheduler)

	return &periodicEventStream{
		action:    action,
		from:      from,
		to:        to,
		interval:  interval,
		scheduler: scheduler,

		currentEvent: firstEvent,
	}
}

func (p *periodicEventStream) popNextEvent() *event {
	if p.closed() {
		panic(ErrEventStreamClosed)
	}

	if !p.currentEvent.actionTime.Before(p.scheduler.Now()) {
		return nil
	}

	currentEvent := p.currentEvent
	p.currentEvent = newEvent(p.action, p.currentEvent.actionTime.Add(p.interval), p.scheduler)

	return currentEvent
}

func (p *periodicEventStream) peakNextTime() time.Time {
	if p.closed() {
		panic(ErrEventStreamClosed)
	}

	return p.currentEvent.actionTime
}

func (p *periodicEventStream) closed() bool {
	if p.to == nil {
		return false
	}

	return !p.scheduler.Now().Before(*p.to)
}
