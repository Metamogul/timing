package simulated_time

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"time"
)

func ptr[T any](t T) *T {
	return &t
}

func Test_newPeriodicEventsStream(t *testing.T) {
	type args struct {
		action    action
		from      time.Time
		to        *time.Time
		interval  time.Duration
		scheduler eventScheduler
	}

	tests := []struct {
		name         string
		args         args
		want         *periodicEventStream
		requirePanic bool
	}{
		{
			name: "no action",
			args: args{
				action:    nil,
				from:      time.UnixMilli(0),
				to:        ptr(time.UnixMilli(1)),
				interval:  time.Second,
				scheduler: NewMockeventScheduler(t),
			},
			requirePanic: true,
		},
		{
			name: "to before from",
			args: args{
				action:    NewMockaction(t),
				from:      time.UnixMilli(1),
				to:        ptr(time.UnixMilli(0)),
				interval:  time.Second,
				scheduler: NewMockeventScheduler(t),
			},
			requirePanic: true,
		},
		{
			name: "to equals from",
			args: args{
				action:    NewMockaction(t),
				from:      time.UnixMilli(1),
				to:        ptr(time.UnixMilli(1)),
				interval:  time.Second,
				scheduler: NewMockeventScheduler(t),
			},
			requirePanic: true,
		},
		{
			name: "interval is zero",
			args: args{
				action:    NewMockaction(t),
				from:      time.UnixMilli(0),
				to:        ptr(time.UnixMilli(1)),
				interval:  0,
				scheduler: NewMockeventScheduler(t),
			},
			requirePanic: true,
		},
		{
			name: "interval is too long",
			args: args{
				action:    NewMockaction(t),
				from:      time.UnixMilli(0),
				to:        ptr(time.UnixMilli(1)),
				interval:  time.Millisecond * 2,
				scheduler: NewMockeventScheduler(t),
			},
			requirePanic: true,
		},
		{
			name: "no eventScheduler",
			args: args{
				action:    NewMockaction(t),
				from:      time.UnixMilli(0),
				to:        ptr(time.UnixMilli(1)),
				interval:  time.Second,
				scheduler: nil,
			},
			requirePanic: true,
		},
		{
			name: "success",
			args: args{
				action:    NewMockaction(t),
				from:      time.UnixMilli(0),
				to:        ptr(time.UnixMilli(1)),
				interval:  time.Second,
				scheduler: NewMockeventScheduler(t),
			},
			want: &periodicEventStream{
				action:    NewMockaction(t),
				from:      time.UnixMilli(0),
				to:        ptr(time.UnixMilli(1)),
				interval:  time.Second,
				scheduler: NewMockeventScheduler(t),

				currentEvent: &event{
					action:     NewMockaction(t),
					actionTime: time.UnixMilli(0).Add(time.Second),
					scheduler:  NewMockeventScheduler(t),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.requirePanic {
				require.Panics(t, func() {
					_ = newPeriodicEventStream(tt.args.action, tt.args.from, tt.args.to, tt.args.interval, tt.args.scheduler)
				})
				return
			}

			if got := newPeriodicEventStream(tt.args.action, tt.args.from, tt.args.to, tt.args.interval, tt.args.scheduler); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newPeriodicEventStream() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_periodicEventStream_popNextEvent(t *testing.T) {
	type fields struct {
		action       action
		from         time.Time
		to           *time.Time
		interval     time.Duration
		scheduler    func() eventScheduler
		currentEvent *event
	}
	tests := []struct {
		name          string
		fields        fields
		want          *event
		requirePanic  bool
		requireClosed bool
	}{
		{
			name: "already closed",
			fields: fields{
				to: ptr(time.UnixMilli(1)),
				scheduler: func() eventScheduler {
					mockedScheduler := NewMockeventScheduler(t)
					mockedScheduler.EXPECT().
						Now().
						Return(time.UnixMilli(2)).
						Once()

					return mockedScheduler
				},
			},
			requirePanic: true,
		},
		{
			name: "success, current event before scheduler.Now()",
			fields: fields{
				action:   NewMockaction(t),
				from:     time.UnixMilli(0),
				to:       nil,
				interval: time.Second,
				scheduler: func() eventScheduler {
					mockedScheduler := NewMockeventScheduler(t)
					mockedScheduler.EXPECT().
						Now().
						Return(time.UnixMilli(1500))

					return mockedScheduler
				},
				currentEvent: newEvent(NewMockaction(t), time.UnixMilli(0).Add(time.Second), NewMockeventScheduler(t)),
			},
			want: newEvent(NewMockaction(t), time.UnixMilli(0).Add(time.Second), NewMockeventScheduler(t)),
		},
		{
			name:   "success, current event at scheduler.Now()",
			fields: fields{},
			want:   nil,
		},
		{
			name:   "success, current event after scheduler.Now()",
			fields: fields{},
			want:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &periodicEventStream{
				action:       tt.fields.action,
				from:         tt.fields.from,
				to:           tt.fields.to,
				interval:     tt.fields.interval,
				scheduler:    tt.fields.scheduler(),
				currentEvent: tt.fields.currentEvent,
			}

			if tt.requirePanic {
				require.Panics(t, func() {
					_ = p.popNextEvent()
				})
				return
			}

			if got := p.popNextEvent(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("popNextEvent() = %v, want %v", got, tt.want)
			}

			if tt.requireClosed {
				require.True(t, p.closed())
			} else {
				require.False(t, p.closed())
			}
		})
	}
}
