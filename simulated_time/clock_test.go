package simulated_time

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewClock(t *testing.T) {
	t.Parallel()

	now := time.Now()

	clock := newClock(now)

	require.NotNil(t, clock)
	require.Equal(t, now, clock.Now())
}

func TestClock_Now(t *testing.T) {
	t.Parallel()

	now := time.Now()

	clock := clock{now}
	require.Equal(t, now, clock.Now())
}

func TestClock_Forward(t *testing.T) {
	t.Parallel()

	now := time.Now()

	clock := clock{now}
	clock.Forward(time.Minute)

	require.Equal(t, now.Add(time.Minute), clock.now)
}

func Test_clock_Set(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name         string
		now          time.Time
		newTime      time.Time
		requirePanic bool
	}{
		{
			name:         "new time in the past",
			now:          now,
			newTime:      now.Add(-time.Second),
			requirePanic: true,
		},
		{
			name:    "new time equals current time",
			now:     now,
			newTime: now,
		},
		{
			name:    "new time after curent time",
			now:     now,
			newTime: now.Add(time.Second),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := &clock{
				now: tt.now,
			}

			if tt.requirePanic {
				require.Panics(t, func() {
					c.Set(tt.newTime)
				})
				return
			}

			c.Set(tt.newTime)
			require.Equal(t, tt.newTime, c.now)
		})
	}
}
