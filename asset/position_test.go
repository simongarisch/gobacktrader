package asset

import "testing"

func TestNewPosition(t *testing.T) {
	asset, err := NewStock("ZZB AU", "AUD")
	if err != nil {
		t.Errorf("Error in NewStock - %s", err)
	}
	position := NewPosition(&asset, 100.0)
	if position.GetTicker() != "ZZB AU" {
		t.Error("Unexpected ticker.")
	}
	if position.GetBaseCurrency() != "AUD" {
		t.Error("Unexpected base currency")
	}
	if position.GetUnits() != 100.0 {
		t.Error("Unexpected units.")
	}
	if position.GetAsset() != &asset {
		t.Error("Unexpected asset")
	}
}

func TestPositionValue(t *testing.T) {
	stock, err := NewStock("ZZB AU", "AUD")
	if err != nil {
		t.Errorf("Error in NewStock - %s", err)
	}
	position := NewPosition(&stock, 100.0)

	if position.GetUnits() != 100.0 {
		t.Error("Expected 100 units.")
	}

	value := position.GetValue() // price not yet set for stock
	if value != nullPrice {
		t.Error("Expecting an uninitialised value.")
	}

	stock.SetPrice(Price{Float64: 2.50, Valid: true})
	value = position.GetValue()
	if !value.Valid {
		t.Error("Expecting a valid position value.")
	}
	if value.Float64 != 250.0 {
		t.Error("Expecting a value of $250.")
	}

	position.Increment(150.0)
	position.Decrement(25.0)
	if position.GetUnits() != 225.0 {
		t.Error("Expected 225 units.")
	}
	stock.SetPrice(Price{Float64: 2.40, Valid: true})

	// now 225 units at 2.40 each = $540 value in total
	value = position.GetValue()
	if !value.Valid {
		t.Error("Expecting a valid position value.")
	}
	if value.Float64 != 540.0 {
		t.Error("Expecting a value of $540.")
	}
}
