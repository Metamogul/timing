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

type Action func(timing.ActionContext)

func (a Action) Perform(ctx timing.ActionContext) { a(ctx) }

func (r *register) incrementAfterOneMinute(scheduler timing.EventScheduler) {
	scheduler.PerformAfter(
		Action(func(timing.ActionContext) {
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
		Action(func(timing.ActionContext) {
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
