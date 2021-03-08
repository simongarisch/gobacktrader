package asset

import (
	"testing"
)

func TestNewStock(t *testing.T) {
	stock := NewStock("ZZB AU", "AUD")
	if stock.GetTicker() != "ZZB AU" {
		t.Error("Unexpected ticker")
	}
	if stock.GetBaseCurrency() != "AUD" {
		t.Error("Unexpected base currency")
	}

	priceFloat := 2.0
	price := Price{Float64: priceFloat, Valid: true}
	stock.SetPrice(price)

	value := stock.GetValue()
	if !value.Valid {
		t.Error("Expecting a valid stock value")
	}

	valueFloat := value.Float64
	if valueFloat != priceFloat {
		t.Errorf("Unexpected stock value: wanted %.2f, got %.2f", priceFloat, valueFloat)
	}
}
