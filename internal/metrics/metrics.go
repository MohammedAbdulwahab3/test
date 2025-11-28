package metrics

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP metrics
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// Database metrics
	dbQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	// Business metrics
	itemsTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "items_total",
			Help: "Total number of items by status",
		},
		[]string{"status"},
	)

	itemOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "item_operations_total",
			Help: "Total number of item operations",
		},
		[]string{"operation"},
	)
)

// MetricsMiddleware tracks HTTP request metrics
func MetricsMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Record metrics
		duration := time.Since(start).Seconds()
		status := c.Response().StatusCode()
		method := c.Method()
		path := c.Route().Path

		// If path is empty (not a registered route), use the actual path
		if path == "" {
			path = c.Path()
		}

		httpRequestsTotal.WithLabelValues(method, path, string(rune(status/100))+"xx").Inc()
		httpRequestDuration.WithLabelValues(method, path).Observe(duration)

		return err
	}
}

// RecordDBQuery records database query duration
func RecordDBQuery(operation string, duration time.Duration) {
	dbQueryDuration.WithLabelValues(operation).Observe(duration.Seconds())
}

// RecordItemOperation records item CRUD operations
func RecordItemOperation(operation string) {
	itemOperationsTotal.WithLabelValues(operation).Inc()
}

// UpdateItemsTotal updates the total items gauge
func UpdateItemsTotal(status string, count float64) {
	itemsTotal.WithLabelValues(status).Set(count)
}
