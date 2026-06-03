package fileserver

import (
	"crypto/md5"
	"ct-go-web-starter/src/infrastructure/reqlog"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func RegisterRoutes(mux *http.ServeMux, dir string) {
	cachedFS := NewCachedFileServer(dir)
	skipped := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqlog.Skip(r.Context())
		cachedFS.ServeHTTP(w, r)
	})
	mux.Handle("/static/", http.StripPrefix("/static/", skipped))
	mux.Handle("GET /sw.js", http.StripPrefix("/", skipped))
	mux.Handle("GET /robots.txt", http.StripPrefix("/", skipped))
}

type CachedFileServer struct {
	dir   string
	etags map[string]string
	mutex sync.RWMutex
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

	if hasETag {
		w.Header().Set("ETag", etag)
		w.Header().Set("Cache-Control", "public, no-cache, must-revalidate")
	}
	http.FileServer(http.Dir(cfs.dir)).ServeHTTP(w, r)
}

func (cfs *CachedFileServer) RefreshETags() {
	cfs.buildETags()
}
