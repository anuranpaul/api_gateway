package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
)

var (
	// Counter for total requests processed by API Gateway
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_gateway_requests_total",
			Help: "Total number of requests handled by the API Gateway",
		},
		[]string{"method", "status"},
	)

	// Histogram to record request latency
	latencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_gateway_latency_seconds",
			Help:    "Histogram of request latency for the API Gateway",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)
)

// Initialize metrics
func InitMetrics() {
	// Register the metrics with Prometheus
	prometheus.MustRegister(requestsTotal)
	prometheus.MustRegister(latencyHistogram)
}

// Middleware to track request metrics
func RequestMetricsMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		// Process the request
		err := c.Next()

		// Record the metrics for Prometheus
		duration := time.Since(start)
		status := c.Response().StatusCode()

		// Increment total requests count
		requestsTotal.WithLabelValues(c.Method(), string(status)).Inc()

		// Record latency for requests
		latencyHistogram.WithLabelValues(c.Method(), string(status)).Observe(duration.Seconds())

		// Optionally log the request details (duration)
		log.Printf("Request processed in %s", duration)
		return err
	}
}

// Expose Prometheus metrics at /metrics endpoint
func ExposeMetrics(app *fiber.App) {
	app.Get("/metrics", promhttp.Handler())
}
