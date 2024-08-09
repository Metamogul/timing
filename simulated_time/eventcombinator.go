package simulated_time

import (
	"slices"
)

type eventCombinator struct {
	activeGenerators   []EventGenerator
	finishedGenerators []EventGenerator
}

func newEventCombinator(inputs ...EventGenerator) *eventCombinator {
	combinator := &eventCombinator{
		activeGenerators:   make([]EventGenerator, 0),
		finishedGenerators: make([]EventGenerator, 0),
	}

	for _, input := range inputs {
		if input.Finished() {
			combinator.finishedGenerators = append(combinator.finishedGenerators, input)
		} else {
			combinator.activeGenerators = append(combinator.activeGenerators, input)
		}
	}

	combinator.sortActiveGenerators()

	return combinator
}

func (e *eventCombinator) add(generator EventGenerator) {
	if generator.Finished() {
		e.finishedGenerators = append(e.finishedGenerators, generator)
		return
	}

	e.activeGenerators = append(e.activeGenerators, generator)

	e.sortActiveGenerators()
}

func (e *eventCombinator) Pop() *event {
	if e.Finished() {
		panic(ErrEventGeneratorFinished)
	}

	nextEvent := e.activeGenerators[0].Pop()

	if e.activeGenerators[0].Finished() {
		e.finishedGenerators = append(e.finishedGenerators, e.activeGenerators[0])
		e.activeGenerators = e.activeGenerators[1:]
	}

	e.sortActiveGenerators()

	return nextEvent
}

func (e *eventCombinator) Peek() event {
	if e.Finished() {
		panic(ErrEventGeneratorFinished)
	}

	return e.activeGenerators[0].Peek()
}

func (e *eventCombinator) Finished() bool {
	return len(e.activeGenerators) == 0
}

func (e *eventCombinator) sortActiveGenerators() {
	slices.SortStableFunc(e.activeGenerators, func(a, b EventGenerator) int {
		return a.Peek().Time.Compare(b.Peek().Time)
	})
}
