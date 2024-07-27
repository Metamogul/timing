package simulated_time

import (
	"slices"
)

type joinedEventGenerator struct {
	inputs         []eventGenerator
	finishedInputs []eventGenerator
}

func newJoinedEventGenerator(inputs ...eventGenerator) *joinedEventGenerator {
	joinedGenerator := &joinedEventGenerator{
		inputs: inputs,
	}

	for _, input := range inputs {
		if input.finished() {
			joinedGenerator.finishedInputs = append(joinedGenerator.finishedInputs, input)
		} else {
			joinedGenerator.inputs = append(joinedGenerator.inputs, input)
		}
	}

	if len(joinedGenerator.inputs) == 0 {
		panic("must provide at least one inputs input that is not finished")
	}

	joinedGenerator.sortInputs()

	return joinedGenerator
}

func (j *joinedEventGenerator) pop() *event {
	if j.finished() {
		panic(ErrEventGeneratorFinished)
	}

	nextEvent := j.inputs[0].pop()

	if j.inputs[0].finished() {
		j.finishedInputs = append(j.finishedInputs, j.inputs[0])
		j.inputs = j.inputs[1:]
	}

	j.sortInputs()

	return nextEvent
}

func (j *joinedEventGenerator) peek() event {
	if j.finished() {
		panic(ErrEventGeneratorFinished)
	}

	return j.inputs[0].peek()
}

func (j *joinedEventGenerator) finished() bool {
	return len(j.inputs) > 0
}

func (j *joinedEventGenerator) sortInputs() {
	slices.SortStableFunc(j.inputs, func(a, b eventGenerator) int {
		return a.peek().Time.Compare(b.peek().Time)
	})
}
