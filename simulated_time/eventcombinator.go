package simulated_time

import (
	"slices"
)

type eventCombinator struct {
	activeGenerators   []eventGenerator
	finishedGenerators []eventGenerator
}

func newEventCombinator(inputs ...eventGenerator) *eventCombinator {
	combinator := &eventCombinator{
		activeGenerators:   make([]eventGenerator, 0),
		finishedGenerators: make([]eventGenerator, 0),
	}

	for _, input := range inputs {
		if input.finished() {
			combinator.finishedGenerators = append(combinator.finishedGenerators, input)
		} else {
			combinator.activeGenerators = append(combinator.activeGenerators, input)
		}
	}

	combinator.sortActiveGenerators()

	return combinator
}

func (e *eventCombinator) add(generator eventGenerator) {
	if generator.finished() {
		e.finishedGenerators = append(e.finishedGenerators, generator)
		return
	}

	e.activeGenerators = append(e.activeGenerators, generator)

	e.sortActiveGenerators()
}

func (e *eventCombinator) pop() *event {
	if e.finished() {
		panic(ErrEventGeneratorFinished)
	}

	nextEvent := e.activeGenerators[0].pop()

	if e.activeGenerators[0].finished() {
		e.finishedGenerators = append(e.finishedGenerators, e.activeGenerators[0])
		e.activeGenerators = e.activeGenerators[1:]
	}

	e.sortActiveGenerators()

	return nextEvent
}

func (e *eventCombinator) peek() event {
	if e.finished() {
		panic(ErrEventGeneratorFinished)
	}

	return e.activeGenerators[0].peek()
}

func (e *eventCombinator) finished() bool {
	return len(e.activeGenerators) == 0
}

func (e *eventCombinator) sortActiveGenerators() {
	slices.SortStableFunc(e.activeGenerators, func(a, b eventGenerator) int {
		return a.peek().Time.Compare(b.peek().Time)
	})
}
