package datasources

import (
	"gobacktrader/asset"
	"testing"
	"time"
)

var (
	testAsset, _  = asset.NewStock("AAPL", "USD")
	testStartDate = time.Date(2021, time.April, 1, 0, 0, 0, 0, time.UTC)
	testEndDate   = time.Date(2021, time.April, 23, 0, 0, 0, 0, time.UTC)
)

func TestQuery(t *testing.T) {
	query := NewQuery(testAsset, testStartDate, testEndDate)

	if query.GetAsset() != testAsset {
		t.Error("Unexpected asset")
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
}

func TestGetSetDetails(t *testing.T) {
	query := NewQuery(testAsset, testStartDate, testEndDate)
	ticker := query.GetTicker()
	apiKey := query.GetAPIKey()
	if ticker != "AAPL" {
		t.Fatalf("Expecting a ticker of 'AAPL', got '%s'", ticker)
	}
	if apiKey != "demo" {
		t.Fatalf("Expecting an API key of 'demo', got '%s'", apiKey)
	}

	query.SetTicker("XXXX").SetAPIKey("YYYY")
	ticker = query.GetTicker()
	apiKey = query.GetAPIKey()
	if ticker != "XXXX" {
		t.Fatalf("Expecting a ticker of 'XXXX', got '%s'", ticker)
	}
	if apiKey != "YYYY" {
		t.Fatalf("Expecting an API key of 'YYYY', got '%s'", apiKey)
	}
}
