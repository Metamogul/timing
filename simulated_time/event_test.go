package simulated_time

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"time"
)

func Test_newEvent(t *testing.T) {
	t.Parallel()

	type args struct {
		action     action
		actionTime time.Time
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
			},
			requirePanic: true,
		},
		{
			name: "success",
			args: args{
				action:     NewMockaction(t),
				actionTime: time.Time{},
			},
			want: &event{
				action: NewMockaction(t),
				Time:   time.Time{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.requirePanic {
				require.Panics(t, func() {
					_ = newEvent(tt.args.action, tt.args.actionTime)
				})
				return
			}

			if got := newEvent(tt.args.action, tt.args.actionTime); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_event_perform(t *testing.T) {
	t.Parallel()

	e := &event{
		action: func() action {
			mockedAction := NewMockaction(t)
			mockedAction.EXPECT().
				perform().
				Once()

			return mockedAction
		}(),
		Time: time.Time{},
	}

	e.perform()
}
