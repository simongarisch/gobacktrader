package asset

import "time"

type iHasGetPrice interface {
	GetPrice() Price
}

// assetSnapshot defines a snapshot in time for a given asset.
type assetSnapshot struct {
	timestamp time.Time
	price Price
}

// newAssetSnapshot returns a new instance of assetSnapshot
func newAssetSnapshot(timestamp time.Time, a iHasGetPrice) assetSnapshot {
	return assetSnapshot {
		timestamp: timestamp,
		price: a.GetPrice(),
	}
}

// GetTimestamp returns the timestamp for this snapshot.
func (s assetSnapshot) GetTime() time.Time {
	return s.timestamp
}

// GetPrice returns the snapshot price.
func (s assetSnapshot) GetPrice() Price {
	return s.price
}

type hasPriceHistory struct {
	history map[time.Time]assetSnapshot
}

type iHasPriceHistory interface {
	TakeSnapshot(time.Time, iHasGetPrice)
	GetHistory() map[time.Time]assetSnapshot
}

// TakeSnapshot records a snapshot at a point in time for future reference.
func (h *hasPriceHistory) TakeSnapshot(timestamp time.Time, asset iHasGetPrice) {
	snap := newAssetSnapshot(timestamp, asset)
	if h.history == nil {
		h.history = make(map[time.Time]assetSnapshot)
	}
	h.history[timestamp] = snap
}

// GetHistory returns the record of asset history.
func (h *hasPriceHistory) GetHistory() map[time.Time]assetSnapshot {
	return h.history
}
