package simulated_time

import (
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
	}

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
			},
			requirePanic: true,
		},
		{
			name: "to before from",
			args: args{
				action:   NewMockAction(t),
				from:     time.Time{}.Add(time.Second),
				to:       ptr(time.Time{}),
				interval: time.Second,
			},
			requirePanic: true,
		},
		{
			name: "to equals from",
			args: args{
				action:   NewMockAction(t),
				from:     time.Time{}.Add(time.Second),
				to:       ptr(time.Time{}.Add(time.Second)),
				interval: time.Second,
			},
			requirePanic: true,
		},
		{
			name: "interval is zero",
			args: args{
				action:   NewMockAction(t),
				from:     time.Time{},
				to:       ptr(time.Time{}.Add(time.Second)),
				interval: 0,
			},
			requirePanic: true,
		},
		{
			name: "interval is too long",
			args: args{
				action:   NewMockAction(t),
				from:     time.Time{},
				to:       ptr(time.Time{}.Add(time.Second)),
				interval: time.Second * 2,
			},
			requirePanic: true,
		},
		{
			name: "success",
			args: args{
				action:   NewMockAction(t),
				from:     time.Time{},
				to:       ptr(time.Time{}.Add(2 * time.Second)),
				interval: time.Second,
			},
			want: &periodicEventGenerator{
				action:   NewMockAction(t),
				from:     time.Time{},
				to:       ptr(time.Time{}.Add(2 * time.Second)),
				interval: time.Second,

				currentEvent: &event{
					Action: NewMockAction(t),
					Time:   time.Time{}.Add(time.Second),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.requirePanic {
				require.Panics(t, func() {
					_ = newPeriodicEventGenerator(tt.args.action, tt.args.from, tt.args.to, tt.args.interval)
				})
				return
			}

			if got := newPeriodicEventGenerator(tt.args.action, tt.args.from, tt.args.to, tt.args.interval); !reflect.DeepEqual(got, tt.want) {
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
		currentEvent *event
	}

	tests := []struct {
		name            string
		fields          fields
		want            *event
		requirePanic    bool
		requireFinished bool
	}{
		{
			name: "already finished",
			fields: fields{
				action:       NewMockAction(t),
				from:         time.Time{},
				to:           ptr(time.Time{}.Add(time.Minute)),
				interval:     10 * time.Second,
				currentEvent: newEvent(NewMockAction(t), time.Time{}.Add(55*time.Second)),
			},
			requirePanic: true,
		},
		{
			name: "success, not finished 1",
			fields: fields{
				action:       NewMockAction(t),
				from:         time.Time{},
				to:           nil,
				interval:     time.Second,
				currentEvent: newEvent(NewMockAction(t), time.Time{}.Add(time.Second)),
			},
			want: newEvent(NewMockAction(t), time.Time{}.Add(time.Second)),
		},
		{
			name: "success, not finished 2",
			fields: fields{
				action:       NewMockAction(t),
				from:         time.Time{},
				to:           ptr(time.Time{}.Add(time.Minute)),
				interval:     10 * time.Second,
				currentEvent: newEvent(NewMockAction(t), time.Time{}.Add(40*time.Second)),
			},
			want: newEvent(NewMockAction(t), time.Time{}.Add(40*time.Second)),
		},
		{
			name: "success, finished",
			fields: fields{
				action:       NewMockAction(t),
				from:         time.Time{},
				to:           ptr(time.Time{}.Add(time.Minute)),
				interval:     10 * time.Second,
				currentEvent: newEvent(NewMockAction(t), time.Time{}.Add(50*time.Second)),
			},
			want:            newEvent(NewMockAction(t), time.Time{}.Add(50*time.Second)),
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
			}

			if tt.requirePanic {
				require.Panics(t, func() {
					_ = p.pop()
				})
				return
			}

			if got := p.pop(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pop() = %v, want %v", got, tt.want)
			}

			if tt.requireFinished {
				require.True(t, p.finished())
			} else {
				require.False(t, p.finished())
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
		currentEvent *event
	}

	tests := []struct {
		name         string
		fields       fields
		want         event
		requirePanic bool
	}{
		{
			name: "already finished",
			fields: fields{
				action:       NewMockAction(t),
				from:         time.Time{},
				to:           ptr(time.Time{}.Add(time.Minute)),
				interval:     10 * time.Second,
				currentEvent: newEvent(NewMockAction(t), time.Time{}.Add(55*time.Second)),
			},
			requirePanic: true,
		},
		{
			name: "success",
			fields: fields{
				action:       NewMockAction(t),
				from:         time.Time{},
				to:           nil,
				interval:     time.Second,
				currentEvent: newEvent(NewMockAction(t), time.Time{}.Add(time.Second)),
			},
			want: *newEvent(NewMockAction(t), time.Time{}.Add(time.Second)),
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
			}

			if tt.requirePanic {
				require.Panics(t, func() {
					_ = p.peek()
				})
				return
			}

			if got := p.peek(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("peek() = %v, want %v", got, tt.want)
			}

			require.False(t, p.finished())
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
		currentEvent *event
	}

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "never finished",
			fields: fields{
				action:       NewMockAction(t),
				from:         time.Time{},
				to:           nil,
				interval:     0,
				currentEvent: newEvent(NewMockAction(t), time.Time{}),
			},
			want: false,
		},
		{
			name: "not finished yet",
			fields: fields{
				action:       NewMockAction(t),
				from:         time.Time{},
				to:           ptr(time.Time{}.Add(time.Minute)),
				interval:     10 * time.Second,
				currentEvent: newEvent(NewMockAction(t), time.Time{}.Add(45*time.Second)),
			},
			want: false,
		},
		{
			name: "finished",
			fields: fields{
				action:       NewMockAction(t),
				from:         time.Time{},
				to:           ptr(time.Time{}.Add(time.Minute)),
				interval:     10 * time.Second,
				currentEvent: newEvent(NewMockAction(t), time.Time{}.Add(55*time.Second)),
			},
			want: true,
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
			}

			require.Equal(t, tt.want, p.finished())
		})
	}
}
