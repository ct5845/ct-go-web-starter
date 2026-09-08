package httpserver

import (
	"context"
	"ct-go-web-starter/internal/infrastructure/config"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

const shutdownTimeout = 10 * time.Second

// Kubernetes removes a pod from Service endpoints asynchronously, so traffic
// keeps arriving for a moment after readiness starts failing. Serving through
// that window is what turns a rolling deploy from 502s into a clean handover.
const drainDelay = 5 * time.Second

var draining atomic.Bool

// Draining reports whether shutdown has begun, so the readiness probe can fail
// before the listener closes.
func Draining() bool {
	return draining.Load()
}

// Run serves handler on addr until SIGINT or SIGTERM, then stops reporting ready,
// waits for load balancers to notice, and drains in-flight requests.
func Run(addr string, handler http.Handler) error {
	server := &http.Server{Addr: addr, Handler: handler}

	errs := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errs <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errs:
		return err
	case <-quit:
	}

	draining.Store(true)

	// Only Kubernetes needs the endpoint-propagation window; locally it would
	// just make every Ctrl+C and air restart five seconds slower.
	if config.AppEnv == "prod" {
		slog.Info("Draining before shutdown", "delay", drainDelay)
		time.Sleep(drainDelay)
	}

	slog.Info("Server shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return err
	}

	slog.Info("Server stopped")
	return nil
}
