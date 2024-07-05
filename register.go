package clock_tester

import (
	"clock-test/clock"
	"time"
)

type register struct {
	counter int
}

func incrementRegisterEveryMinute(register *register, clock clock.Clock) {
	go func() {
		for {
			<-clock.After(time.Minute)
			time.Sleep(time.Microsecond * 1120)
			register.counter++
		}
	}()

	time.Sleep(time.Millisecond)
}
