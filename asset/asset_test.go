package asset

import (
	"database/sql"
	"testing"
)

func TestNewAsset(t *testing.T) {
	asset := NewAsset("ZZB AU")
	if asset.GetTicker() != "ZZB AU" {
		t.Error("Unexpected ticker")
	}
	if asset.GetPrice().Valid {
		t.Error("Expecting an uninitialised price")
	}
}

func TestGetSetPrice(t *testing.T) {
	asset := NewAsset("ZZB AU")
	price := sql.NullFloat64{Float64: 2.75, Valid: true}
	asset.SetPrice(price)
	if asset.GetPrice().Float64 != 2.75 {
		t.Error("Unexpected price")
	}
}
