package fileserver

import (
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type CachedFileServer struct {
	dir   string
	etags map[string]string
	mutex sync.RWMutex
}

type gzipResponseWriter struct {
	http.ResponseWriter
	gzWriter  *gzip.Writer
	etag      string
	hasETag   bool
	headerSet bool
}

func (rw *gzipResponseWriter) Write(b []byte) (int, error) {
	if !rw.headerSet {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.gzWriter.Write(b)
}

func (rw *gzipResponseWriter) WriteHeader(statusCode int) {
	if !rw.headerSet {
		if rw.hasETag {
			rw.Header().Set("ETag", rw.etag)
		}
		rw.Header().Set("Cache-Control", "public, no-cache, must-revalidate")
		rw.Header().Set("Content-Encoding", "gzip")
		rw.Header().Del("Content-Length")
		rw.headerSet = true
	}
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *gzipResponseWriter) Close() error {
	return rw.gzWriter.Close()
}

func NewCachedFileServer(dir string) *CachedFileServer {
	cfs := &CachedFileServer{
		dir:   dir,
		etags: make(map[string]string),
	}
	cfs.buildETags()
	return cfs
}

func (cfs *CachedFileServer) buildETags() {
	cfs.mutex.Lock()
	defer cfs.mutex.Unlock()

	err := filepath.Walk(cfs.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				slog.Error("Error opening file", "path", path, "error", err)
				return nil
			}
			defer file.Close()

			hash := md5.New()
			if _, err := io.Copy(hash, file); err != nil {
				slog.Error("Error hashing file", "path", path, "error", err)
				return nil
			}

			// Calculate relative path from dir and convert to URL path
			relPath, err := filepath.Rel(cfs.dir, path)
			if err != nil {
				return err
			}

			// Convert Windows paths to URL paths
			urlPath := strings.ReplaceAll(relPath, "\\", "/")
			etag := fmt.Sprintf(`"%x"`, hash.Sum(nil))
			cfs.etags[urlPath] = etag
		}
		return nil
	})

	if err != nil {
		slog.Error("Error building ETags", "error", err)
	}
}

func (cfs *CachedFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cfs.mutex.RLock()
	etag, hasETag := cfs.etags[r.URL.Path]
	cfs.mutex.RUnlock()

	if hasETag {
		// Handle If-None-Match header
		if match := r.Header.Get("If-None-Match"); match == etag {
			w.Header().Set("ETag", etag)
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	// Check if client accepts gzip
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		gzWriter := gzip.NewWriter(w)
		defer gzWriter.Close()

		gzWrapped := &gzipResponseWriter{
			ResponseWriter: w,
			gzWriter:       gzWriter,
			etag:           etag,
			hasETag:        hasETag,
		}

		fileServer := http.FileServer(http.Dir(cfs.dir))
		fileServer.ServeHTTP(gzWrapped, r)
		return
	}

	// Fallback: serve without gzip
	if hasETag {
		w.Header().Set("ETag", etag)
		w.Header().Set("Cache-Control", "public, no-cache, must-revalidate")
	}
	fileServer := http.FileServer(http.Dir(cfs.dir))
	fileServer.ServeHTTP(w, r)
}

func (cfs *CachedFileServer) RefreshETags() {
	cfs.buildETags()
}
