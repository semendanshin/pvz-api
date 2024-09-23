package pgx

import (
	"github.com/jackc/pgx/v5/pgtype"
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

	ReceivedAt  pgtype.Timestamptz `db:"received_at"`
	StorageTime pgtype.Interval    `db:"storage_time"`

	IssuedAt   pgtype.Timestamptz `db:"issued_at"`
	ReturnedAt pgtype.Timestamptz `db:"returned_at"`

	DeletedAt pgtype.Timestamptz `db:"deleted_at"`
}

func newTimestamptz(t time.Time) pgtype.Timestamptz {
	var valid bool
	if !t.IsZero() {
		valid = true
	}

	return pgtype.Timestamptz{Time: t, Valid: valid}
}

func newInterval(d time.Duration) pgtype.Interval {
	var valid bool
	if d != 0 {
		valid = true
	}

	return pgtype.Interval{Microseconds: d.Microseconds(), Valid: valid}
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

		ReceivedAt:  newTimestamptz(order.ReceivedAt),
		StorageTime: newInterval(order.StorageTime),

		IssuedAt:   newTimestamptz(order.IssuedAt),
		ReturnedAt: newTimestamptz(order.ReturnedAt),

		DeletedAt: newTimestamptz(time.Time{}),
	}
}

func intervalToDuration(i pgtype.Interval) time.Duration {
	const (
		microsecondsPerSecond = 1000000
		microsecondsPerMinute = 60 * microsecondsPerSecond
		microsecondsPerHour   = 60 * microsecondsPerMinute
		microsecondsPerDay    = 24 * microsecondsPerHour
		microsecondsPerMonth  = 30 * microsecondsPerDay
	)

	return time.Duration((int64(i.Months)*microsecondsPerMonth + int64(i.Days)*microsecondsPerDay + i.Microseconds) * 1000)
}

func (p *pgxPvzOrder) ToDomain() domain.PVZOrder {
	println(p.StorageTime.Value())
	return domain.PVZOrder{
		OrderID:     p.OrderID,
		PVZID:       p.PVZID,
		RecipientID: p.RecipientID,

		Cost:   p.Cost,
		Weight: p.Weight,

		Packaging:      domain.PackagingType(p.Packaging),
		AdditionalFilm: p.AdditionalFilm,

		ReceivedAt:  p.ReceivedAt.Time,
		StorageTime: intervalToDuration(p.StorageTime),

		IssuedAt:   p.IssuedAt.Time,
		ReturnedAt: p.ReturnedAt.Time,
	}
}
