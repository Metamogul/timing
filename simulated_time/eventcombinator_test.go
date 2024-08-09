package simulated_time

import (
	"context"
	"github.com/metamogul/timing"
	"github.com/stretchr/testify/require"
	"reflect"
	"slices"
	"testing"
	"time"
)

func Test_newEventCombinator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                  string
		activeGenerators      func() []EventGenerator
		lenActiveGenerators   int
		lenFinishedGenerators int
	}{
		{
			name:                  "no generators passed",
			activeGenerators:      func() []EventGenerator { return nil },
			lenActiveGenerators:   0,
			lenFinishedGenerators: 0,
		},
		{
			name: "all generators finished",
			activeGenerators: func() []EventGenerator {
				mockEventGenerator := NewMockEventGenerator(t)
				mockEventGenerator.EXPECT().
					Finished().
					Return(true).
					Once()

				return []EventGenerator{mockEventGenerator}
			},
			lenActiveGenerators:   0,
			lenFinishedGenerators: 1,
		},
		{
			name: "two mixed generators",
			activeGenerators: func() []EventGenerator {
				mockEventGenerator1 := NewMockEventGenerator(t)
				mockEventGenerator1.EXPECT().
					Finished().
					Return(true).
					Once()

				mockEventGenerator2 := NewMockEventGenerator(t)
				mockEventGenerator2.EXPECT().
					Finished().
					Return(false).
					Once()

				return []EventGenerator{
					mockEventGenerator1,
					mockEventGenerator2,
				}
			},
			lenActiveGenerators:   1,
			lenFinishedGenerators: 1,
		},
		{
			name: "two unfinished generators",
			activeGenerators: func() []EventGenerator {
				mockEventGenerator1 := NewMockEventGenerator(t)
				mockEventGenerator1.EXPECT().
					Finished().
					Return(false).
					Once()
				mockEventGenerator1.EXPECT().
					Peek().
					Return(event{
						Action: timing.NewMockAction(t),
						Time:   time.Time{},
					}).
					Maybe()

				mockEventGenerator2 := NewMockEventGenerator(t)
				mockEventGenerator2.EXPECT().
					Finished().
					Return(false).
					Once()
				mockEventGenerator2.EXPECT().
					Peek().
					Return(event{
						Action: timing.NewMockAction(t),
						Time:   time.Time{}.Add(time.Second),
					}).
					Maybe()

				return []EventGenerator{
					mockEventGenerator1,
					mockEventGenerator2,
				}
			},
			lenActiveGenerators:   2,
			lenFinishedGenerators: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := newEventCombinator(tt.activeGenerators()...)

			require.NotNil(t, got.activeGenerators)
			require.NotNil(t, got.finishedGenerators)

			require.Len(t, got.activeGenerators, tt.lenActiveGenerators)
			require.Len(t, got.finishedGenerators, tt.lenFinishedGenerators)

			sorted := slices.IsSortedFunc(got.activeGenerators, func(a, b EventGenerator) int {
				return a.Peek().Time.Compare(b.Peek().Time)
			})
			require.True(t, sorted)
		})
	}
}

func Test_eventCombinator_add(t *testing.T) {
	t.Parallel()

	type fields struct {
		activeGenerators   []EventGenerator
		finishedGenerators []EventGenerator
	}

	tests := []struct {
		name                string
		fields              fields
		generator           func() EventGenerator
		generatorIsFinished bool
	}{
		{
			name:   "generator finished",
			fields: fields{activeGenerators: []EventGenerator{}, finishedGenerators: []EventGenerator{}},
			generator: func() EventGenerator {
				mockEventGenerator := NewMockEventGenerator(t)
				mockEventGenerator.EXPECT().
					Finished().
					Return(true).
					Once()

				return mockEventGenerator
			},
			generatorIsFinished: true,
		},
		{
			name:   "generator not finished",
			fields: fields{activeGenerators: []EventGenerator{}, finishedGenerators: []EventGenerator{}},
			generator: func() EventGenerator {
				mockEventGenerator := NewMockEventGenerator(t)
				mockEventGenerator.EXPECT().
					Finished().
					Return(false).
					Once()

				return mockEventGenerator
			},
			generatorIsFinished: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := &eventCombinator{
				activeGenerators:   tt.fields.activeGenerators,
				finishedGenerators: tt.fields.finishedGenerators,
			}

			e.add(tt.generator())

			if !tt.generatorIsFinished {
				require.Len(t, e.activeGenerators, len(tt.fields.activeGenerators)+1)
				require.Len(t, e.finishedGenerators, len(tt.fields.finishedGenerators))
			} else {
				require.Len(t, e.activeGenerators, len(tt.fields.activeGenerators))
				require.Len(t, e.finishedGenerators, len(tt.fields.finishedGenerators)+1)
			}

			sorted := slices.IsSortedFunc(e.activeGenerators, func(a, b EventGenerator) int {
				return a.Peek().Time.Compare(b.Peek().Time)
			})
			require.True(t, sorted)
		})
	}
}

func Test_eventCombinator_pop(t *testing.T) {
	t.Parallel()

	type fields struct {
		activeGenerators   func() []EventGenerator
		finishedGenerators func() []EventGenerator
	}

	tests := []struct {
		name              string
		fields            fields
		finishesGenerator bool
		want              *event
		requirePanic      bool
	}{
		{
			name: "all generators finished",
			fields: fields{
				activeGenerators: func() []EventGenerator {
					return make([]EventGenerator, 0)
				},
				finishedGenerators: func() []EventGenerator {
					return make([]EventGenerator, 0)
				},
			},
			requirePanic: true,
		},
		{
			name: "success, generator not finished",
			fields: fields{
				activeGenerators: func() []EventGenerator {
					eventGenerator1 := newPeriodicEventGenerator(timing.NewMockAction(t), time.Time{}, nil, time.Minute, context.Background())
					eventGenerator2 := newPeriodicEventGenerator(timing.NewMockAction(t), time.Time{}, nil, time.Second, context.Background())
					return []EventGenerator{eventGenerator1, eventGenerator2}
				},
				finishedGenerators: func() []EventGenerator {
					return make([]EventGenerator, 0)
				},
			},
			finishesGenerator: false,
			want: &event{
				Action: timing.NewMockAction(t),
				Time:   time.Time{}.Add(time.Second),
			},
		},
		{
			name: "success, generator finished",
			fields: fields{
				activeGenerators: func() []EventGenerator {
					eventGenerator1 := newSingleEventGenerator(timing.NewMockAction(t), time.Time{}, context.Background())
					eventGenerator2 := newPeriodicEventGenerator(timing.NewMockAction(t), time.Time{}, nil, time.Second, context.Background())
					return []EventGenerator{eventGenerator1, eventGenerator2}
				},
				finishedGenerators: func() []EventGenerator {
					return make([]EventGenerator, 0)
				},
			},
			finishesGenerator: true,
			want: &event{
				Action: timing.NewMockAction(t),
				Time:   time.Time{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := &eventCombinator{
				activeGenerators:   tt.fields.activeGenerators(),
				finishedGenerators: tt.fields.finishedGenerators(),
			}
			e.sortActiveGenerators()

			if tt.requirePanic {
				require.Panics(t, func() {
					_ = e.Pop()
				})
				return
			}

			if got := e.Pop(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pop() = %v, want %v", got, tt.want)
			}

			if !tt.finishesGenerator {
				require.Len(t, e.activeGenerators, len(tt.fields.activeGenerators()))
				require.Len(t, e.finishedGenerators, len(tt.fields.finishedGenerators()))
			} else {
				require.Len(t, e.activeGenerators, len(tt.fields.activeGenerators())-1)
				require.Len(t, e.finishedGenerators, len(tt.fields.finishedGenerators())+1)
			}
		})
	}
}

func Test_eventCombinator_peek(t *testing.T) {
	t.Parallel()

	type fields struct {
		activeGenerators   func() []EventGenerator
		finishedGenerators func() []EventGenerator
	}

	tests := []struct {
		name         string
		fields       fields
		want         event
		requirePanic bool
	}{
		{
			name: "all generators finished",
			fields: fields{
				activeGenerators: func() []EventGenerator {
					return make([]EventGenerator, 0)
				},
				finishedGenerators: func() []EventGenerator {
					return make([]EventGenerator, 0)
				},
			},
			requirePanic: true,
		},
		{
			name: "success",
			fields: fields{
				activeGenerators: func() []EventGenerator {
					eventGenerator1 := newPeriodicEventGenerator(timing.NewMockAction(t), time.Time{}, nil, time.Minute, context.Background())
					eventGenerator2 := newPeriodicEventGenerator(timing.NewMockAction(t), time.Time{}, nil, time.Second, context.Background())
					return []EventGenerator{eventGenerator1, eventGenerator2}
				},
				finishedGenerators: func() []EventGenerator {
					return make([]EventGenerator, 0)
				},
			},
			want: event{
				Action: timing.NewMockAction(t),
				Time:   time.Time{}.Add(time.Second),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := &eventCombinator{
				activeGenerators:   tt.fields.activeGenerators(),
				finishedGenerators: tt.fields.finishedGenerators(),
			}
			e.sortActiveGenerators()

			if tt.requirePanic {
				require.Panics(t, func() {
					_ = e.Peek()
				})
				return
			}

			if got := e.Peek(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Peek() = %v, want %v", got, tt.want)
			}

			require.Len(t, e.activeGenerators, len(tt.fields.activeGenerators()))
			require.Len(t, e.finishedGenerators, len(tt.fields.finishedGenerators()))

		})
	}
}

func Test_eventCombinator_finished(t *testing.T) {
	t.Parallel()

	type fields struct {
		activeGenerators   []EventGenerator
		finishedGenerators []EventGenerator
	}

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "not finished",
			fields: fields{
				activeGenerators:   []EventGenerator{NewMockEventGenerator(t)},
				finishedGenerators: make([]EventGenerator, 0),
			},
			want: false,
		},
		{
			name: "finished",
			fields: fields{
				activeGenerators:   make([]EventGenerator, 0),
				finishedGenerators: make([]EventGenerator, 0),
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := &eventCombinator{
				activeGenerators:   tt.fields.activeGenerators,
				finishedGenerators: tt.fields.finishedGenerators,
			}

			if got := e.Finished(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Finished() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_eventCombinator_sortActiveGeneratos(t *testing.T) {
	t.Parallel()

	eventGenerator1 := newPeriodicEventGenerator(timing.NewMockAction(t), time.Time{}, nil, time.Minute, context.Background())
	eventGenerator2 := newPeriodicEventGenerator(timing.NewMockAction(t), time.Time{}, nil, time.Second, context.Background())
	eventGenerator3 := newPeriodicEventGenerator(timing.NewMockAction(t), time.Time{}, nil, time.Hour, context.Background())

	activeGenerators := []EventGenerator{eventGenerator1, eventGenerator2, eventGenerator3}

	e := &eventCombinator{
		activeGenerators:   activeGenerators,
		finishedGenerators: make([]EventGenerator, 0),
	}
	e.sortActiveGenerators()

	sorted := slices.IsSortedFunc(e.activeGenerators, func(a, b EventGenerator) int {
		return a.Peek().Time.Compare(b.Peek().Time)
	})
	require.True(t, sorted)
}
