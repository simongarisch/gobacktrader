package asset

import (
	"testing"
	"time"
)

func TestNewCash(t *testing.T) {
	// start with a vaid currency
	currency := "usd"
	cash, err := NewCash(currency)
	if err != nil {
		t.Errorf("Error in NewCash - %s", err)
	}
	if cash.GetCurrency() != "USD" {
		t.Error("Expected a 'USD' currency code.")
	}
	if cash.GetBaseCurrency() != "USD" {
		t.Error("Expected a 'USD' base currency code")
	}

	// and an invalid currency
	currency = "usda"
	_, err = NewCash(currency)
	errstr := err.Error()
	if errstr != "'USDA' is not a valid currency code" {
		t.Error("Unexpected error string.")
	}
}

func TestPriceValue(t *testing.T) {
	cash, err := NewCash("USD")
	if err != nil {
		t.Errorf("Error in NewCash - %s", err)
	}

	price := cash.GetPrice()
	value := cash.GetValue()

	if !price.Valid {
		t.Error("Price is invalid")
	}
	if !value.Valid {
		t.Error("Value is invalid")
	}

	if price.Float64 != 1.0 {
		t.Error("Price should be 1.0")
	}
	if value.Float64 != 1.0 {
		t.Error("Valud should be 1.0")
	}
}

func TestCashHistory(t *testing.T) {
	cash, err := NewCash("USD")
	if err != nil {
		t.Errorf("Error in NewCash - %s", err)
	}

	time1 := time.Date(2020, time.December, 14, 0, 0, 0, 0, time.UTC)
	time2 := time.Date(2020, time.December, 15, 0, 0, 0, 0, time.UTC)

	cash.TakeSnapshot(time1, &cash)
	cash.TakeSnapshot(time2, &cash)

	history := cash.GetHistory()

	snap1 := history[time1]
	if !snap1.GetTime().Equal(time1) {
		t.Error("snap1 - unexpected time.")
	}
	if snap1.GetPrice().Float64 != 1.0 {
		t.Error("snap1 - unexpected price.")
	}


	snap2 := history[time2]
	if !snap2.GetTime().Equal(time2) {
		t.Error("snap2 - unexpected time.")
	}
	if snap2.GetPrice().Float64 != 1.0 {
		t.Error("snap2 - unexpected price.")
	}
}
