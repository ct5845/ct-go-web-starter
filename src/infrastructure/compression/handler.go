package compression

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

func Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gz := gzip.NewWriter(w)
		defer gz.Close()

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Set("Vary", "Accept-Encoding")

			w = gzipResponseWriter{ResponseWriter: w, Writer: gz}
		}

		h.ServeHTTP(w, r)
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (g gzipResponseWriter) Write(b []byte) (int, error) {
	return g.Writer.Write(b)
}

func (g gzipResponseWriter) Flush() {
	if flusher, ok := g.Writer.(interface{ Flush() error }); ok {
		flusher.Flush()
	}

	if flusher, ok := g.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}
