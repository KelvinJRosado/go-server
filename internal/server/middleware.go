package server

import (
	"log/slog"
	"net/http"
	"sync/atomic"

	"github.com/KelvinJRosado/go-server/internal/database"
)

type apiConfig struct {
	fileserverHits  atomic.Int32
	databaseQueries *database.Queries
	platform        string
	jwtSecret       string
	polkaKey        string
}

func (ac *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Increment counter
		ac.fileserverHits.Add(1)

		// Call next handler
		next.ServeHTTP(res, req)
	})
}

func (ac *apiConfig) middlewareRequestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		slog.Info("Received request", "method", req.Method, "path", req.URL.Path)
		next.ServeHTTP(res, req)
	})
}
