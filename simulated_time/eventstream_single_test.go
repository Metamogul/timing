package simulated_time

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"time"
)

func Test_newSingleEventStream(t *testing.T) {
	type args struct {
		action     action
		actionTime time.Time
		scheduler  eventScheduler
	}

	tests := []struct {
		name         string
		args         args
		want         *singleEventStream
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
			want: &singleEventStream{
				scheduler: NewMockeventScheduler(t),
				event: &event{
					action:     NewMockaction(t),
					actionTime: time.Time{},
					scheduler:  NewMockeventScheduler(t),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.requirePanic {
				require.Panics(t, func() {
					_ = newSingleEventStream(tt.args.action, tt.args.actionTime, tt.args.scheduler)
				})
				return
			}

			if got := newSingleEventStream(tt.args.action, tt.args.actionTime, tt.args.scheduler); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newSingleEventStream() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_singleEventStream_popNextEvent(t *testing.T) {
	type fields struct {
		scheduler func() eventScheduler
		event     *event
	}

	tests := []struct {
		name         string
		fields       fields
		want         *event
		requirePanic bool
	}{
		{
			name: "already closed",
			fields: fields{
				scheduler: func() eventScheduler { return NewMockeventScheduler(t) },
				event:     nil,
			},
			requirePanic: true,
		},
		{
			name: "success, event before scheduler.Now()",
			fields: fields{
				scheduler: func() eventScheduler {
					mockedScheduler := NewMockeventScheduler(t)
					mockedScheduler.EXPECT().
						Now().
						Return(time.UnixMilli(1))

					return mockedScheduler
				},
				event: newEvent(NewMockaction(t), time.UnixMilli(0), NewMockeventScheduler(t)),
			},
			want: newEvent(NewMockaction(t), time.UnixMilli(0), NewMockeventScheduler(t)),
		},
		{
			name: "success, event at scheduler.Now()",
			fields: fields{
				scheduler: func() eventScheduler {
					mockedScheduler := NewMockeventScheduler(t)
					mockedScheduler.EXPECT().
						Now().
						Return(time.UnixMilli(1))

					return mockedScheduler
				},
				event: newEvent(NewMockaction(t), time.UnixMilli(1), NewMockeventScheduler(t)),
			},
			want: nil,
		},
		{
			name: "success, event after scheduler.Now()",
			fields: fields{
				scheduler: func() eventScheduler {
					mockedScheduler := NewMockeventScheduler(t)
					mockedScheduler.EXPECT().
						Now().
						Return(time.UnixMilli(1))

					return mockedScheduler
				},
				event: newEvent(NewMockaction(t), time.UnixMilli(2), NewMockeventScheduler(t)),
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &singleEventStream{
				scheduler: tt.fields.scheduler(),
				event:     tt.fields.event,
			}

			if tt.requirePanic {
				require.Panics(t, func() {
					_ = s.popNextEvent()
				})
				return
			}

			if got := s.popNextEvent(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("popNextEvent() = %v, want %v", got, tt.want)
			}

			if tt.want != nil {
				require.True(t, s.closed())
			} else {
				require.False(t, s.closed())
			}
		})
	}
}

func Test_singleEventStream_peakNextTime(t *testing.T) {
	type fields struct {
		event *event
	}

	tests := []struct {
		name         string
		fields       fields
		want         time.Time
		requirePanic bool
	}{
		{
			name: "already closed",
			fields: fields{
				event: nil,
			},
			requirePanic: true,
		},
		{
			name: "success",
			fields: fields{
				event: newEvent(NewMockaction(t), time.Time{}, NewMockeventScheduler(t)),
			},
			want: time.Time{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &singleEventStream{
				event: tt.fields.event,
			}

			if tt.requirePanic {
				require.Panics(t, func() {
					_ = s.peakNextTime()
				})
				return
			}

			if got := s.peakNextTime(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("peakNextTime() = %v, want %v", got, tt.want)
			}

			require.False(t, s.closed())
		})
	}
}

func Test_singleEventStream_closed(t *testing.T) {
	type fields struct {
		event *event
	}

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "closed",
			fields: fields{
				event: nil,
			},
			want: true,
		},
		{
			name: "not closed",
			fields: fields{
				event: newEvent(NewMockaction(t), time.Time{}, NewMockeventScheduler(t)),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &singleEventStream{
				event: tt.fields.event,
			}

			if got := s.closed(); got != tt.want {
				t.Errorf("closed() = %v, want %v", got, tt.want)
			}
		})
	}
}
