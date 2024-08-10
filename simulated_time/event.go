//go:generate go run github.com/vektra/mockery/v2@v2.43.2
package simulated_time

import (
	"github.com/metamogul/timing"
	"time"
)

type Event struct {
	timing.Action
	time.Time
}

func NewEvent(action timing.Action, time time.Time) *Event {
	if action == nil {
		panic("action can't be nil")
	}

	return &Event{
		Action: action,
		Time:   time,
	}
}
