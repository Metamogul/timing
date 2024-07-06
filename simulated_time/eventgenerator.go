package simulated_time

type eventGenerator interface {
	EventScheduler() EventScheduler
	// After calling SimulatedClock.Forward, nextEvents is
	// supposed to return an eventStream that will yield
	// all  events from the actionTime before SimulatedClock.Forward
	// was called until the new now-actionTime.
	nextEvents() eventStream
}

type singleEventGenerator struct {
	scheduler EventScheduler
}

type repeatedEventGenerator struct {
	scheduler EventScheduler
}
