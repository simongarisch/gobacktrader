package asset

import "testing"

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
	if errstr != "'USDA' is an invalid currency code" {
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
