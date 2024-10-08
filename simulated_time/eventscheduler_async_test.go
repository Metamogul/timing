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

func TestNewAsyncEventScheduler(t *testing.T) {
	t.Parallel()

	now := time.Now()

	newEventScheduler := NewAsyncEventScheduler(now)

	require.NotNil(t, newEventScheduler)
	require.IsType(t, &AsyncEventScheduler{}, newEventScheduler)

	require.NotNil(t, newEventScheduler.clock)
	require.Equal(t, now, newEventScheduler.Now())

	require.NotNil(t, newEventScheduler.eventGenerators)
}

func TestAsyncEventScheduler_Forward(t *testing.T) {
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

	a := &AsyncEventScheduler{
		clock:           newClock(now),
		eventGenerators: newEventCombinator(eventGenerators...),
	}

	a.Forward(3 * time.Millisecond)

	require.Equal(t, now.Add(2*time.Millisecond), eventTimes[0])
	require.Equal(t, now.Add(1*time.Millisecond), eventTimes[1])
}

func TestAsyncEventScheduler_Forward_PeriodicSchedulingAction(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	eventTimes := make([]time.Time, 0)

	action := timing.NewMockAction(t)
	action.EXPECT().
		Perform(mock.Anything).
		Run(func(ctx timing.ActionContext) {
			eventTimes = append(eventTimes, ctx.Clock().Now())
			ctx.DoneSchedulingNewEvents()
		}).
		Times(3)
	schedulingAction := NewSchedulingAction(action)

	eventGenerators := []EventGenerator{
		newPeriodicEventGenerator(schedulingAction, now, nil, time.Millisecond, context.Background()),
	}

	a := &AsyncEventScheduler{
		clock:           newClock(now),
		eventGenerators: newEventCombinator(eventGenerators...),
	}

	a.Forward(3 * time.Millisecond)

	require.Equal(t, now.Add(1*time.Millisecond), eventTimes[0])
	require.Equal(t, now.Add(2*time.Millisecond), eventTimes[1])
	require.Equal(t, now.Add(3*time.Millisecond), eventTimes[2])
}

func TestAsyncEventScheduler_Forward_RecursiveScheduling(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	eventTimes := make([]time.Time, 0)

	s := &AsyncEventScheduler{
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
			ctx.DoneSchedulingNewEvents()
			eventTimes = append(eventTimes, ctx.Clock().Now())
		}).
		Once()

	s.PerformAfter(NewSchedulingAction(outerAction), time.Second, context.Background())

	s.Forward(3 * time.Second)

	sorted := slices.IsSortedFunc(eventTimes, func(a, b time.Time) int {
		return a.Compare(b)
	})
	require.True(t, sorted)
}

func TestAsyncEventScheduler_performNextEvent(t *testing.T) {
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

			a := &AsyncEventScheduler{
				clock:           newClock(now),
				eventGenerators: newEventCombinator(tt.eventGenerators()...),
			}

			if gotShouldContinue := a.performNextEvent(targetTime); gotShouldContinue != tt.wantShouldContinue {
				t.Errorf("performNextEvent() = %v, want %v", gotShouldContinue, tt.wantShouldContinue)
			}
			a.wg.Wait()

			if tt.wantShouldContinue == true {
				require.Equal(t, now.Add(time.Second), a.Now())
			} else {
				require.Equal(t, targetTime, a.Now())
			}
		})
	}
}

func TestAsyncEventScheduler_ForwardToNextEvent(t *testing.T) {
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
		newSingleEventGenerator(longRunningAction1, now.Add(1*time.Second), context.Background()),
		newSingleEventGenerator(longRunningAction2, now.Add(2*time.Second), context.Background()),
	}

	a := &AsyncEventScheduler{
		clock:           newClock(now),
		eventGenerators: newEventCombinator(eventGenerators...),
	}

	a.ForwardToNextEvent()
	require.Len(t, eventTimes, 1)
	require.Equal(t, now.Add(1*time.Second), eventTimes[0])
	require.Equal(t, now.Add(1*time.Second), a.Now())

	a.ForwardToNextEvent()
	require.Len(t, eventTimes, 2)
	require.Equal(t, now.Add(2*time.Second), eventTimes[1])
	require.Equal(t, now.Add(2*time.Second), a.Now())
}

func TestAsyncEventScheduler_ForwardToNextEvent_SchedulingAction(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	mu := sync.Mutex{}
	eventTimes := make([]time.Time, 0)

	a := &AsyncEventScheduler{
		clock: newClock(now),
	}

	innerAction := timing.NewMockAction(t)
	innerAction.EXPECT().
		Perform(mock.Anything).
		Run(func(ctx timing.ActionContext) {
			mu.Lock()
			eventTimes = append(eventTimes, ctx.Clock().Now())
			mu.Unlock()
		}).
		Once()

	outerAction := timing.NewMockAction(t)
	outerAction.EXPECT().
		Perform(mock.Anything).
		Run(func(ctx timing.ActionContext) {
			a.PerformAfter(innerAction, time.Second, ctx)
			ctx.DoneSchedulingNewEvents()

			mu.Lock()
			eventTimes = append(eventTimes, ctx.Clock().Now())
			mu.Unlock()
		}).
		Once()

	outerSchedulingAction := NewSchedulingAction(outerAction)

	eventGenerators := []EventGenerator{
		newSingleEventGenerator(outerSchedulingAction, now.Add(1*time.Second), context.Background()),
	}

	a.eventGenerators = newEventCombinator(eventGenerators...)

	a.ForwardToNextEvent()
	require.Len(t, eventTimes, 1)
	require.Equal(t, now.Add(1*time.Second), eventTimes[0])
	require.Equal(t, now.Add(1*time.Second), a.Now())

	a.ForwardToNextEvent()
	require.Len(t, eventTimes, 2)
	require.Equal(t, now.Add(2*time.Second), eventTimes[1])
	require.Equal(t, now.Add(2*time.Second), a.Now())
}

func TestAsyncEventScheduler_PerformNow(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	a := &AsyncEventScheduler{
		clock:           newClock(now),
		eventGenerators: newEventCombinator(),
	}
	a.PerformNow(timing.NewMockAction(t), context.Background())

	require.Len(t, a.eventGenerators.activeGenerators, 1)
	require.IsType(t, &singleEventGenerator{}, a.eventGenerators.activeGenerators[0])
}

func TestAsyncEventScheduler_PerformAfter(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	a := &AsyncEventScheduler{
		clock:           newClock(now),
		eventGenerators: newEventCombinator(),
	}
	a.PerformAfter(timing.NewMockAction(t), time.Second, context.Background())

	require.Len(t, a.eventGenerators.activeGenerators, 1)
	require.IsType(t, &singleEventGenerator{}, a.eventGenerators.activeGenerators[0])
}

func TestAsyncEventScheduler_PerformRepeatedly(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	a := &AsyncEventScheduler{
		clock:           newClock(now),
		eventGenerators: newEventCombinator(),
	}
	a.PerformRepeatedly(timing.NewMockAction(t), nil, time.Second, context.Background())

	require.Len(t, a.eventGenerators.activeGenerators, 1)
	require.IsType(t, &periodicEventGenerator{}, a.eventGenerators.activeGenerators[0])
}

func TestAsyncEventScheduler_AddGenerator(t *testing.T) {
	t.Parallel()

	mockEventGenerator := NewMockEventGenerator(t)
	mockEventGenerator.EXPECT().
		Finished().
		Return(false).
		Once()

	a := &AsyncEventScheduler{
		eventGenerators: newEventCombinator(),
	}
	a.AddGenerator(mockEventGenerator)

	require.Len(t, a.eventGenerators.activeGenerators, 1)
}
