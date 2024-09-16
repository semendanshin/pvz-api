package domain

import (
	"fmt"
	"time"
)

type PackagingType string

const (
	PackagingTypeUnknown PackagingType = "unknown"
	PackagingTypeBox     PackagingType = "box"
	PackagingTypeBag     PackagingType = "bag"
	PackagingTypeFilm    PackagingType = "film"
)

func (p PackagingType) String() string {
	return string(p)
}

func NewPackagingType(p string) (PackagingType, error) {
	switch p {
	case "box":
		return PackagingTypeBox, nil
	case "bag":
		return PackagingTypeBag, nil
	case "film":
		return PackagingTypeFilm, nil
	default:
		return PackagingTypeUnknown, fmt.Errorf(
			"unknown packaging type (available types: box, bag, film): %s", ErrInvalidArgument,
		)
	}
}

// PVZOrder is a struct for PVZ order
type PVZOrder struct {
	OrderID     string
	PVZID       string
	RecipientID string

	Cost   int
	Weight int

	Packaging      PackagingType
	AdditionalFilm bool

	ReceivedAt  time.Time
	StorageTime time.Duration

	IssuedAt   time.Time
	ReturnedAt time.Time
}

func NewPVZOrder(orderID, pvzID, recipientID string, cost, weight int, storageTime time.Duration, packaging PackagingType, additionalFilm bool) PVZOrder {
	return PVZOrder{
		OrderID:        orderID,
		PVZID:          pvzID,
		RecipientID:    recipientID,
		Cost:           cost,
		Weight:         weight,
		Packaging:      packaging,
		AdditionalFilm: additionalFilm,
		ReceivedAt:     time.Now(),
		StorageTime:    storageTime,
	}
}
