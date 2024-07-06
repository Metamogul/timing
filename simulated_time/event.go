//go:generate go run github.com/vektra/mockery/v2@v2.43.2
package simulated_time

import (
	"github.com/metamogul/timing"
	"time"
)

type eventScheduler interface {
	timing.Clock
	eventCompletionWaitGroupAdd(delta int)
	eventCompletionWaitGroupDone()
}

type action interface {
	perform()
}

type event struct {
	action     action
	actionTime time.Time
	scheduler  eventScheduler
}

func newEvent(action action, actionTime time.Time, scheduler eventScheduler) *event {
	if action == nil {
		panic("action can't be nil")
	}

	if scheduler == nil {
		panic("scheduler can't be nil")
	}

	return &event{
		action:     action,
		actionTime: actionTime,
		scheduler:  scheduler,
	}
}

func (e *event) performAsync() {
	if e.actionTime.After(e.scheduler.Now()) {
		panic("trying to perform event ahead of actionTime")
	}

	e.scheduler.eventCompletionWaitGroupAdd(1)
	go func() {
		e.action.perform()
		e.scheduler.eventCompletionWaitGroupDone()
	}()
}

func (e *event) perform() {
	if e.actionTime.After(e.scheduler.Now()) {
		panic("trying to perform event ahead of actionTime")
	}

	e.action.perform()
}
