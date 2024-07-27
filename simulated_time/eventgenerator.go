package simulated_time

import (
	"errors"
)

var ErrEventGeneratorFinished = errors.New("event generator is finished")

type eventGenerator interface {
	pop() *event
	peek() event

	finished() bool
}
