package simulated_time

import (
	"context"
	"github.com/metamogul/timing"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"slices"
	"sync"
	"testing"
	"time"
)

func TestNewSerialEventScheduler(t *testing.T) {
	t.Parallel()

	now := time.Now()

	newEventScheduler := NewSerialEventScheduler(now)

	require.NotNil(t, newEventScheduler)

	require.NotNil(t, newEventScheduler.clock)
	require.Equal(t, now, newEventScheduler.Now())

	require.NotNil(t, newEventScheduler.eventGenerators)
}

func TestSerialEventScheduler_Forward(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	mu := sync.Mutex{}
	eventTimes := make([]time.Time, 0)

	longRunningAction1 := timing.NewMockAction(t)
	longRunningAction1.EXPECT().
		Perform(mock.Anything).
		Run(func(clock timing.Clock) {
			time.Sleep(100 * time.Millisecond)

			mu.Lock()
			eventTimes = append(eventTimes, clock.Now())
			mu.Unlock()
		}).
		Once()

	longRunningAction2 := timing.NewMockAction(t)
	longRunningAction2.EXPECT().
		Perform(mock.Anything).
		Run(func(clock timing.Clock) {
			time.Sleep(50 * time.Millisecond)

			mu.Lock()
			eventTimes = append(eventTimes, clock.Now())
			mu.Unlock()
		}).
		Once()

	eventGenerators := []eventGenerator{
		newSingleEventGenerator(longRunningAction1, now.Add(1*time.Millisecond), context.Background()),
		newSingleEventGenerator(longRunningAction2, now.Add(2*time.Millisecond), context.Background()),
	}

	e := &SerialEventScheduler{
		clock:           newClock(now),
		eventGenerators: newEventCombinator(eventGenerators...),
	}

	e.Forward(3 * time.Millisecond)

	sorted := slices.IsSortedFunc(eventTimes, func(a, b time.Time) int {
		return a.Compare(b)
	})
	require.True(t, sorted)
}

func TestSerialEventScheduler_Forward_RecursiveScheduling(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	eventTimes := make([]time.Time, 0)

	e := &SerialEventScheduler{
		clock:           newClock(now),
		eventGenerators: newEventCombinator(),
	}

	innerAction := timing.NewMockAction(t)
	innerAction.EXPECT().
		Perform(mock.Anything).
		Run(func(clock timing.Clock) {
			eventTimes = append(eventTimes, clock.Now())
		}).
		Once()

	outerAction := timing.NewMockAction(t)
	outerAction.EXPECT().
		Perform(mock.Anything).
		Run(func(clock timing.Clock) {
			e.PerformAfter(innerAction, time.Second, context.Background())
			eventTimes = append(eventTimes, clock.Now())
		}).
		Once()

	e.PerformAfter(outerAction, time.Second, context.Background())

	e.Forward(3 * time.Second)

	sorted := slices.IsSortedFunc(eventTimes, func(a, b time.Time) int {
		return a.Compare(b)
	})
	require.True(t, sorted)
}

func TestSerialEventScheduler_performNextEvent(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	targetTime := now.Add(time.Minute)

	tests := []struct {
		name               string
		eventGenerators    func() []eventGenerator
		wantShouldContinue bool
	}{
		{
			name:            "all event generators finished",
			eventGenerators: func() []eventGenerator { return nil },
		},
		{
			name: "next event after target time",
			eventGenerators: func() []eventGenerator {
				return []eventGenerator{newSingleEventGenerator(timing.NewMockAction(t), now.Add(1*time.Hour), context.Background())}
			},
		},
		{
			name: "event dispatched successfully",
			eventGenerators: func() []eventGenerator {
				mockAction := timing.NewMockAction(t)
				mockAction.EXPECT().
					Perform(mock.Anything).
					Once()
				return []eventGenerator{newSingleEventGenerator(mockAction, now.Add(1*time.Second), context.Background())}
			},
			wantShouldContinue: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := &SerialEventScheduler{
				clock:           newClock(now),
				eventGenerators: newEventCombinator(tt.eventGenerators()...),
			}

			if gotShouldContinue := e.performNextEvent(targetTime); gotShouldContinue != tt.wantShouldContinue {
				t.Errorf("performNextEventSerially() = %v, want %v", gotShouldContinue, tt.wantShouldContinue)
			}

			if tt.wantShouldContinue == true {
				require.Equal(t, now.Add(time.Second), e.Now())
			} else {
				require.Equal(t, targetTime, e.Now())
			}
		})
	}
}

func TestSerialEventScheduler_PerformAfter(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	e := &SerialEventScheduler{
		clock:           newClock(now),
		eventGenerators: newEventCombinator(),
	}
	e.PerformAfter(timing.NewMockAction(t), time.Second, context.Background())

	require.Len(t, e.eventGenerators.activeGenerators, 1)
	require.IsType(t, &singleEventGenerator{}, e.eventGenerators.activeGenerators[0])
}

func TestSerialEventScheduler_PerformRepeatedly(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	e := &SerialEventScheduler{
		clock:           newClock(now),
		eventGenerators: newEventCombinator(),
	}
	e.PerformRepeatedly(timing.NewMockAction(t), nil, time.Second, context.Background())

	require.Len(t, e.eventGenerators.activeGenerators, 1)
	require.IsType(t, &periodicEventGenerator{}, e.eventGenerators.activeGenerators[0])
}
