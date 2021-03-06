package asset

import (
	"testing"
)

func TestNewAsset(t *testing.T) {
	asset, err := NewAsset("zzb au", "AUD")
	if err != nil {
		t.Errorf("Error in NewAsset - %s", err)
	}
	if asset.GetTicker() != "ZZB AU" {
		t.Error("Unexpected ticker")
	}
	if asset.GetBaseCurrency() != "AUD" {
		t.Error("Unexpected base currency")
	}
	if asset.GetPrice().Valid {
		t.Error("Expecting an uninitialised price")
	}

	multiplier := asset.GetMultiplier()
	if multiplier != defaultMultiplier {
		t.Errorf("Unexpected multiplier: wanted %0.2f, got %0.2f", defaultMultiplier, multiplier)
	}
}

func TestNewAssetWithMultiplier(t *testing.T) {
	asset, err := NewAssetWithMultiplier("ZZB AU", "AUD", 100.0)
	if err != nil {
		t.Errorf("Error in NewAssetWithMultiplier - %s", err)
	}
	if asset.GetTicker() != "ZZB AU" {
		t.Error("Unexpected ticker")
	}
	if asset.GetMultiplier() != 100.0 {
		t.Errorf("Unexpected multiplier")
	}
	if asset.GetPrice().Valid {
		t.Error("Expecting an uninitialised price")
	}
}

func TestGetSetPrice(t *testing.T) {
	asset, err := NewAsset("ZZB AU", "AUD")
	if err != nil {
		t.Errorf("Error in NewAsset - %s", err)
	}
	priceFloat := 2.75
	price := Price{Float64: priceFloat, Valid: true}
	asset.SetPrice(price)

	actualPrice := asset.GetPrice().Float64
	if actualPrice != priceFloat {
		t.Errorf("Unexpected price: wanted %0.2f, got %0.2f", priceFloat, actualPrice)
	}
}

func TestRevalue(t *testing.T) {
	asset, err := NewAssetWithMultiplier("ZZB AU", "AUD", 100.0)
	if err != nil {
		t.Errorf("Error in NewAssetWithMultiplier - %s", err)
	}
	price := Price{Float64: 2.0, Valid: true}
	asset.SetPrice(price)

	value := asset.GetValue()
	if !value.Valid {
		t.Error("We should have a valid value")
	}
	valueFloat := value.Float64
	expectedValueFloat := 200.0
	if valueFloat != expectedValueFloat {
		t.Errorf("Unexpected value: wanted %0.2f, got %0.2f", expectedValueFloat, valueFloat)
	}
}
