package simulated_time

import (
	"context"
	"github.com/metamogul/timing"
	"sync"
)

const ActionContextEventLoopBlockerKey = "actionContextEventLoopBlocker"

type actionContext struct {
	context.Context

	clock            timing.Clock
	eventLoopBlocker *sync.WaitGroup
}

func newActionContext(ctx context.Context, clock timing.Clock, eventLoopBlocker *sync.WaitGroup) *actionContext {
	return &actionContext{
		Context: ctx,

		clock:            clock,
		eventLoopBlocker: eventLoopBlocker,
	}
}

func (a *actionContext) Clock() timing.Clock {
	return a.clock
}

func (a *actionContext) DoneSchedulingNewEvents() {
	if a.eventLoopBlocker == nil {
		return
	}

	a.eventLoopBlocker.Done()
}

func (a *actionContext) Value(key any) any {
	switch key {
	case timing.ActionContextClockKey:
		return a.clock
	case ActionContextEventLoopBlockerKey:
		return a.eventLoopBlocker
	default:
		return a.Context.Value(key)
	}
}
