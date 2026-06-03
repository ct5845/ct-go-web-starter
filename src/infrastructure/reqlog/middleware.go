package reqlog

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const waterfallWidth = 40

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := generateRequestID()
			rc := newRequestContext(requestID)

			ctx := context.WithValue(r.Context(), ctxKey{}, rc)
			ctx = withLogAttrs(ctx, slog.String("request_id", requestID))
			r = r.WithContext(ctx)

			w.Header().Set("X-Request-ID", requestID)
			rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(rw, r)

			logRequest(r, requestID, rw.status, rc)
		})
	}
}

func logRequest(r *http.Request, requestID string, status int, rc *requestContext) {
	rc.mu.Lock()
	if rc.skipped {
		rc.mu.Unlock()
		return
	}
	duration := time.Since(rc.startTime)
	spans := rc.spans
	rc.mu.Unlock()

	level := slog.LevelInfo
	if status >= 500 {
		level = slog.LevelError
	} else if status >= 400 || duration > time.Second {
		level = slog.LevelWarn
	}

	attrs := []any{
		"request_id", requestID,
		"method", r.Method,
		"path", r.URL.Path,
		"status", status,
		"duration_ms", duration.Milliseconds(),
		"remote_ip", clientIP(r),
	}

	if len(spans) > 0 {
		attrs = append(attrs, "timing", buildWaterfall(spans, duration))
	}

	slog.Log(r.Context(), level, fmt.Sprintf("%s %s %d", r.Method, r.URL.Path, status), attrs...)
}

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.wroteHeader {
		rw.status = code
		rw.wroteHeader = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.wroteHeader = true
	}
	return rw.ResponseWriter.Write(b)
}

func (rw *responseWriter) Flush() {
	if f, ok := rw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func buildWaterfall(spans []Span, total time.Duration) string {
	sorted := make([]Span, len(spans))
	copy(sorted, spans)
	// sort by start time so the waterfall reads top-to-bottom chronologically
	for i := 1; i < len(sorted); i++ {
		for j := i; j > 0 && sorted[j].Start < sorted[j-1].Start; j-- {
			sorted[j], sorted[j-1] = sorted[j-1], sorted[j]
		}
	}

	totalMs := total.Milliseconds()
	if totalMs == 0 {
		totalMs = 1
	}

	maxKey := 0
	for _, s := range sorted {
		if len(s.Key) > maxKey {
			maxKey = len(s.Key)
		}
	}

	var b strings.Builder
	for _, s := range sorted {
		startMs := s.Start.Milliseconds()
		durMs := s.Duration.Milliseconds()

		startCol := int(float64(startMs) / float64(totalMs) * waterfallWidth)
		barLen := int(float64(durMs) / float64(totalMs) * waterfallWidth)
		if barLen < 1 {
			barLen = 1
		}
		if startCol+barLen > waterfallWidth {
			barLen = waterfallWidth - startCol
		}

		bar := strings.Repeat(" ", startCol) + strings.Repeat("█", barLen) + strings.Repeat(" ", waterfallWidth-startCol-barLen)
		fmt.Fprintf(&b, "%-*s |%s| %dms (start;%dms)\n", maxKey, s.Key, bar, durMs, startMs)
	}
	return b.String()
}

func generateRequestID() string {
	return fmt.Sprintf("req-%d", time.Now().UnixNano())
}

func clientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		parts := strings.Split(fwd, ",")
		return strings.TrimSpace(parts[0])
	}
	if real := r.Header.Get("X-Real-IP"); real != "" {
		return real
	}
	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}
