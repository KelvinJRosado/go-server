package server

import (
	"fmt"
	"net/http"
)

func (ac *apiConfig) metricsHandler() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Update response header
		res.Header().Add("Content-Type", "text/plain; charset=utf-8")

		// Update response body and status
		res.WriteHeader(http.StatusOK)
		res.Write(fmt.Appendf(nil, "Hits: %d", ac.fileserverHits.Load()))
	})
}

func (ac *apiConfig) resetMetricsHandler() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Reset counter back to 0
		ac.fileserverHits.Store(0)

		// Update response header
		res.Header().Add("Content-Type", "text/plain; charset=utf-8")

		// Update response body and status
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("Counter was reset"))
	})
}
