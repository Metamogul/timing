package system

import (
	"context"
	"github.com/metamogul/timing"
	"sync"
	"testing"
	"time"
)

func TestEventScheduler_PerformNow(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	clock := Clock{}

	wg := &sync.WaitGroup{}

	mockAction := timing.NewMockAction(t)
	mockAction.EXPECT().
		Perform(newActionContext(ctx, clock)).
		Run(func(timing.ActionContext) { wg.Done() }).
		Once()

	eventSchedulerUnderTest := &EventScheduler{Clock: clock}
	wg.Add(1)
	eventSchedulerUnderTest.PerformNow(mockAction, ctx)
	wg.Wait()
}

func TestEventScheduler_PerformNow_cancelled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	clock := Clock{}

	eventSchedulerUnderTest := &EventScheduler{Clock: clock}
	eventSchedulerUnderTest.PerformNow(timing.NewMockAction(t), ctx)
	time.Sleep(2 * time.Millisecond)
}

func TestEventScheduler_PerformAfter(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	clock := Clock{}

	wg := &sync.WaitGroup{}

	mockAction := timing.NewMockAction(t)
	mockAction.EXPECT().
		Perform(newActionContext(ctx, clock)).
		Run(func(timing.ActionContext) { wg.Done() }).
		Once()

	eventSchedulerUnderTest := &EventScheduler{Clock: clock}
	wg.Add(1)
	eventSchedulerUnderTest.PerformAfter(mockAction, time.Millisecond, ctx)
	wg.Wait()
}

func TestEventScheduler_PerformAfter_cancelled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	clock := Clock{}

	eventSchedulerUnderTest := &EventScheduler{Clock: clock}
	eventSchedulerUnderTest.PerformAfter(timing.NewMockAction(t), time.Millisecond, ctx)
	time.Sleep(2 * time.Millisecond)
}

func TestEventScheduler_PerformRepeatedly_until(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	clock := Clock{}

	wg := &sync.WaitGroup{}

	mockAction := timing.NewMockAction(t)
	mockAction.EXPECT().
		Perform(newActionContext(ctx, clock)).
		Run(func(timing.ActionContext) { wg.Done() }).
		Twice()

	eventSchedulerUnderTest := &EventScheduler{Clock: Clock{}}
	wg.Add(2)
	eventSchedulerUnderTest.PerformRepeatedly(mockAction, ptr(clock.Now().Add(3*time.Millisecond)), time.Millisecond, ctx)
	wg.Wait()
}

func TestEventScheduler_PerformRepeatedly_indefinitely(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	clock := Clock{}

	mockAction := timing.NewMockAction(t)
	mockAction.EXPECT().
		Perform(newActionContext(ctx, clock)).
		Twice()

	eventSchedulerUnderTest := &EventScheduler{Clock: Clock{}}
	eventSchedulerUnderTest.PerformRepeatedly(mockAction, nil, time.Millisecond, ctx)
	time.Sleep(3 * time.Millisecond)
}

func TestEventScheduler_PerformRepeatedly_cancelled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	clock := Clock{}

	eventSchedulerUnderTest := &EventScheduler{Clock: Clock{}}
	eventSchedulerUnderTest.PerformRepeatedly(timing.NewMockAction(t), ptr(clock.Now().Add(3*time.Millisecond)), time.Millisecond, ctx)
	time.Sleep(2 * time.Millisecond)
}
