package asset

import (
	"database/sql"
	"testing"
)

func TestNewStock(t *testing.T) {
	stock := NewStock("ZZB AU")
	if stock.GetTicker() != "ZZB AU" {
		t.Error("Unexpected ticker")
	}

	priceFloat := 2.0
	price := sql.NullFloat64{Float64: priceFloat, Valid: true}
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
