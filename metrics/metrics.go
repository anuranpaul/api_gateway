package metrics

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func init() {
	prometheus.MustRegister(RequestsTotal)
	prometheus.MustRegister(RequestDuration)
}

// PrometheusMiddleware records metrics for each request
func PrometheusMiddleware(c *fiber.Ctx) error {
	// Record metrics after the request is processed
	err := c.Next()

	// Record request count
	RequestsTotal.WithLabelValues(
		c.Method(),
		c.Path(),
		fmt.Sprintf("%d", c.Response().StatusCode()),
	).Inc()

	return err
}

// MetricsHandler returns a Fiber handler that serves Prometheus metrics
func MetricsHandler() fiber.Handler {
	return adaptor.HTTPHandler(promhttp.Handler())
}