package reqlog

import (
	"context"
	"sync"
	"time"
)

type requestContext struct {
	mu        sync.Mutex
	requestID string
	startTime time.Time
	skipped   bool
	spans     []Span
}

// Span records the timing of a named operation within a request.
type Span struct {
	Key         string
	Description string
	Start       time.Duration
	Duration    time.Duration
}

type ctxKey struct{}

func newRequestContext(requestID string) *requestContext {
	return &requestContext{
		requestID: requestID,
		startTime: time.Now(),
	}
}

func fromContext(ctx context.Context) *requestContext {
	rc, _ := ctx.Value(ctxKey{}).(*requestContext)
	return rc
}

// Track records how long the calling scope takes. Use as: defer reqlog.Track(ctx, "key", "description")()
func Track(ctx context.Context, key, description string) func() {
	rc := fromContext(ctx)
	if rc == nil {
		return func() {}
	}
	start := time.Now()
	return func() {
		rc.mu.Lock()
		rc.spans = append(rc.spans, Span{
			Key:         key,
			Description: description,
			Start:       start.Sub(rc.startTime),
			Duration:    time.Since(start),
		})
		rc.mu.Unlock()
	}
}

// Skip marks the request so the middleware omits it from the request log.
func Skip(ctx context.Context) {
	if rc := fromContext(ctx); rc != nil {
		rc.mu.Lock()
		rc.skipped = true
		rc.mu.Unlock()
	}
}

// RequestID returns the request ID stored in ctx, or empty string if none.
func RequestID(ctx context.Context) string {
	if rc := fromContext(ctx); rc != nil {
		return rc.requestID
	}
	return ""
}
