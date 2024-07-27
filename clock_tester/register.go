package clock_tester

import (
	"github.com/metamogul/timing"
	"sync"
	"time"
)

type register struct {
	counter int
}

type Action func()

func (a Action) Perform() { a() }

func (r *register) incrementAfterOneMinute(scheduler timing.EventScheduler) {
	scheduler.PerformAfter(Action(func() {
		// Simulate execution time
		time.Sleep(100 * time.Millisecond)

		r.counter++
	}), time.Minute)
}

func (r *register) incrementEveryMinute(scheduler timing.EventScheduler) {
	mu := sync.Mutex{}

	scheduler.PerformRepeatedly(Action(func() {
		mu.Lock()

		// Simulate execution time
		time.Sleep(10 * time.Millisecond)

		r.counter++

		mu.Unlock()
	}), nil, time.Minute)
}
