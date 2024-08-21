package system

import (
	"context"
	"github.com/metamogul/timing"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewActionContext(t *testing.T) {
	t.Parallel()

	clock := Clock{}
	ctx := context.Background()

	actionContextUnderTest := newActionContext(ctx, clock)
	require.NotNil(t, actionContextUnderTest)
	require.NotNil(t, actionContextUnderTest.Context)
	require.NotNil(t, actionContextUnderTest.clock)
}

func TestActionContext_Clock(t *testing.T) {
	t.Parallel()

	clock := Clock{}

	actionContextUnderTest := newActionContext(context.Background(), clock)
	gotClock := actionContextUnderTest.Clock()
	require.Equal(t, clock, gotClock)
}

func TestActionContext_DoneSchedulingNewEvents(t *testing.T) {
	t.Parallel()

	actionContextUnderTest := newActionContext(context.Background(), Clock{})
	actionContextUnderTest.DoneSchedulingNewEvents()
}

func TestActionContext_Value_Clock(t *testing.T) {
	t.Parallel()

	clock := Clock{}

	actionContextUnderTest := newActionContext(context.Background(), clock)
	gotClock := actionContextUnderTest.Value(timing.ActionContextClockKey)
	require.Equal(t, clock, gotClock)
}

func TestActionContext_Value_Default(t *testing.T) {
	t.Parallel()

	actionContextUnderTest := newActionContext(context.Background(), Clock{})
	gotValue := actionContextUnderTest.Value("someNoneExistentKey")
	require.Nil(t, gotValue)
}
