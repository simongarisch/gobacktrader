package datasources

import (
	"fmt"
	"testing"
)

func TestFmpCloudQueryRun(t *testing.T) {
	query := NewFmpCloudQuery(testAsset, testStartDate, testEndDate)
	result, err := query.Run()
	if err != nil {
		t.Fatalf("Error in query.Run() - %s", err)
	}

	for _, item := range result.Historical {
		// fmt.Println(item)
		date, close := item.Date, item.Close
		fmt.Println(fmt.Sprintf("%s: %0.4f", date, close))
	}
}

func TestGenerateEvents(t *testing.T) {
	query := NewFmpCloudQuery(testAsset, testStartDate, testEndDate)
	events, err := query.GenerateEvents()
	if err != nil {
		t.Fatalf("Error in GenerateEvents - %s", err)
	}

	numEvents := len(events)
	if numEvents != 16 {
		t.Errorf("Unexpected number of events - wanted 16, got %d", numEvents)
	}

	found := false
	for _, event := range events {
		if event.GetTime().Equal(testEndDate) {
			found = true
			price := event.GetPrice()
			if !price.Valid {
				t.Error("Expecting a valid price")
			}
			if price.Float64 != 134.32 {
				t.Error("Unexpected price")
			}
		}
	}

	if !found {
		t.Errorf("Unable to find price event for 23 Apr 2021")
	}
}

func TestSetGetTicker(t *testing.T) {
	query := NewFmpCloudQuery(testAsset, testStartDate, testEndDate)
	ticker := query.GetTicker()
	apiKey := query.GetAPIKey()
	if ticker != "AAPL" {
		t.Fatalf("Unexpected ticker - wanted 'AAPL', got '%s'", ticker)
	}
	if apiKey != "demo" {
		t.Fatalf("Unexpected API key - wanted 'demo', got '%s'", apiKey)
	}

	query.SetTicker("XXXX").SetAPIKey("YYYY")
	ticker = query.GetTicker()
	apiKey = query.GetAPIKey()
	if ticker != "XXXX" {
		t.Fatalf("Unexpected ticker - wanted 'XXXX', got '%s'", ticker)
	}
	if apiKey != "YYYY" {
		t.Fatalf("Unexpected API key - wanted 'YYYY', got '%s'", apiKey)
	}
}
