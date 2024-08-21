package simulated_time

import (
	"context"
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
		ctx        context.Context
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
				ctx:        context.Background(),
			},
			want: &Event{
				Action:  timing.NewMockAction(t),
				Time:    time.Time{},
				Context: context.Background(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.requirePanic {
				require.Panics(t, func() {
					_ = NewEvent(tt.args.action, tt.args.actionTime, tt.args.ctx)
				})
				return
			}

			if got := NewEvent(tt.args.action, tt.args.actionTime, tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_event_perform(t *testing.T) {
	t.Parallel()

	actionContextArg := newActionContext(context.Background(), newClock(time.Now()), nil)

	e := &Event{
		Action: func() timing.Action {
			mockedAction := timing.NewMockAction(t)
			mockedAction.EXPECT().
				Perform(actionContextArg).
				Once()

			return mockedAction
		}(),
		Time: time.Time{},
	}

	e.Perform(actionContextArg)
}
