package clock_tester

import (
	"github.com/metamogul/timing"
	"time"
)

type register struct {
	counter int
}

func incrementAfterOneMinute(register *register, scheduler timing.EventScheduler) {
	scheduler.DoAfter(time.Minute, func() {
		time.Sleep(time.Microsecond * 1120)
		register.counter++
	})
}

func incrementEveryMinute(register *register, scheduler timing.EventScheduler) {
	scheduler.DoRepeatedly(time.Minute, func() {
		time.Sleep(time.Microsecond * 1120)
		register.counter++
	})
}
