//go:generate go run github.com/vektra/mockery/v2@v2.43.2
package simulated_time

import (
	"github.com/metamogul/timing"
	"sync"
	"time"
)

type eventScheduler interface {
	timing.Clock
	eventCompletionWaitGroup() *sync.WaitGroup
}

type event struct {
	action     func()
	actionTime time.Time
	scheduler  eventScheduler
}

func newEvent(action func(), actionTime time.Time, scheduler eventScheduler) *event {
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

	e.scheduler.eventCompletionWaitGroup().Add(1)
	go func() {
		e.action()
		e.scheduler.eventCompletionWaitGroup().Done()
	}()
}

func (e *event) perform() {
	if e.actionTime.After(e.scheduler.Now()) {
		panic("trying to perform event ahead of actionTime")
	}

	e.action()
}
