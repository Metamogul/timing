package simulated_time

import (
	"github.com/metamogul/timing"
	"sync"
	"time"
)

type EventScheduler struct {
	*Clock

	wg sync.WaitGroup
}

func NewEventScheduler(now time.Time) *EventScheduler {
	return &EventScheduler{
		Clock: NewClock(now),
	}
}

func (e *EventScheduler) Forward(duration time.Duration) {
	/*

		- set new time of clock
		- while peek event time before current time:
			- add one to wait group
			- dispatch new event in go routine and wait with waitgroup
	*/

	/*e.now.Add(duration)
	// forward all timers as much as needed:
	for {
		slices.SortStableFunc(e.timedCall, func(a, b *timedCalls) int {
			return a.performAt.Compare(b.performAt)
		})
		// pop next timer and forward until not possible anymore

	}*/
}

// PerformAfter spawns a go routine and unblocks it after duration. Don't use
// to repeatedly call f, but use DoRepeatedly instead.
func (e *EventScheduler) PerformAfter(duration time.Duration, action timing.Action) {
	/*e.wg.Add(1)

	e.timedCall = append(e.timedCall, e.newDelayedCall(duration, f))
	// Block caller until ready
	// Block EventScheduler if call was in the past -> use wg
	// Release EventScheduler if call is in the future -> use wg
	*/
}

// PerformRepeatedly spawns a go routine and unblocks if periodically every
// waiting for interval.
func (e *EventScheduler) PerformRepeatedly(interval time.Duration, action timing.Action) {

}
