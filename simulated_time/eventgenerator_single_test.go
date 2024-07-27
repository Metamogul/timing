package simulated_time

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"time"
)

func Test_newSingleEventGenerator(t *testing.T) {
	t.Parallel()

	type args struct {
		action     action
		actionTime time.Time
	}

	tests := []struct {
		name         string
		args         args
		want         *singleEventGenerator
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
			want: &singleEventGenerator{
				event: &event{
					action: NewMockaction(t),
					Time:   time.Time{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.requirePanic {
				require.Panics(t, func() {
					_ = newSingleEventGenerator(tt.args.action, tt.args.actionTime)
				})
				return
			}

			if got := newSingleEventGenerator(tt.args.action, tt.args.actionTime); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newSingleEventGenerator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_singleEventStream_pop(t *testing.T) {
	t.Parallel()

	type fields struct {
		event *event
	}

	tests := []struct {
		name         string
		fields       fields
		want         *event
		requirePanic bool
	}{
		{
			name: "already finished",
			fields: fields{
				event: nil,
			},
			requirePanic: true,
		},
		{
			name: "success",
			fields: fields{
				event: newEvent(NewMockaction(t), time.Time{}),
			},
			want: newEvent(NewMockaction(t), time.Time{}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &singleEventGenerator{
				event: tt.fields.event,
			}

			if tt.requirePanic {
				require.Panics(t, func() {
					_ = s.pop()
				})
				return
			}

			if got := s.pop(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pop() = %v, want %v", got, tt.want)
			}

			if tt.want != nil {
				require.True(t, s.finished())
			} else {
				require.False(t, s.finished())
			}
		})
	}
}

func Test_singleEventStream_peek(t *testing.T) {
	t.Parallel()

	type fields struct {
		event *event
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
				event: nil,
			},
			requirePanic: true,
		},
		{
			name: "success",
			fields: fields{
				event: newEvent(NewMockaction(t), time.Time{}),
			},
			want: *newEvent(NewMockaction(t), time.Time{}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &singleEventGenerator{
				event: tt.fields.event,
			}

			if tt.requirePanic {
				require.Panics(t, func() {
					_ = s.peek()
				})
				return
			}

			if got := s.peek(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("peek() = %v, want %v", got, tt.want)
			}

			require.False(t, s.finished())
		})
	}
}

func Test_singleEventStream_finished(t *testing.T) {
	t.Parallel()

	type fields struct {
		event *event
	}

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "finished",
			fields: fields{
				event: nil,
			},
			want: true,
		},
		{
			name: "not finished",
			fields: fields{
				event: newEvent(NewMockaction(t), time.Time{}),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &singleEventGenerator{
				event: tt.fields.event,
			}

			if got := s.finished(); got != tt.want {
				t.Errorf("finished() = %v, want %v", got, tt.want)
			}
		})
	}
}
