package simulated_time

import (
	"github.com/metamogul/timing"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewSchedulingAction(t *testing.T) {
	t.Parallel()

	mockAction := timing.NewMockAction(t)

	schedulingActionUnderTest := NewSchedulingAction(mockAction)
	require.NotNil(t, schedulingActionUnderTest)
	require.NotNil(t, schedulingActionUnderTest.Action)
	require.NotNil(t, schedulingActionUnderTest.eventLoopBlocker)

	// Ensure that the eventLoopBlocker has a wait count of 0
	schedulingActionUnderTest.eventLoopBlocker.Wait()
}
