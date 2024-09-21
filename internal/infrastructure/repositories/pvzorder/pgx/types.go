package pgx

import (
	"homework/internal/domain"
	"time"
)

type pgxPvzOrder struct {
	OrderID     string `db:"order_id"`
	PVZID       string `db:"pvz_id"`
	RecipientID string `db:"recipient_id"`

	Cost   int `db:"cost"`
	Weight int `db:"weight"`

	Packaging      string `db:"packaging"`
	AdditionalFilm bool   `db:"additional_film"`

	ReceivedAt  time.Time     `db:"received_at"`
	StorageTime time.Duration `db:"storage_time"`

	IssuedAt   time.Time `db:"issued_at"`
	ReturnedAt time.Time `db:"returned_at"`

	DeletedAt time.Time `db:"deleted_at"`
}

func newPgxPvzOrder(order domain.PVZOrder) pgxPvzOrder {
	return pgxPvzOrder{
		OrderID:     order.OrderID,
		PVZID:       order.PVZID,
		RecipientID: order.RecipientID,

		Cost:   order.Cost,
		Weight: order.Weight,

		Packaging:      order.Packaging.String(),
		AdditionalFilm: order.AdditionalFilm,

		ReceivedAt:  order.ReceivedAt,
		StorageTime: order.StorageTime,

		IssuedAt:   order.IssuedAt,
		ReturnedAt: order.ReturnedAt,

		DeletedAt: time.Time{},
	}
}

func (p *pgxPvzOrder) ToDomain() domain.PVZOrder {
	return domain.PVZOrder{
		OrderID:     p.OrderID,
		PVZID:       p.PVZID,
		RecipientID: p.RecipientID,

		Cost:   p.Cost,
		Weight: p.Weight,

		Packaging:      domain.PackagingType(p.Packaging),
		AdditionalFilm: p.AdditionalFilm,

		ReceivedAt:  p.ReceivedAt,
		StorageTime: p.StorageTime,

		IssuedAt:   p.IssuedAt,
		ReturnedAt: p.ReturnedAt,
	}
}
