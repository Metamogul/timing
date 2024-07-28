package system

import (
	"github.com/metamogul/timing"
	"reflect"
	"testing"
	"time"
)

func TestClock_Now(t *testing.T) {
	tests := []struct {
		name string
		want time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Clock{}
			if got := s.Now(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Now() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventScheduler_PerformAfter(t *testing.T) {
	type fields struct {
		Clock Clock
	}
	type args struct {
		duration time.Duration
		action   timing.Action
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &EventScheduler{
				Clock: tt.fields.Clock,
			}
			s.PerformAfter(tt.args.duration, tt.args.action)
		})
	}
}

func TestEventScheduler_PerformRepeatedly(t *testing.T) {
	type fields struct {
		Clock Clock
	}
	type args struct {
		duration time.Duration
		action   timing.Action
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &EventScheduler{
				Clock: tt.fields.Clock,
			}
			s.PerformRepeatedly(tt.args.duration, tt.args.action)
		})
	}
}
