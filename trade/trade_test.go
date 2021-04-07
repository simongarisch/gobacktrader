package trade

import (
	"gobacktrader/asset"
	"gobacktrader/btutil"
	"testing"
)

func TestTradeInit(t *testing.T) {
	portfolio, err1 := asset.NewPortfolio("XXX", "AUD")
	stock, err2 := asset.NewStock("ZZB AU", "AUD")
	if err := btutil.AnyValidError(err1, err2); err != nil {
		t.Errorf("Error in asset init - %s", err)
	}

	trade := Trade{
		portfolio:   portfolio,
		targetAsset: stock,
		units:       100.0,
	}

	if trade.GetPortfolio() != portfolio {
		t.Error("Unexpected portfolio")
	}
	if trade.GetAsset() != stock {
		t.Error("Unexpected asset")
	}
	if trade.GetUnits() != 100 {
		t.Error("Unexpected units")
	}

	// check the trade value
	stock.SetPrice(asset.Price{Float64: 2.50, Valid: true})
	value := trade.GetLocalCurrencyValue()
	if !value.Valid {
		t.Error("Trade should have a valid value")
	}
	if value.Float64 != 250 {
		t.Errorf("Unexpected trade value: wanted 250, got %0.2f", value.Float64)
	}
}
