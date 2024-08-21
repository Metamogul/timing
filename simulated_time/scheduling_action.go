package simulated_time

import (
	"github.com/metamogul/timing"
	"sync"
)

type SchedulingAction struct {
	timing.Action
	eventLoopBlocker *sync.WaitGroup
}

func NewSchedulingAction(action timing.Action) SchedulingAction {
	eventLoopBlocker := &sync.WaitGroup{}
	eventLoopBlocker.Add(1)

	return SchedulingAction{
		Action:           action,
		eventLoopBlocker: eventLoopBlocker,
	}
}

func (s SchedulingAction) WaitForEventSchedulingCompletion() {
	s.eventLoopBlocker.Wait()
}
