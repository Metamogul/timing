//go:generate go run github.com/vektra/mockery/v2@v2.43.2
package simulated_time

import (
	"time"
)

type action interface {
	perform()
}

type event struct {
	action
	time.Time
}

func newEvent(action action, time time.Time) *event {
	if action == nil {
		panic("action can't be nil")
	}

	return &event{
		action: action,
		Time:   time,
	}
}
