package simulated_time

import (
	"slices"
)

type eventCombinator struct {
	inputs         []eventGenerator
	finishedInputs []eventGenerator
}

func newEventCombinator(inputs ...eventGenerator) *eventCombinator {
	combinator := &eventCombinator{
		inputs:         make([]eventGenerator, 0),
		finishedInputs: make([]eventGenerator, 0),
	}

	for _, input := range inputs {
		if input.finished() {
			combinator.finishedInputs = append(combinator.finishedInputs, input)
		} else {
			combinator.inputs = append(combinator.inputs, input)
		}
	}

	if len(combinator.inputs) == 0 {
		panic("must provide at least one inputs input that is not finished")
	}

	combinator.sortInputs()

	return combinator
}

func (e *eventCombinator) pop() *event {
	if e.finished() {
		panic(ErrEventGeneratorFinished)
	}

	nextEvent := e.inputs[0].pop()

	if e.inputs[0].finished() {
		e.finishedInputs = append(e.finishedInputs, e.inputs[0])
		e.inputs = e.inputs[1:]
	}

	e.sortInputs()

	return nextEvent
}

func (e *eventCombinator) peek() event {
	if e.finished() {
		panic(ErrEventGeneratorFinished)
	}

	return e.inputs[0].peek()
}

func (e *eventCombinator) finished() bool {
	return len(e.inputs) == 0
}

func (e *eventCombinator) sortInputs() {
	slices.SortStableFunc(e.inputs, func(a, b eventGenerator) int {
		return a.peek().Time.Compare(b.peek().Time)
	})
}
