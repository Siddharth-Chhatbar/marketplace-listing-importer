// Package api
package api

import (
	"context"
	"net/http"
	"time"
)

type databaseChecker interface {
	PingContext(context.Context) error
}

func isServerLive(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func isServerReady(w http.ResponseWriter, r *http.Request, db databaseChecker) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func NewHandler(db databaseChecker) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /livez", isServerLive)
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		isServerReady(w, r, db)
	})

	return mux
}
