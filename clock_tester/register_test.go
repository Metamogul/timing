package clock_tester

import (
	"github.com/metamogul/timing/simulated_time"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_incrementAfterOneMinute(t *testing.T) {
	mockClock := simulated_time.NewEventScheduler(time.Now())

	r := &register{}

	incrementAfterOneMinute(r, mockClock)

	mockClock.Forward(time.Minute * 15)
	require.Equal(t, 1, r.counter)
}

func Test_incrementEveryMinute(t *testing.T) {
	mockClock := simulated_time.NewEventScheduler(time.Now())

	r := &register{}

	incrementEveryMinute(r, mockClock)

	mockClock.Forward(time.Minute * 15)
	require.Equal(t, 15, r.counter)
}
