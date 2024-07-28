package clock_tester

import (
	"github.com/metamogul/timing/simulated_time"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_incrementAfterOneMinute(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	scheduler := simulated_time.NewAsyncEventScheduler(now)

	r := &register{}

	r.incrementAfterOneMinute(scheduler)

	scheduler.Forward(time.Minute * 15)
	require.Equal(t, 1, r.counter)
}

func Test_incrementEveryMinute(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	scheduler := simulated_time.NewAsyncEventScheduler(now)

	r := &register{}

	r.incrementEveryMinute(scheduler)

	scheduler.Forward(time.Minute * 60)
	require.Equal(t, 60, r.counter)
}
