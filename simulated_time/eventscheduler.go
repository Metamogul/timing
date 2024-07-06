package simulated_time

import (
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
	/*e.now.Add(duration)
	// forward all timers as much as needed:
	for {
		slices.SortStableFunc(e.timedCall, func(a, b *timedCalls) int {
			return a.performAt.Compare(b.performAt)
		})
		// pop next timer and forward until not possible anymore

	}*/
}

// DoAfter spawns a go routine and unblocks it after duration. Don't use
// to repeatedly call f, but use DoRepeatedly instead.
func (e *EventScheduler) DoAfter(duration time.Duration, f func()) {
	/*e.wg.Add(1)

	e.timedCall = append(e.timedCall, e.newDelayedCall(duration, f))
	// Block caller until ready
	// Block EventScheduler if call was in the past -> use wg
	// Release EventScheduler if call is in the future -> use wg
	*/
}

// DoRepeatedly spawns a go routine and unblocks if periodically every
// waiting for interval.
func (e *EventScheduler) DoRepeatedly(interval time.Duration, f func()) {

}

func (e *EventScheduler) eventCompletionWaitGroup() *sync.WaitGroup {
	return &e.wg
}
