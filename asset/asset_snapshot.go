package asset

import "time"

// assetSnapshot defines a snapshot in time for a given asset.
type assetSnapshot struct {
	timestamp time.Time
	price Price
	value Price
}

// newAssetSnapshot returns a new instance of assetSnapshot
func newAssetSnapshot(timestamp time.Time, a IAssetReadOnly) assetSnapshot {
	return assetSnapshot {
		timestamp: timestamp,
		price: a.GetPrice(),
		value: a.GetValue(),
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

// GetValue returns the snapshot value.
func (s assetSnapshot) GetValue() Price {
	return s.value
}
