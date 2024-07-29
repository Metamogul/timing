package simulated_time

import (
	"context"
	"github.com/metamogul/timing"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"time"
)

func Test_newSingleEventGenerator(t *testing.T) {
	t.Parallel()

	type args struct {
		action     timing.Action
		actionTime time.Time
		ctx        context.Context
	}

	tests := []struct {
		name         string
		args         args
		want         *singleEventGenerator
		requirePanic bool
	}{
		{
			name: "no Action",
			args: args{
				action:     nil,
				actionTime: time.Time{},
				ctx:        context.Background(),
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
			want: &singleEventGenerator{
				event: &event{
					Action: timing.NewMockAction(t),
					Time:   time.Time{},
				},
				ctx: context.Background(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.requirePanic {
				require.Panics(t, func() {
					_ = newSingleEventGenerator(tt.args.action, tt.args.actionTime, tt.args.ctx)
				})
				return
			}

			if got := newSingleEventGenerator(tt.args.action, tt.args.actionTime, tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newSingleEventGenerator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_singleEventStream_pop(t *testing.T) {
	t.Parallel()

	type fields struct {
		event *event
		ctx   context.Context
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
				ctx:   context.Background(),
			},
			requirePanic: true,
		},
		{
			name: "success",
			fields: fields{
				event: newEvent(timing.NewMockAction(t), time.Time{}),
				ctx:   context.Background(),
			},
			want: newEvent(timing.NewMockAction(t), time.Time{}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &singleEventGenerator{
				event: tt.fields.event,
				ctx:   tt.fields.ctx,
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
		ctx   context.Context
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
				ctx:   context.Background(),
			},
			requirePanic: true,
		},
		{
			name: "success",
			fields: fields{
				event: newEvent(timing.NewMockAction(t), time.Time{}),
				ctx:   context.Background(),
			},
			want: *newEvent(timing.NewMockAction(t), time.Time{}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &singleEventGenerator{
				event: tt.fields.event,
				ctx:   tt.fields.ctx,
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
		ctx   context.Context
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "no event",
			fields: fields{
				event: nil,
				ctx:   context.Background(),
			},
			want: true,
		},
		{
			name: "context is done",
			fields: fields{
				event: newEvent(timing.NewMockAction(t), time.Time{}),
				ctx:   ctx,
			},
			want: true,
		},
		{
			name: "not finished",
			fields: fields{
				event: newEvent(timing.NewMockAction(t), time.Time{}),
				ctx:   context.Background(),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &singleEventGenerator{
				event: tt.fields.event,
				ctx:   tt.fields.ctx,
			}

			if got := s.finished(); got != tt.want {
				t.Errorf("finished() = %v, want %v", got, tt.want)
			}
		})
	}
}
