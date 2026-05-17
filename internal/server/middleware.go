package server

import (
	"log/slog"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (ac *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		// Log request info
		slog.Info("Received fileServer request", "path", req.URL.Path)

		// Increment counter
		ac.fileserverHits.Add(1)

		// Call next handler
		next.ServeHTTP(res, req)
	})
}
