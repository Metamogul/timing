package system

import (
	"context"
	"github.com/metamogul/timing"
)

type actionContext struct {
	context.Context
	clock timing.Clock
}

func newActionContext(ctx context.Context, clock timing.Clock) *actionContext {
	return &actionContext{
		Context: ctx,
		clock:   clock,
	}
}

func (a *actionContext) Clock() timing.Clock {
	return a.clock
}

func (a *actionContext) DoneSchedulingNewEvents() { /*Noop*/ }

func (a *actionContext) Value(key any) any {
	switch key {
	case timing.ActionContextClockKey:
		return a.clock
	default:
		return a.Context.Value(key)
	}
}
