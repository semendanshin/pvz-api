package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	pvzIDLabel = "pvz_id"
)

var (
	ordersIssuedCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "orders_issued_total",
		Help: "The total number of orders issued",
	}, []string{
		pvzIDLabel,
	})
)

// IncOrdersIssued increments the orders issued counter
func IncOrdersIssued(pvzID string) {
	ordersIssuedCounter.WithLabelValues(pvzID).Inc()
}
