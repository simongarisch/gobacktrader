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
}
