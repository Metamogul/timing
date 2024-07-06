package simulated_time

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"sync"
	"testing"
	"time"
)

func Test_newEvent(t *testing.T) {
	type args struct {
		action     action
		actionTime time.Time
		scheduler  eventScheduler
	}

	tests := []struct {
		name         string
		args         args
		want         *event
		requirePanic bool
	}{
		{
			name: "no action",
			args: args{
				action:     nil,
				actionTime: time.Time{},
				scheduler:  NewMockeventScheduler(t),
			},
			requirePanic: true,
		},
		{
			name: "no eventScheduler",
			args: args{
				action:     NewMockaction(t),
				actionTime: time.Time{},
				scheduler:  nil,
			},
			requirePanic: true,
		},
		{
			name: "success",
			args: args{
				action:     NewMockaction(t),
				actionTime: time.Time{},
				scheduler:  NewMockeventScheduler(t),
			},
			want: &event{
				action:     NewMockaction(t),
				actionTime: time.Time{},
				scheduler:  NewMockeventScheduler(t),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.requirePanic {
				require.Panics(t, func() {
					_ = newEvent(tt.args.action, tt.args.actionTime, tt.args.scheduler)
				})
				return
			}

			if got := newEvent(tt.args.action, tt.args.actionTime, tt.args.scheduler); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_event_perform(t *testing.T) {
	type fields struct {
		action     func() action
		actionTime time.Time
		scheduler  func() eventScheduler
	}

	tests := []struct {
		name         string
		fields       fields
		requirePanic bool
	}{
		{
			name: "success",
			fields: fields{
				action: func() action {
					mockedAction := NewMockaction(t)
					mockedAction.EXPECT().
						perform().
						Once()

					return mockedAction
				},
				actionTime: time.UnixMilli(0),
				scheduler: func() eventScheduler {
					mockedScheduler := NewMockeventScheduler(t)

					mockedScheduler.EXPECT().
						Now().
						Return(time.UnixMilli(1)).
						Once()

					return mockedScheduler
				},
			},
		},
		{
			name: "panic on perform ahead of time",
			fields: fields{
				action: func() action {
					return NewMockaction(t)
				},
				actionTime: time.UnixMilli(2),
				scheduler: func() eventScheduler {
					mockedScheduler := NewMockeventScheduler(t)

					mockedScheduler.EXPECT().
						Now().
						Return(time.UnixMilli(1)).
						Once()

					return mockedScheduler
				},
			},
			requirePanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				action:     tt.fields.action(),
				actionTime: tt.fields.actionTime,
				scheduler:  tt.fields.scheduler(),
			}

			if tt.requirePanic {
				require.Panics(t, e.perform)
				return

			}

			e.perform()
		})
	}
}

func Test_event_performAsync(t *testing.T) {
	type fields struct {
		action     func() action
		actionTime time.Time
		scheduler  func() eventScheduler
	}

	wg := sync.WaitGroup{}

	tests := []struct {
		name         string
		fields       fields
		requirePanic bool
	}{
		{
			name: "success",
			fields: fields{
				action: func() action {
					mockedAction := NewMockaction(t)

					mockedAction.EXPECT().
						perform().
						Once()

					return mockedAction
				},
				actionTime: time.UnixMilli(0),
				scheduler: func() eventScheduler {
					mockedScheduler := NewMockeventScheduler(t)

					mockedScheduler.EXPECT().
						Now().
						Return(time.UnixMilli(1)).
						Once()

					mockedScheduler.EXPECT().
						eventCompletionWaitGroupAdd(1).
						Run(func(delta int) {
							wg.Add(delta)
						}).
						Once()

					mockedScheduler.EXPECT().
						eventCompletionWaitGroupDone().
						Run(func() {
							wg.Done()
						}).
						Once()

					return mockedScheduler
				},
			},
		},
		{
			name: "panic on perform ahead of time",
			fields: fields{
				action: func() action {
					return NewMockaction(t)
				},
				actionTime: time.UnixMilli(2),
				scheduler: func() eventScheduler {
					mockedScheduler := NewMockeventScheduler(t)
					mockedScheduler.EXPECT().
						Now().
						Return(time.UnixMilli(1)).
						Once()

					return mockedScheduler
				},
			},
			requirePanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				action:     tt.fields.action(),
				actionTime: tt.fields.actionTime,
				scheduler:  tt.fields.scheduler(),
			}

			if tt.requirePanic {
				require.Panics(t, e.perform)
				return
			}

			e.performAsync()
			wg.Wait()
		})
	}
}
