package api

import (
	"ct-go-web-starter/internal/infrastructure/reqlog"
	"ct-go-web-starter/internal/service"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/greet", handleGreet)
}

func handleGreet(w http.ResponseWriter, r *http.Request) {
	defer reqlog.Track(r.Context(), "api.handleGreet", "")()

	greeting, err := service.Greet(r.URL.Query().Get("name"))
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, greeting)
}

type errorResponse struct {
	Error string `json:"error"`
}

// Domain errors are UI states, not failures — they must not become 500s.
func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrEmptyName):
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
	default:
		slog.Error("Unhandled service error", "error", err)
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal server error"})
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
