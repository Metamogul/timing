package simulated_time

import (
	"context"
	"github.com/metamogul/timing"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"time"
)

func Test_newPeriodicEventGenerator(t *testing.T) {
	t.Parallel()

	type args struct {
		action   timing.Action
		from     time.Time
		to       *time.Time
		interval time.Duration
		ctx      context.Context
	}

	ctx := context.Background()

	tests := []struct {
		name         string
		args         args
		want         *periodicEventGenerator
		requirePanic bool
	}{
		{
			name: "no Action",
			args: args{
				action:   nil,
				from:     time.Time{},
				to:       ptr(time.Time{}.Add(time.Second)),
				interval: time.Second,
				ctx:      ctx,
			},
			requirePanic: true,
		},
		{
			name: "to before from",
			args: args{
				action:   timing.NewMockAction(t),
				from:     time.Time{}.Add(time.Second),
				to:       ptr(time.Time{}),
				interval: time.Second,
				ctx:      ctx,
			},
			requirePanic: true,
		},
		{
			name: "to equals from",
			args: args{
				action:   timing.NewMockAction(t),
				from:     time.Time{}.Add(time.Second),
				to:       ptr(time.Time{}.Add(time.Second)),
				interval: time.Second,
				ctx:      ctx,
			},
			requirePanic: true,
		},
		{
			name: "interval is zero",
			args: args{
				action:   timing.NewMockAction(t),
				from:     time.Time{},
				to:       ptr(time.Time{}.Add(time.Second)),
				interval: 0,
				ctx:      ctx,
			},
			requirePanic: true,
		},
		{
			name: "interval is too long",
			args: args{
				action:   timing.NewMockAction(t),
				from:     time.Time{},
				to:       ptr(time.Time{}.Add(time.Second)),
				interval: time.Second * 2,
				ctx:      ctx,
			},
			requirePanic: true,
		},
		{
			name: "success",
			args: args{
				action:   timing.NewMockAction(t),
				from:     time.Time{},
				to:       ptr(time.Time{}.Add(2 * time.Second)),
				interval: time.Second,
				ctx:      ctx,
			},
			want: &periodicEventGenerator{
				action:   timing.NewMockAction(t),
				from:     time.Time{},
				to:       ptr(time.Time{}.Add(2 * time.Second)),
				interval: time.Second,
				currentEvent: &Event{
					Action:  timing.NewMockAction(t),
					Time:    time.Time{}.Add(time.Second),
					Context: ctx,
				},
				ctx: ctx,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.requirePanic {
				require.Panics(t, func() {
					_ = newPeriodicEventGenerator(tt.args.action, tt.args.from, tt.args.to, tt.args.interval, tt.args.ctx)
				})
				return
			}

			if got := newPeriodicEventGenerator(tt.args.action, tt.args.from, tt.args.to, tt.args.interval, tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newPeriodicEventGenerator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_periodicEventGenerator_pop(t *testing.T) {
	t.Parallel()

	type fields struct {
		action       timing.Action
		from         time.Time
		to           *time.Time
		interval     time.Duration
		currentEvent *Event
		ctx          context.Context
	}

	ctx := context.Background()

	tests := []struct {
		name            string
		fields          fields
		want            *Event
		requirePanic    bool
		requireFinished bool
	}{
		{
			name: "already finished",
			fields: fields{
				action:       timing.NewMockAction(t),
				from:         time.Time{},
				to:           ptr(time.Time{}.Add(time.Minute)),
				interval:     10 * time.Second,
				currentEvent: NewEvent(timing.NewMockAction(t), time.Time{}.Add(55*time.Second), ctx),
				ctx:          context.Background(),
			},
			requirePanic: true,
		},
		{
			name: "success, not finished 1",
			fields: fields{
				action:       timing.NewMockAction(t),
				from:         time.Time{},
				to:           nil,
				interval:     time.Second,
				currentEvent: NewEvent(timing.NewMockAction(t), time.Time{}.Add(time.Second), ctx),
				ctx:          context.Background(),
			},
			want: NewEvent(timing.NewMockAction(t), time.Time{}.Add(time.Second), ctx),
		},
		{
			name: "success, not finished 2",
			fields: fields{
				action:       timing.NewMockAction(t),
				from:         time.Time{},
				to:           ptr(time.Time{}.Add(time.Minute)),
				interval:     10 * time.Second,
				currentEvent: NewEvent(timing.NewMockAction(t), time.Time{}.Add(40*time.Second), ctx),
				ctx:          context.Background(),
			},
			want: NewEvent(timing.NewMockAction(t), time.Time{}.Add(40*time.Second), ctx),
		},
		{
			name: "success, finished",
			fields: fields{
				action:       timing.NewMockAction(t),
				from:         time.Time{},
				to:           ptr(time.Time{}.Add(time.Minute)),
				interval:     10 * time.Second,
				currentEvent: NewEvent(timing.NewMockAction(t), time.Time{}.Add(50*time.Second), ctx),
				ctx:          context.Background(),
			},
			want:            NewEvent(timing.NewMockAction(t), time.Time{}.Add(50*time.Second), ctx),
			requireFinished: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := &periodicEventGenerator{
				action:       tt.fields.action,
				from:         tt.fields.from,
				to:           tt.fields.to,
				interval:     tt.fields.interval,
				currentEvent: tt.fields.currentEvent,
				ctx:          tt.fields.ctx,
			}

			if tt.requirePanic {
				require.Panics(t, func() {
					_ = p.Pop()
				})
				return
			}

			if got := p.Pop(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pop() = %v, want %v", got, tt.want)
			}

			if tt.requireFinished {
				require.True(t, p.Finished())
			} else {
				require.False(t, p.Finished())
			}
		})
	}
}

func Test_periodicEventGenerator_peek(t *testing.T) {
	t.Parallel()

	type fields struct {
		action       timing.Action
		from         time.Time
		to           *time.Time
		interval     time.Duration
		currentEvent *Event
		ctx          context.Context
	}

	ctx := context.Background()

	tests := []struct {
		name         string
		fields       fields
		want         Event
		requirePanic bool
	}{
		{
			name: "already finished",
			fields: fields{
				action:       timing.NewMockAction(t),
				from:         time.Time{},
				to:           ptr(time.Time{}.Add(time.Minute)),
				interval:     10 * time.Second,
				currentEvent: NewEvent(timing.NewMockAction(t), time.Time{}.Add(55*time.Second), ctx),
				ctx:          context.Background(),
			},
			requirePanic: true,
		},
		{
			name: "success",
			fields: fields{
				action:       timing.NewMockAction(t),
				from:         time.Time{},
				to:           nil,
				interval:     time.Second,
				currentEvent: NewEvent(timing.NewMockAction(t), time.Time{}.Add(time.Second), ctx),
				ctx:          context.Background(),
			},
			want: *NewEvent(timing.NewMockAction(t), time.Time{}.Add(time.Second), ctx),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := &periodicEventGenerator{
				action:       tt.fields.action,
				from:         tt.fields.from,
				to:           tt.fields.to,
				interval:     tt.fields.interval,
				currentEvent: tt.fields.currentEvent,
				ctx:          tt.fields.ctx,
			}

			if tt.requirePanic {
				require.Panics(t, func() {
					_ = p.Peek()
				})
				return
			}

			if got := p.Peek(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Peek() = %v, want %v", got, tt.want)
			}

			require.False(t, p.Finished())
		})
	}
}

func Test_periodicEventGenerator_finished(t *testing.T) {
	t.Parallel()

	type fields struct {
		action       timing.Action
		from         time.Time
		to           *time.Time
		interval     time.Duration
		currentEvent *Event
		ctx          context.Context
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "context is done",
			fields: fields{
				action:       timing.NewMockAction(t),
				from:         time.Time{},
				to:           ptr(time.Time{}.Add(time.Minute)),
				interval:     10 * time.Second,
				currentEvent: NewEvent(timing.NewMockAction(t), time.Time{}.Add(45*time.Second), ctx),
				ctx:          ctx,
			},
			want: true,
		},
		{
			name: "to is nil",
			fields: fields{
				action:       timing.NewMockAction(t),
				from:         time.Time{},
				to:           nil,
				interval:     0,
				currentEvent: NewEvent(timing.NewMockAction(t), time.Time{}, ctx),
				ctx:          context.Background(),
			},
			want: false,
		},
		{
			name: "to is set, finished",
			fields: fields{
				action:       timing.NewMockAction(t),
				from:         time.Time{},
				to:           ptr(time.Time{}.Add(time.Minute)),
				interval:     10 * time.Second,
				currentEvent: NewEvent(timing.NewMockAction(t), time.Time{}.Add(55*time.Second), ctx),
				ctx:          context.Background(),
			},
			want: true,
		},
		{
			name: "to is set, not finished yet",
			fields: fields{
				action:       timing.NewMockAction(t),
				from:         time.Time{},
				to:           ptr(time.Time{}.Add(time.Minute)),
				interval:     10 * time.Second,
				currentEvent: NewEvent(timing.NewMockAction(t), time.Time{}.Add(45*time.Second), ctx),
				ctx:          context.Background(),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := &periodicEventGenerator{
				action:       tt.fields.action,
				from:         tt.fields.from,
				to:           tt.fields.to,
				interval:     tt.fields.interval,
				currentEvent: tt.fields.currentEvent,
				ctx:          tt.fields.ctx,
			}

			require.Equal(t, tt.want, p.Finished())
		})
	}
}
