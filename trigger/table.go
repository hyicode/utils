package trigger

import "sync"

type EventTable struct {
	cbList map[string][]CB
}

func NewEventTable() EventTableI {
	return &EventTable{cbList: make(map[string][]CB)}
}

func (t *EventTable) RegisterCB(key string, cb CB) {
	t.cbList[key] = append(t.cbList[key], cb)
}

func (t *EventTable) CBList(key string) []CB {
	return t.cbList[key]
}

type EventTableMutex struct {
	mu sync.RWMutex
	EventTable
}

func NewEventTableMutex() EventTableI {
	return &EventTableMutex{EventTable: EventTable{cbList: make(map[string][]CB)}}
}

func (t *EventTableMutex) RegisterCB(key string, cb CB) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.EventTable.RegisterCB(key, cb)
}

func (t *EventTableMutex) CBList(key string) []CB {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.EventTable.CBList(key)
}
