package simulated_time

import (
	"github.com/metamogul/timing"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewSchedulingAction(t *testing.T) {
	t.Parallel()

	mockAction := timing.NewMockAction(t)

	schedulingActionUnderTest := NewSchedulingAction(mockAction)
	require.NotNil(t, schedulingActionUnderTest)
	require.NotNil(t, schedulingActionUnderTest.Action)
	require.NotNil(t, schedulingActionUnderTest.eventLoopBlocker)

	// Ensure that the eventLoopBlocker has a wait count of 1
	go func() {
		schedulingActionUnderTest.eventLoopBlocker.Done()
	}()
	schedulingActionUnderTest.WaitForEventSchedulingCompletion()
}

func TestSchedulingAction_WaitForEventSchedulingCompletion(t *testing.T) {
	t.Parallel()

	mockAction := timing.NewMockAction(t)
	schedulingActionUnderTest := NewSchedulingAction(mockAction)

	go func() {
		time.Sleep(time.Millisecond * 100)
		schedulingActionUnderTest.eventLoopBlocker.Done()
	}()

	time1 := time.Now()
	schedulingActionUnderTest.WaitForEventSchedulingCompletion()
	time2 := time.Now()

	require.Greater(t, time2.Sub(time1), time.Millisecond*99)
}
