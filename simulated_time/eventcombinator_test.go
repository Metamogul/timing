package simulated_time

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"slices"
	"testing"
	"time"
)

func Test_newEventCombinator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		inputs            func() []eventGenerator
		lenInputs         int
		lenFinishedInputs int
		requirePanic      bool
	}{
		{
			name:         "no inputs passed",
			requirePanic: true,
		},
		{
			name: "all inputs finished",
			inputs: func() []eventGenerator {
				mockEventGenerator := NewMockeventGenerator(t)
				mockEventGenerator.EXPECT().
					finished().
					Return(true).
					Once()

				return []eventGenerator{mockEventGenerator}
			},
			lenInputs:         0,
			lenFinishedInputs: 1,
			requirePanic:      true,
		},
		{
			name: "two mixed inputs",
			inputs: func() []eventGenerator {
				mockEventGenerator1 := NewMockeventGenerator(t)
				mockEventGenerator1.EXPECT().
					finished().
					Return(true).
					Once()

				mockEventGenerator2 := NewMockeventGenerator(t)
				mockEventGenerator2.EXPECT().
					finished().
					Return(false).
					Once()

				return []eventGenerator{
					mockEventGenerator1,
					mockEventGenerator2,
				}
			},
			lenInputs:         1,
			lenFinishedInputs: 1,
		},
		{
			name: "two unfinished inputs",
			inputs: func() []eventGenerator {
				mockEventGenerator1 := NewMockeventGenerator(t)
				mockEventGenerator1.EXPECT().
					finished().
					Return(false).
					Once()
				mockEventGenerator1.EXPECT().
					peek().
					Return(event{
						action: noAction{},
						Time:   time.Time{},
					}).
					Maybe()

				mockEventGenerator2 := NewMockeventGenerator(t)
				mockEventGenerator2.EXPECT().
					finished().
					Return(false).
					Once()
				mockEventGenerator2.EXPECT().
					peek().
					Return(event{
						action: noAction{},
						Time:   time.Time{}.Add(time.Second),
					}).
					Maybe()

				return []eventGenerator{
					mockEventGenerator1,
					mockEventGenerator2,
				}
			},
			lenInputs:         2,
			lenFinishedInputs: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.requirePanic {
				require.Panics(t, func() {
					_ = newEventCombinator(tt.inputs()...)
				})
				return
			}

			got := newEventCombinator(tt.inputs()...)

			require.NotNil(t, got.inputs)
			require.NotNil(t, got.finishedInputs)

			require.Len(t, got.inputs, tt.lenInputs)
			require.Len(t, got.finishedInputs, tt.lenFinishedInputs)
		})
	}
}

func Test_eventCombinator_pop(t *testing.T) {
	t.Parallel()

	type fields struct {
		inputs         func() []eventGenerator
		finishedInputs func() []eventGenerator
	}

	tests := []struct {
		name          string
		fields        fields
		finishesInput bool
		want          *event
		requirePanic  bool
	}{
		{
			name: "all inputs finished",
			fields: fields{
				inputs: func() []eventGenerator {
					return make([]eventGenerator, 0)
				},
				finishedInputs: func() []eventGenerator {
					return make([]eventGenerator, 0)
				},
			},
			requirePanic: true,
		},
		{
			name: "success, generator not finished",
			fields: fields{
				inputs: func() []eventGenerator {
					eventGenerator1 := newPeriodicEventGenerator(noAction{}, time.Time{}, nil, time.Minute)
					eventGenerator2 := newPeriodicEventGenerator(noAction{}, time.Time{}, nil, time.Second)
					return []eventGenerator{eventGenerator1, eventGenerator2}
				},
				finishedInputs: func() []eventGenerator {
					return make([]eventGenerator, 0)
				},
			},
			finishesInput: false,
			want: &event{
				action: noAction{},
				Time:   time.Time{}.Add(time.Second),
			},
		},
		{
			name: "success, generator finished",
			fields: fields{
				inputs: func() []eventGenerator {
					eventGenerator1 := newSingleEventGenerator(noAction{}, time.Time{})
					eventGenerator2 := newPeriodicEventGenerator(noAction{}, time.Time{}, nil, time.Second)
					return []eventGenerator{eventGenerator1, eventGenerator2}
				},
				finishedInputs: func() []eventGenerator {
					return make([]eventGenerator, 0)
				},
			},
			finishesInput: true,
			want: &event{
				action: noAction{},
				Time:   time.Time{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			j := &eventCombinator{
				inputs:         tt.fields.inputs(),
				finishedInputs: tt.fields.finishedInputs(),
			}
			j.sortInputs()

			if tt.requirePanic {
				require.Panics(t, func() {
					_ = j.pop()
				})
				return
			}

			if got := j.pop(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pop() = %v, want %v", got, tt.want)
			}

			if !tt.finishesInput {
				require.Len(t, j.inputs, len(tt.fields.inputs()))
				require.Len(t, j.finishedInputs, len(tt.fields.finishedInputs()))
			} else {
				require.Len(t, j.inputs, len(tt.fields.inputs())-1)
				require.Len(t, j.finishedInputs, len(tt.fields.finishedInputs())+1)
			}

		})
	}
}

func Test_eventCombinator_peek(t *testing.T) {
	t.Parallel()

	type fields struct {
		inputs         func() []eventGenerator
		finishedInputs func() []eventGenerator
	}

	tests := []struct {
		name         string
		fields       fields
		want         event
		requirePanic bool
	}{
		{
			name: "all inputs finished",
			fields: fields{
				inputs: func() []eventGenerator {
					return make([]eventGenerator, 0)
				},
				finishedInputs: func() []eventGenerator {
					return make([]eventGenerator, 0)
				},
			},
			requirePanic: true,
		},
		{
			name: "success",
			fields: fields{
				inputs: func() []eventGenerator {
					eventGenerator1 := newPeriodicEventGenerator(noAction{}, time.Time{}, nil, time.Minute)
					eventGenerator2 := newPeriodicEventGenerator(noAction{}, time.Time{}, nil, time.Second)
					return []eventGenerator{eventGenerator1, eventGenerator2}
				},
				finishedInputs: func() []eventGenerator {
					return make([]eventGenerator, 0)
				},
			},
			want: event{
				action: noAction{},
				Time:   time.Time{}.Add(time.Second),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			j := &eventCombinator{
				inputs:         tt.fields.inputs(),
				finishedInputs: tt.fields.finishedInputs(),
			}
			j.sortInputs()

			if tt.requirePanic {
				require.Panics(t, func() {
					_ = j.peek()
				})
				return
			}

			if got := j.peek(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("peek() = %v, want %v", got, tt.want)
			}

			require.Len(t, j.inputs, len(tt.fields.inputs()))
			require.Len(t, j.finishedInputs, len(tt.fields.finishedInputs()))

		})
	}
}

func Test_eventCombinator_finished(t *testing.T) {
	t.Parallel()

	type fields struct {
		inputs         []eventGenerator
		finishedInputs []eventGenerator
	}

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "not finished",
			fields: fields{
				inputs:         []eventGenerator{NewMockeventGenerator(t)},
				finishedInputs: make([]eventGenerator, 0),
			},
			want: false,
		},
		{
			name: "finished",
			fields: fields{
				inputs:         make([]eventGenerator, 0),
				finishedInputs: make([]eventGenerator, 0),
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			j := &eventCombinator{
				inputs:         tt.fields.inputs,
				finishedInputs: tt.fields.finishedInputs,
			}

			if got := j.finished(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("finished() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_eventCombinator_sortInputs(t *testing.T) {
	t.Parallel()

	eventGenerator1 := newPeriodicEventGenerator(noAction{}, time.Time{}, nil, time.Minute)
	eventGenerator2 := newPeriodicEventGenerator(noAction{}, time.Time{}, nil, time.Second)
	eventGenerator3 := newPeriodicEventGenerator(noAction{}, time.Time{}, nil, time.Hour)

	inputs := []eventGenerator{eventGenerator1, eventGenerator2, eventGenerator3}

	j := &eventCombinator{
		inputs:         inputs,
		finishedInputs: make([]eventGenerator, 0),
	}
	j.sortInputs()

	sorted := slices.IsSortedFunc(j.inputs, func(a, b eventGenerator) int {
		return a.peek().Time.Compare(b.peek().Time)
	})
	require.True(t, sorted)
}
