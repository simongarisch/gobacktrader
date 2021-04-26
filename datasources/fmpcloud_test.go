package datasources

import (
	"fmt"
	"testing"
	"time"
)

var (
	testStock     = "AAPL"
	testStartDate = time.Date(2021, time.April, 1, 0, 0, 0, 0, time.UTC)
	testEndDate   = time.Date(2021, time.April, 23, 0, 0, 0, 0, time.UTC)
)

func TestFmpCloudQuery(t *testing.T) {
	query := NewFmpCloudQuery(testStock, testStartDate, testEndDate)

	if query.GetStock() != testStock {
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
	query := NewFmpCloudQuery(testStock, testStartDate, testEndDate)
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
