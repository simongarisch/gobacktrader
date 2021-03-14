package asset

import "time"

// AssetSnapshot defines a snapshot in time for a given asset.
type assetSnapshot struct {
	timestamp time.Time
	price Price
	value Price
}

// NewAssetSnapshot returns a new instance of assetSnapshot
func NewAssetSnapshot(timestamp time.Time, a IAssetReadOnly) assetSnapshot {
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
