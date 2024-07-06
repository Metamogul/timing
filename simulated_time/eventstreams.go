package simulated_time

import (
	"errors"
	"github.com/metamogul/timing"
	"slices"
	"sync"
	"time"
)

var ErrEventStreamClosed = errors.New("event stream is closed")

type eventStream interface {
	popNextEvent() *event
	peakNextTime() time.Time
	closed() bool
}

type eventStreamEventScheduler interface {
	timing.Clock
	eventCompletionWaitGroup() sync.WaitGroup
}

/////////////////////////
/// singleEventStream ///
/////////////////////////

type singleEventStream struct {
	*event
}

func newSingleEventStream(action func(), actionTime time.Time, scheduler eventScheduler) singleEventStream {
	return singleEventStream{
		event: newEvent(action, actionTime, scheduler),
	}
}

func (s *singleEventStream) popNextEvent() *event {
	if s.closed() {
		panic(ErrEventStreamClosed)
	}

	event := s.event
	s.event = nil

	return event
}

func (s *singleEventStream) peakNextTime() time.Time {
	if s.closed() {
		panic(ErrEventStreamClosed)
	}

	return s.event.actionTime
}

func (s *singleEventStream) closed() bool {
	return s.event == nil
}

/////////////////////////
/// multiEventsStream ///
/////////////////////////

type multiEventsStream struct {
	action    func()
	from      time.Time
	to        time.Time
	interval  time.Duration
	scheduler eventScheduler

	currentEvent *event
}

func newMultiEventsStream(action func(), from, to time.Time, interval time.Duration, scheduler eventScheduler) multiEventsStream {
	firstEvent := newEvent(action, from, scheduler)

	return multiEventsStream{
		action:    action,
		from:      from,
		to:        to,
		interval:  interval,
		scheduler: scheduler,

		currentEvent: firstEvent,
	}
}

func (m *multiEventsStream) popNextEvent() *event {
	if m.closed() {
		panic(ErrEventStreamClosed)
	}

	currentEvent := m.currentEvent
	m.currentEvent = newEvent(m.action, m.currentEvent.actionTime.Add(m.interval), m.scheduler)

	return currentEvent
}

func (m *multiEventsStream) peakNextTime() time.Time {
	if m.closed() {
		panic(ErrEventStreamClosed)
	}

	return m.currentEvent.actionTime
}

func (m *multiEventsStream) closed() bool {
	return !m.currentEvent.actionTime.Before(m.scheduler.Now())
}

//////////////////////////
/// joinedEventsStream ///
//////////////////////////

type joinedEventsStream struct {
	inputStreams  []eventStream
	closedStreams []eventStream
}

func newJoinedEventsStream(input ...eventStream) *joinedEventsStream {
	joinedEventsStream := &joinedEventsStream{
		inputStreams: input,
	}

	for _, stream := range input {
		if stream.closed() {
			joinedEventsStream.closedStreams = append(joinedEventsStream.closedStreams, stream)
			continue
		}

		joinedEventsStream.inputStreams = append(joinedEventsStream.inputStreams, stream)
	}

	joinedEventsStream.sortInput()

	return joinedEventsStream
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
