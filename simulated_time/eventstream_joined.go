package simulated_time

import (
	"slices"
	"time"
)

type joinedEventsStream struct {
	inputStreams  []eventStream
	closedStreams []eventStream
}

func newJoinedEventsStream(input ...eventStream) *joinedEventsStream {
	joinedStream := &joinedEventsStream{
		inputStreams: input,
	}

	for _, stream := range input {
		if stream.closed() {
			joinedStream.closedStreams = append(joinedStream.closedStreams, stream)
		} else {
			joinedStream.inputStreams = append(joinedStream.inputStreams, stream)
		}
	}

	if len(joinedStream.inputStreams) == 0 {
		panic("must provide at least one input stream that is not closed")
	}

	joinedStream.sortInput()

	return joinedStream
}

func (j *joinedEventsStream) popNextEvent() *event {
	if j.closed() {
		panic(ErrEventStreamClosed)
	}

	nextEvent := j.inputStreams[0].popNextEvent()

	if j.inputStreams[0].closed() {
		j.closedStreams = append(j.closedStreams, j.inputStreams[0])
		j.inputStreams = j.inputStreams[1:]
	}

	j.sortInput()

	return nextEvent
}

func (j *joinedEventsStream) peakNextTime() time.Time {
	if j.closed() {
		panic(ErrEventStreamClosed)
	}

	return j.inputStreams[0].peakNextTime()
}

func (j *joinedEventsStream) closed() bool {
	return len(j.inputStreams) > 0
}

func (j *joinedEventsStream) sortInput() {
	slices.SortStableFunc(j.inputStreams, func(a, b eventStream) int {
		return a.peakNextTime().Compare(b.peakNextTime())
	})
}
