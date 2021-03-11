package events

import (
	"testing"
	"time"
)

type testEvent struct {
	eventTime   time.Time
	isProcessed bool
}

func newTestEvent(eventTime time.Time) testEvent {
	return testEvent{eventTime: eventTime, isProcessed: false}
}

func (t testEvent) GetTime() time.Time {
	return t.eventTime
}

func (t testEvent) IsProcessed() bool {
	return t.isProcessed
}

func (t *testEvent) Process() {
	t.isProcessed = true
}

func TestEventsLen(t *testing.T) {
	events := NewEvents()
	if events.Len() != 0 {
		t.Error("Expecting a new events collection to be empty.")
	}

	eventTime1 := time.Date(2020, time.December, 14, 0, 0, 0, 0, time.UTC)
	eventTime2 := time.Date(2020, time.December, 15, 0, 0, 0, 0, time.UTC)
	event1 := newTestEvent(eventTime1)
	event2 := newTestEvent(eventTime2)

	events.Add(&event1)
	if events.Len() != 1 {
		t.Errorf("Expecting an events.Len of 1, got %d", events.Len())
	}
	events.Add(&event1) // adding the same event again should do nothing
	if events.Len() != 1 {
		t.Errorf("Expecting an events.Len of 1, got %d", events.Len())
	}

	events.Add(&event2)
	if events.Len() != 2 {
		t.Errorf("Expecting an events.Len of 2, got %d", events.Len())
	}

	events.GetEventsForTime(eventTime1)
	if events.Len() != 1 {
		t.Errorf("Expecting an events.Len of 1, got %d", events.Len())
	}
}
