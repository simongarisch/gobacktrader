package datasources

import (
	"fmt"
	"gobacktrader/asset"
	"testing"
	"time"
)

var (
	testAsset, _  = asset.NewStock("AAPL", "USD")
	testStartDate = time.Date(2021, time.April, 1, 0, 0, 0, 0, time.UTC)
	testEndDate   = time.Date(2021, time.April, 23, 0, 0, 0, 0, time.UTC)
)

func TestFmpCloudQuery(t *testing.T) {
	query := NewFmpCloudQuery(testAsset, testStartDate, testEndDate)

	if query.GetAsset() != testAsset {
		t.Error("Unexpected stock")
	}
	if !query.GetStartDate().Equal(testStartDate) {
		t.Error("Unexpected start date")
	}
	if !query.GetEndDate().Equal(testEndDate) {
		t.Error("Unexpected end date")
	}
	if query.GetAPIKey() != "demo" {
		t.Error("Unexpected api key")
	}

	expectedURLTail := "AAPL?from=2021-04-01&to=2021-04-23&apikey=demo"
	expectedURL := fmpBaseURL + expectedURLTail
	actualURL := query.GetURL()
	if actualURL != expectedURL {
		t.Errorf("Unexpected URL - wanted '%s', got '%s'", expectedURL, actualURL)
	}
}

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
		t.Errorf("Unexpected number of events - wanted 10, got %d", numEvents)
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
