package health

import (
	"ct-go-web-starter/internal/infrastructure/httpserver"
	"ct-go-web-starter/internal/infrastructure/reqlog"
	"ct-go-web-starter/internal/service"
	"encoding/json"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /healthz", handleLive)
	mux.HandleFunc("GET /ready", handleReady)
}

// Liveness answers "is this process wedged?" and deliberately checks no
// dependencies: a failing database would otherwise restart every pod instead of
// just taking them out of service.
func handleLive(w http.ResponseWriter, r *http.Request) {
	reqlog.Skip(r.Context())
	writeJSON(w, http.StatusOK, service.CheckHealth())
}

// Readiness answers "should this pod receive traffic?" and fails as soon as
// shutdown begins, so Kubernetes stops routing before the listener closes.
func handleReady(w http.ResponseWriter, r *http.Request) {
	reqlog.Skip(r.Context())

	status := service.CheckHealth()
	if httpserver.Draining() {
		status.Status = "draining"
		writeJSON(w, http.StatusServiceUnavailable, status)
		return
	}

	writeJSON(w, http.StatusOK, status)
}

func writeJSON(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(body)
}
