package events

import (
	"gobacktrader/btutil"
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
	if !events.IsEmpty() {
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

	events.Get()
	if events.Len() != 1 {
		t.Errorf("Expecting an events.Len of 1, got %d", events.Len())
	}
}

func TestEventsAddGet(t *testing.T) {
	events := NewEvents()
	if !events.IsEmpty() {
		t.Error("Expecting a new events collection to be empty.")
	}

	// we cannot get an event if events is empty
	event, err := events.Get()
	errStr := btutil.GetErrorString(err)
	if errStr != "the events list is empty" {
		t.Errorf("Unexpected error string, got '%s'", errStr)
	}
	if event != nil {
		t.Error("Expecting no event to be passed back.")
	}

	time1 := time.Date(2020, time.December, 14, 0, 0, 0, 0, time.UTC)
	time2 := time.Date(2020, time.December, 15, 0, 0, 0, 0, time.UTC)
	time3 := time.Date(2020, time.December, 16, 0, 0, 0, 0, time.UTC)

	event1 := newTestEvent(time1)
	event2 := newTestEvent(time2)
	event3 := newTestEvent(time3)

	// add the events
	events.Add(&event2)
	if events.Len() != 1 {
		t.Errorf("Expecting a length of one.")
	}
	events.Add(&event3)
	if events.Len() != 2 {
		t.Errorf("Expecting a length of two.")
	}
	events.Add(&event1)
	if events.Len() != 3 {
		t.Errorf("Expecting a length of three.")
	}

	// and get them back
	firstEvent, err := events.Get()
	if err != nil {
		t.Errorf("Error in events.Get - %s", err)
	}
	if !firstEvent.GetTime().Equal(time1) {
		t.Error("Expecting event1 to be returned")
	}

	eventsList := events.FetchAll()
	if len(eventsList) != 2 {
		t.Error("Expecting two additional events to be returned.")
	}
	if !eventsList[0].GetTime().Equal(time2) {
		t.Error("Expecting event2 to be returned")
	}
	if !eventsList[1].GetTime().Equal(time3) {
		t.Error("Expecting event3 to be returned")
	}
}
