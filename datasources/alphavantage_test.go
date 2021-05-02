package datasources

import (
	"testing"
)

func TestAlphaVantageQueryGetURL(t *testing.T) {
	query := NewAlphaVantageQuery(testAsset, testStartDate, testEndDate)

	url := query.GetURL()
	expectedURL := "https://www.alphavantage.co/query?function=TIME_SERIES_DAILY_ADJUSTED&symbol=AAPL&outputsize=full&apikey=demo"
	if url != expectedURL {
		t.Errorf("Expecting a URL of '%s', got '%s'", expectedURL, url)
	}
}

func TestAlphaVantageRun(t *testing.T) {
	query := NewAlphaVantageQuery(testAsset, testStartDate, testEndDate)
	_, err := query.Run()
	if err != nil {
		t.Errorf("Error in query.Run() - %s", err)
	}
}

func TestAlphaVantageGenerateEvents(t *testing.T) {
	query := NewAlphaVantageQuery(testAsset, testStartDate, testEndDate)
	// query.SetAPIKey("----")
	if query.GetAPIKey() == "demo" {
		t.Skip() // data won't always get passed back for the demo account
	}
	events, err := query.GenerateEvents()
	if err != nil {
		t.Fatalf("Error in GenerateEvents - %s", err)
	}

	numEvents := len(events)
	if numEvents != 16 {
		t.Errorf("Unexpected number of events - wanted 16, got %d", numEvents)
	}
}
