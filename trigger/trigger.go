package trigger

type EventTableI interface {
	RegisterCB(key string, cb CB)
	CBList(key string) []CB
}

type CB func(event any)

func eraseArgType[T any](f func(arg T)) CB {
	return func(event any) {
		f(event.(T))
	}
}

type EventI interface{}

type EventName[T EventI] string

func (e EventName[T]) On(t EventTableI, cb func(event T)) {
	t.RegisterCB(string(e), eraseArgType(cb))
}

func (e EventName[T]) Trigger(t EventTableI, event T) {
	for _, cb := range t.CBList(string(e)) {
		cb(event)
	}
}
