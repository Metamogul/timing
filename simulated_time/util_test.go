package simulated_time

func ptr[T any](t T) *T {
	return &t
}

type noAction struct{}

func (n noAction) perform() {}
