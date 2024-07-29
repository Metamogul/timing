package system

import (
	"context"
	"github.com/metamogul/timing"
	"sync"
	"testing"
	"time"
)

func TestEventScheduler_PerformAfter(t *testing.T) {
	t.Parallel()

	clock := Clock{}
	s := &EventScheduler{Clock: clock}

	wg := &sync.WaitGroup{}

	mockAction := timing.NewMockAction(t)
	mockAction.EXPECT().
		Perform(clock).
		Run(func(timing.Clock) { wg.Done() }).
		Once()

	wg.Add(1)
	s.PerformAfter(mockAction, time.Millisecond, context.Background())
	wg.Wait()
}

func TestEventScheduler_PerformAfter_cancelled(t *testing.T) {
	t.Parallel()

	clock := Clock{}
	s := &EventScheduler{Clock: clock}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	s.PerformAfter(timing.NewMockAction(t), time.Millisecond, ctx)
	time.Sleep(2 * time.Millisecond)
}

func TestEventScheduler_PerformRepeatedly_until(t *testing.T) {
	t.Parallel()

	clock := Clock{}
	s := &EventScheduler{Clock: Clock{}}

	wg := &sync.WaitGroup{}

	mockAction := timing.NewMockAction(t)
	mockAction.EXPECT().
		Perform(clock).
		Run(func(timing.Clock) { wg.Done() }).
		Twice()

	wg.Add(2)
	s.PerformRepeatedly(mockAction, ptr(clock.Now().Add(3*time.Millisecond)), time.Millisecond, context.Background())
	wg.Wait()
}

func TestEventScheduler_PerformRepeatedly_indefinitely(t *testing.T) {
	t.Parallel()

	clock := Clock{}
	s := &EventScheduler{Clock: Clock{}}

	mockAction := timing.NewMockAction(t)
	mockAction.EXPECT().
		Perform(clock).
		Twice()

	s.PerformRepeatedly(mockAction, nil, time.Millisecond, context.Background())
	time.Sleep(3 * time.Millisecond)
}

func TestEventScheduler_PerformRepeatedly_cancelled(t *testing.T) {
	t.Parallel()

	clock := Clock{}
	s := &EventScheduler{Clock: Clock{}}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	s.PerformRepeatedly(timing.NewMockAction(t), ptr(clock.Now().Add(3*time.Millisecond)), time.Millisecond, ctx)
	time.Sleep(2 * time.Millisecond)
}
