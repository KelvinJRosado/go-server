package server

import (
	"log/slog"
	"net/http"
)

func (ac *apiConfig) healthHandler() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		slog.Info("Received health check request")

		// Update response header
		res.Header().Add("Content-Type", "text/plain; charset=utf-8")

		// Update response body and status
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("OK"))
	})
}
