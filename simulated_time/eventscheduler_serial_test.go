package simulated_time

import (
	"context"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/metamogul/timing"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewSerialEventScheduler(t *testing.T) {
	t.Parallel()

	now := time.Now()

	newEventScheduler := NewSerialEventScheduler(now)

	require.NotNil(t, newEventScheduler)
	require.IsType(t, &SerialEventScheduler{}, newEventScheduler)

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
		Run(func(ctx timing.ActionContext) {
			time.Sleep(100 * time.Millisecond)

			mu.Lock()
			eventTimes = append(eventTimes, ctx.Clock().Now())
			mu.Unlock()
		}).
		Once()

	longRunningAction2 := timing.NewMockAction(t)
	longRunningAction2.EXPECT().
		Perform(mock.Anything).
		Run(func(ctx timing.ActionContext) {
			time.Sleep(50 * time.Millisecond)

			mu.Lock()
			eventTimes = append(eventTimes, ctx.Clock().Now())
			mu.Unlock()
		}).
		Once()

	eventGenerators := []EventGenerator{
		newSingleEventGenerator(longRunningAction1, now.Add(1*time.Millisecond), context.Background()),
		newSingleEventGenerator(longRunningAction2, now.Add(2*time.Millisecond), context.Background()),
	}

	s := &SerialEventScheduler{
		clock:           newClock(now),
		eventGenerators: newEventCombinator(eventGenerators...),
	}

	s.Forward(3 * time.Millisecond)

	sorted := slices.IsSortedFunc(eventTimes, func(a, b time.Time) int {
		return a.Compare(b)
	})
	require.True(t, sorted)
}

func TestSerialEventScheduler_Forward_RecursiveScheduling(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	eventTimes := make([]time.Time, 0)

	s := &SerialEventScheduler{
		clock:           newClock(now),
		eventGenerators: newEventCombinator(),
	}

	innerAction := timing.NewMockAction(t)
	innerAction.EXPECT().
		Perform(mock.Anything).
		Run(func(ctx timing.ActionContext) {
			eventTimes = append(eventTimes, ctx.Clock().Now())
		}).
		Once()

	outerAction := timing.NewMockAction(t)
	outerAction.EXPECT().
		Perform(mock.Anything).
		Run(func(ctx timing.ActionContext) {
			s.PerformAfter(innerAction, time.Second, context.Background())
			eventTimes = append(eventTimes, ctx.Clock().Now())
		}).
		Once()

	s.PerformAfter(outerAction, time.Second, context.Background())

	s.Forward(3 * time.Second)

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
		eventGenerators    func() []EventGenerator
		wantShouldContinue bool
	}{
		{
			name:            "all event generators finished",
			eventGenerators: func() []EventGenerator { return nil },
		},
		{
			name: "next event after target time",
			eventGenerators: func() []EventGenerator {
				return []EventGenerator{newSingleEventGenerator(timing.NewMockAction(t), now.Add(1*time.Hour), context.Background())}
			},
		},
		{
			name: "event dispatched successfully",
			eventGenerators: func() []EventGenerator {
				mockAction := timing.NewMockAction(t)
				mockAction.EXPECT().
					Perform(mock.Anything).
					Once()
				return []EventGenerator{newSingleEventGenerator(mockAction, now.Add(1*time.Second), context.Background())}
			},
			wantShouldContinue: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &SerialEventScheduler{
				clock:           newClock(now),
				eventGenerators: newEventCombinator(tt.eventGenerators()...),
			}

			if gotShouldContinue := s.performNextEvent(targetTime); gotShouldContinue != tt.wantShouldContinue {
				t.Errorf("performNextEventSerially() = %v, want %v", gotShouldContinue, tt.wantShouldContinue)
			}

			if tt.wantShouldContinue == true {
				require.Equal(t, now.Add(time.Second), s.Now())
			} else {
				require.Equal(t, targetTime, s.Now())
			}
		})
	}
}

func TestSerialEventScheduler_ForwardToNextEvent(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	eventTimes := make([]time.Time, 0)

	longRunningAction1 := timing.NewMockAction(t)
	longRunningAction1.EXPECT().
		Perform(mock.Anything).
		Run(func(ctx timing.ActionContext) {
			time.Sleep(100 * time.Millisecond)
			eventTimes = append(eventTimes, ctx.Clock().Now())
		}).
		Once()

	longRunningAction2 := timing.NewMockAction(t)
	longRunningAction2.EXPECT().
		Perform(mock.Anything).
		Run(func(ctx timing.ActionContext) {
			time.Sleep(50 * time.Millisecond)
			eventTimes = append(eventTimes, ctx.Clock().Now())
		}).
		Once()

	eventGenerators := []EventGenerator{
		newSingleEventGenerator(longRunningAction1, now.Add(1*time.Second), context.Background()),
		newSingleEventGenerator(longRunningAction2, now.Add(2*time.Second), context.Background()),
	}

	s := &SerialEventScheduler{
		clock:           newClock(now),
		eventGenerators: newEventCombinator(eventGenerators...),
	}

	s.ForwardToNextEvent()
	require.Len(t, eventTimes, 1)
	require.Equal(t, now.Add(1*time.Second), eventTimes[0])
	require.Equal(t, now.Add(1*time.Second), s.Now())

	s.ForwardToNextEvent()
	require.Len(t, eventTimes, 2)
	require.Equal(t, now.Add(2*time.Second), eventTimes[1])
	require.Equal(t, now.Add(2*time.Second), s.Now())
}

func TestSerialEventScheduler_PerformNow(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	s := &SerialEventScheduler{
		clock:           newClock(now),
		eventGenerators: newEventCombinator(),
	}
	s.PerformNow(timing.NewMockAction(t), context.Background())

	require.Len(t, s.eventGenerators.activeGenerators, 1)
	require.IsType(t, &singleEventGenerator{}, s.eventGenerators.activeGenerators[0])
}

func TestSerialEventScheduler_PerformAfter(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	s := &SerialEventScheduler{
		clock:           newClock(now),
		eventGenerators: newEventCombinator(),
	}
	s.PerformAfter(timing.NewMockAction(t), time.Second, context.Background())

	require.Len(t, s.eventGenerators.activeGenerators, 1)
	require.IsType(t, &singleEventGenerator{}, s.eventGenerators.activeGenerators[0])
}

func TestSerialEventScheduler_PerformRepeatedly(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	s := &SerialEventScheduler{
		clock:           newClock(now),
		eventGenerators: newEventCombinator(),
	}
	s.PerformRepeatedly(timing.NewMockAction(t), nil, time.Second, context.Background())

	require.Len(t, s.eventGenerators.activeGenerators, 1)
	require.IsType(t, &periodicEventGenerator{}, s.eventGenerators.activeGenerators[0])
}

func TestSerialEventScheduler_AddGenerator(t *testing.T) {
	t.Parallel()

	mockEventGenerator := NewMockEventGenerator(t)
	mockEventGenerator.EXPECT().
		Finished().
		Return(false).
		Once()

	s := &SerialEventScheduler{
		eventGenerators: newEventCombinator(),
	}
	s.AddGenerator(mockEventGenerator)

	require.Len(t, s.eventGenerators.activeGenerators, 1)
}
