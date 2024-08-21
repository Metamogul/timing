package simulated_time

import (
	"context"
	"github.com/metamogul/timing"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestNewActionContext(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	actionContextUnderTest := newActionContext(context.Background(), newClock(now), &sync.WaitGroup{})
	require.NotNil(t, actionContextUnderTest)
	require.NotNil(t, actionContextUnderTest.Context)
	require.NotNil(t, actionContextUnderTest.clock)
	require.NotNil(t, actionContextUnderTest.eventLoopBlocker)
}

func TestActionContext_Clock(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	clock := newClock(now)

	actionContextUnderTest := newActionContext(context.Background(), clock, &sync.WaitGroup{})
	gotClock := actionContextUnderTest.Clock()
	require.Equal(t, clock, gotClock)
}

func TestActionContext_DoneSchedulingNewEvents_blockerIsNil(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	actionContextUnderTest := newActionContext(context.Background(), newClock(now), nil)
	actionContextUnderTest.DoneSchedulingNewEvents()
}

func TestActionContext_DoneSchedulingNewEvents_blockerNotNil(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	eventLoopBlocker := &sync.WaitGroup{}

	actionContextUnderTest := newActionContext(context.Background(), newClock(now), eventLoopBlocker)

	eventLoopBlocker.Add(1)
	go func() {
		defer actionContextUnderTest.DoneSchedulingNewEvents()
	}()
	eventLoopBlocker.Wait()
}

func TestActionContext_Value_Clock(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	clock := newClock(now)

	actionContextUnderTest := newActionContext(context.Background(), clock, &sync.WaitGroup{})
	gotClock := actionContextUnderTest.Value(timing.ActionContextClockKey)
	require.Equal(t, clock, gotClock)
}

func TestActionContext_Value_EventLoopBlocker(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	eventLoopBlocker := &sync.WaitGroup{}

	actionContextUnderTest := newActionContext(context.Background(), newClock(now), eventLoopBlocker)
	gotEventLoopBlocker := actionContextUnderTest.Value(ActionContextEventLoopBlockerKey)
	require.Equal(t, eventLoopBlocker, gotEventLoopBlocker)
}

func TestActionContext_Value_Default(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	actionContextUnderTest := newActionContext(context.Background(), newClock(now), &sync.WaitGroup{})
	gotValue := actionContextUnderTest.Value("someNoneExistentKey")
	require.Nil(t, gotValue)
}
