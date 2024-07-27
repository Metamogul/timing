package simulated_time

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewEventScheduler(t *testing.T) {
	t.Parallel()

	now := time.Now()

	newEventScheduler := NewEventScheduler(now)

	require.NotNil(t, newEventScheduler)

	require.NotNil(t, newEventScheduler.clock)
	require.Equal(t, now, newEventScheduler.Now())

	require.NotNil(t, newEventScheduler.eventGenerators)
}

func TestEventScheduler_Forward(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	longRunningAction1 := NewMockAction(t)
	longRunningAction1.EXPECT().
		Perform().
		Run(func() {
			time.Sleep(50 * time.Millisecond)
		}).
		Once()

	longRunningAction2 := NewMockAction(t)
	longRunningAction2.EXPECT().
		Perform().
		Run(func() {
			time.Sleep(60 * time.Millisecond)
		}).
		Once()

	eventGenerators := []eventGenerator{
		newSingleEventGenerator(longRunningAction1, now.Add(1*time.Millisecond)),
		newSingleEventGenerator(longRunningAction2, now.Add(2*time.Millisecond)),
	}

	e := &EventScheduler{
		clock:           newClock(now),
		eventGenerators: newEventCombinator(eventGenerators...),
	}

	e.Forward(3 * time.Millisecond)
}

func TestEventScheduler_Forward_RecursiveScheduling(t *testing.T) {
	// TODO: move this test to strictly monotonous scheduler

	/*
		t.Parallel()

		now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		e := &EventScheduler{
			clock:           newClock(now),
			eventGenerators: newEventCombinator(),
		}

		innerAction := NewMockAction(t)
		innerAction.EXPECT().
			Perform().
			Run(func() {
				fmt.Println("called inner action")
			}).
			Once()

		outerAction := NewMockAction(t)
		outerAction.EXPECT().
			Perform().
			Run(func() {
				fmt.Println("called outer action")
				e.PerformAfter(innerAction, time.Second)
			}).
			Once()

		e.PerformAfter(outerAction, time.Second)

		e.Forward(3 * time.Second)
	*/
}

func TestEventScheduler_dispatchNextEvent(t *testing.T) {
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
				return []eventGenerator{newSingleEventGenerator(NewMockAction(t), now.Add(1*time.Hour))}
			},
		},
		{
			name: "event dispatched successfully",
			eventGenerators: func() []eventGenerator {
				mockAction := NewMockAction(t)
				mockAction.EXPECT().
					Perform().
					Once()
				return []eventGenerator{newSingleEventGenerator(mockAction, now.Add(1*time.Second))}
			},
			wantShouldContinue: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := &EventScheduler{
				clock:           newClock(now),
				eventGenerators: newEventCombinator(tt.eventGenerators()...),
			}

			if gotShouldContinue := e.dispatchNextEvent(targetTime); gotShouldContinue != tt.wantShouldContinue {
				t.Errorf("dispatchNextEvent() = %v, want %v", gotShouldContinue, tt.wantShouldContinue)
			}
			e.wg.Wait()

			if tt.wantShouldContinue == true {
				require.Equal(t, now.Add(time.Second), e.Now())
			} else {
				require.Equal(t, targetTime, e.Now())
			}
		})
	}
}

func TestEventScheduler_PerformAfter(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	e := &EventScheduler{
		clock:           newClock(now),
		eventGenerators: newEventCombinator(),
	}
	e.PerformAfter(NewMockAction(t), time.Second)

	require.Len(t, e.eventGenerators.inputs, 1)
	require.IsType(t, &singleEventGenerator{}, e.eventGenerators.inputs[0])
}

func TestEventScheduler_PerformRepeatedly(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	e := &EventScheduler{
		clock:           newClock(now),
		eventGenerators: newEventCombinator(),
	}
	e.PerformRepeatedly(NewMockAction(t), nil, time.Second)

	require.Len(t, e.eventGenerators.inputs, 1)
	require.IsType(t, &periodicEventGenerator{}, e.eventGenerators.inputs[0])
}
