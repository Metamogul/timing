package simulated_time

import (
	"time"
)

type Clock struct {
	now time.Time
}

func NewClock(now time.Time) *Clock {
	return &Clock{
		now: now,
	}
}

func (c *Clock) Now() time.Time {
	return c.now
}

func (c *Clock) Forward(d time.Duration) {
	c.now = c.now.Add(d)
}
