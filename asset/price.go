package asset

import (
	"database/sql"
	"time"
)

// Price is the unit of measurement for price and value.
type Price sql.NullFloat64

var (
	nullPrice = Price{Float64: 0.0, Valid: false}
	nullValue = Price{Float64: 0.0, Valid: false}
	unitPrice = Price{Float64: 1.0, Valid: true}
)

// PriceSnapshot defines a snapshot in time for a given price.
type PriceSnapshot struct {
	timestamp time.Time
	price     Price
}

type iHasGetPrice interface {
	GetPrice() Price
}

// NewPriceSnapshot returns a new instance of PriceSnapshot.
func NewPriceSnapshot(timestamp time.Time, a iHasGetPrice) PriceSnapshot {
	return PriceSnapshot{
		timestamp: timestamp,
		price:     a.GetPrice(),
	}
}

// GetTime returns the timestamp for this snapshot.
func (s PriceSnapshot) GetTime() time.Time {
	return s.timestamp
}

// GetPrice returns the snapshot price.
func (s PriceSnapshot) GetPrice() Price {
	return s.price
}

type priceHistory struct {
	price   Price
	history map[time.Time]PriceSnapshot
}

type iPriceHistory interface {
	GetPrice() Price
	TakeSnapshot(time.Time, iHasGetPrice)
	GetHistory() map[time.Time]PriceSnapshot
}

// GetPrice returns the asset's price.
func (h *priceHistory) GetPrice() Price {
	return h.price
}

// TakeSnapshot records a snapshot at a point in time for future reference.
func (h *priceHistory) TakeSnapshot(timestamp time.Time, asset iHasGetPrice) {
	snap := NewPriceSnapshot(timestamp, asset)
	if h.history == nil {
		h.history = make(map[time.Time]PriceSnapshot)
	}
	h.history[timestamp] = snap
}

// GetHistory returns the record of asset history.
func (h *priceHistory) GetHistory() map[time.Time]PriceSnapshot {
	return h.history
}
