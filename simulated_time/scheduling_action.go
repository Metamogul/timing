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
	return SchedulingAction{
		Action:           action,
		eventLoopBlocker: &sync.WaitGroup{},
	}
}
