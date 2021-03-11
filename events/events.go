package events

import (
	"errors"
	"sort"
	"time"
)

// IEvent defines the interface for event objects.
type IEvent interface {
	GetTime() time.Time
	IsProcessed() bool
	Process()
}

// Events represents a collection of events.
type Events struct {
	list []IEvent
}

// NewEvents returns a new empty events collection.
func NewEvents() Events {
	return Events{}
}

// Len returns the number of event dates.
func (e Events) Len() int {
	return len(e.list)
}

// IsEmpty returns true if the events list is empty.
func (e Events) IsEmpty() bool {
	if len(e.list) == 0 {
		return true
	}
	return false
}

// Add will add an event to the events collection.
func (e *Events) Add(event IEvent) {
	if e.Contains(event) {
		return // event is already in the list
	}
	e.list = insertSort(e.list, event)
}

// Get will return the most recent event.
func (e *Events) Get() (IEvent, error) {
	if e.IsEmpty() {
		return nil, errors.New("the events list is empty")
	}
	event := e.list[0]
	e.list = e.list[1:]
	return event, nil
}

// FetchAll returns all events from the list.
func (e *Events) FetchAll() []IEvent {
	list := e.list
	e.list = nil
	return list
}

// Contains returns true if some event is in the list, false otherwise
func (e Events) Contains(event IEvent) bool {
	for _, listEvent := range e.list {
		if listEvent == event {
			return true
		}
	}
	return false
}

func insertSort(events []IEvent, event IEvent) []IEvent {
	eventTime := event.GetTime()
	index := sort.Search(len(events), func(i int) bool { return events[i].GetTime().After(eventTime) })
	events = append(events, event)
	copy(events[index+1:], events[index:])
	events[index] = event
	return events
}
