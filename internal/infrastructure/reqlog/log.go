package reqlog

import (
	"context"
	"log/slog"
)

type logAttrsKey struct{}

func withLogAttrs(ctx context.Context, attrs ...slog.Attr) context.Context {
	existing, _ := ctx.Value(logAttrsKey{}).([]slog.Attr)
	return context.WithValue(ctx, logAttrsKey{}, append(existing, attrs...))
}

// ContextAttrs returns the slog attributes stored in ctx by the middleware.
// Used by the slog handler to inject them into every log record.
func ContextAttrs(ctx context.Context) []slog.Attr {
	attrs, _ := ctx.Value(logAttrsKey{}).([]slog.Attr)
	return attrs
}

// ContextHandler wraps an slog.Handler and injects per-request attributes
// stored in context (e.g. request_id) into every log record.
type ContextHandler struct {
	Inner slog.Handler
}

func (h *ContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Inner.Enabled(ctx, level)
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if attrs := ContextAttrs(ctx); len(attrs) > 0 {
		r.AddAttrs(attrs...)
	}
	return h.Inner.Handle(ctx, r)
}

func (h *ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ContextHandler{Inner: h.Inner.WithAttrs(attrs)}
}

func (h *ContextHandler) WithGroup(name string) slog.Handler {
	return &ContextHandler{Inner: h.Inner.WithGroup(name)}
}
