package domain

import (
	"time"
)

// PVZOrder is a struct for PVZ order
type PVZOrder struct {
	OrderID     string
	PVZID       string
	RecipientID string

	ReceivedAt  time.Time
	StorageTime time.Duration

	IssuedAt   time.Time
	ReturnedAt time.Time
}
