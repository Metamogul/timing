package simulated_time

import (
	"github.com/metamogul/timing"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"time"
)

func Test_newEvent(t *testing.T) {
	t.Parallel()

	type args struct {
		action     timing.Action
		actionTime time.Time
	}

	tests := []struct {
		name         string
		args         args
		want         *Event
		requirePanic bool
	}{
		{
			name: "no Action",
			args: args{
				action:     nil,
				actionTime: time.Time{},
			},
			requirePanic: true,
		},
		{
			name: "success",
			args: args{
				action:     timing.NewMockAction(t),
				actionTime: time.Time{},
			},
			want: &Event{
				Action: timing.NewMockAction(t),
				Time:   time.Time{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.requirePanic {
				require.Panics(t, func() {
					_ = NewEvent(tt.args.action, tt.args.actionTime)
				})
				return
			}

			if got := NewEvent(tt.args.action, tt.args.actionTime); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_event_perform(t *testing.T) {
	t.Parallel()

	clockArg := newClock(time.Now())

	e := &Event{
		Action: func() timing.Action {
			mockedAction := timing.NewMockAction(t)
			mockedAction.EXPECT().
				Perform(clockArg).
				Once()

			return mockedAction
		}(),
		Time: time.Time{},
	}

	e.Perform(clockArg)
}
