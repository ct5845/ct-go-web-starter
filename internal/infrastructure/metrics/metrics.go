package metrics

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var requests = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "http_requests_total",
	Help: "Total HTTP requests handled.",
}, []string{"route", "status"})

var requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "http_request_duration_seconds",
	Help:    "HTTP request latency in seconds.",
	Buckets: prometheus.DefBuckets,
}, []string{"route"})

// RecordRequest is keyed on the matched route pattern, never on the request
// path: a path carries IDs, and one series per ID would eventually take the
// scraper down.
func RecordRequest(route string, status int, elapsed time.Duration) {
	if route == "" {
		route = "unmatched"
	}

	requests.WithLabelValues(route, strconv.Itoa(status)).Inc()
	requestDuration.WithLabelValues(route).Observe(elapsed.Seconds())
}

// Serve exposes /metrics on its own port, so scraping never rides on the port
// the application publishes. Losing metrics costs observability rather than
// service, so a bind failure is logged loudly and left to the operator.
func Serve(addr string) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	slog.Info("Metrics server starting", "addr", "http://localhost"+addr+"/metrics")

	go func() {
		if err := http.ListenAndServe(addr, mux); err != nil {
			slog.Error("Metrics server error", "error", err)
		}
	}()
}
