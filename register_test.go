package clock_tester

import (
	"clock-test/clock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_incrementRegisterEveryMinute(t *testing.T) {
	mockClock := clock.NewMock()
	mockClock.Set(time.Now())

	r := &register{}

	incrementRegisterEveryMinute(r, mockClock)

	mockClock.Add(time.Minute * 15)
	require.Equal(t, 15, r.counter)
}
