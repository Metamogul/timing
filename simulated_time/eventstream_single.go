package simulated_time

import "time"

type singleEventStream struct {
	scheduler eventScheduler
	*event
}

func newSingleEventStream(action action, actionTime time.Time, scheduler eventScheduler) *singleEventStream {
	return &singleEventStream{
		scheduler: scheduler,
		event:     newEvent(action, actionTime, scheduler),
	}
}

func (s *singleEventStream) popNextEvent() *event {
	if s.closed() {
		panic(ErrEventStreamClosed)
	}

	if !s.event.actionTime.Before(s.scheduler.Now()) {
		return nil
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
