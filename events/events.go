package events

import "time"

// IEvent defines the interface for event objects.
type IEvent interface {
	GetTime() time.Time
	IsProcessed() bool
	Process()
}

// Events represents a collection of events.
type Events struct {
	times    []time.Time
	eventMap map[time.Time][]IEvent
}

// NewEvents returns a new empty events collection.
func NewEvents() Events {
	eventMap := make(map[time.Time][]IEvent)
	return Events{eventMap: eventMap}
}

// Len returns the number of event dates.
func (e Events) Len() int {
	return len(e.eventMap)
}

// Add will add an event to the events collection.
func (e *Events) Add(event IEvent) {
	eventTime := event.GetTime()
	if _, ok := e.eventMap[eventTime]; ok {
		for _, existingEvent := range e.eventMap[eventTime] {
			if existingEvent == event {
				return // don't add the same event twice
			}
		}
		e.eventMap[eventTime] = append(e.eventMap[eventTime], event)
	} else {
		e.eventMap[eventTime] = []IEvent{event}
		e.addTime(eventTime)
	}
}

// GetEventsForTime returns all the events for a given time
// and deletes them from the events collection
func (e *Events) GetEventsForTime(t time.Time) []IEvent {
	eventList, ok := e.eventMap[t]
	if ok {
		delete(e.eventMap, t)
	}
	return eventList
}

func (e *Events) addTime(eventTime time.Time) {
	targetIndex := 0
	for _, currentTime := range e.times {
		if currentTime.After(eventTime) {
			break
		}
		targetIndex++
	}
	//e.times = append(e.times[:targetIndex], eventTime, e.times[targetIndex:]...)
}
