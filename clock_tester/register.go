package clock_tester

import (
	"context"
	"github.com/metamogul/timing"
	"sync"
	"time"
)

type register struct {
	counter int
}

type Action func(timing.Clock)

func (a Action) Perform(clock timing.Clock) { a(clock) }

func (r *register) incrementAfterOneMinute(scheduler timing.EventScheduler) {
	scheduler.PerformAfter(
		Action(func(timing.Clock) {
			// Simulate execution time
			time.Sleep(100 * time.Millisecond)

			r.counter++
		}),
		time.Minute,
		context.Background(),
	)
}

func (r *register) incrementEveryMinute(scheduler timing.EventScheduler) {
	mu := sync.Mutex{}

	scheduler.PerformRepeatedly(
		Action(func(timing.Clock) {
			mu.Lock()

			// Simulate execution time
			time.Sleep(10 * time.Millisecond)

			r.counter++

			mu.Unlock()
		}),
		nil,
		time.Minute,
		context.Background(),
	)
}
