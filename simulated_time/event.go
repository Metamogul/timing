//go:generate go run github.com/vektra/mockery/v2@v2.43.2
package simulated_time

import (
	"github.com/metamogul/timing"
	"time"
)

type event struct {
	timing.Action
	time.Time
}

func newEvent(action timing.Action, time time.Time) *event {
	if action == nil {
		panic("action can't be nil")
	}

	return &event{
		Action: action,
		Time:   time,
	}
}
