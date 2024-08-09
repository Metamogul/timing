package simulated_time

import (
	"errors"
)

var ErrEventGeneratorFinished = errors.New("event generator is finished")

type EventGenerator interface {
	Pop() *event
	Peek() event

	Finished() bool
}
