package server

import (
	"log/slog"
	"net/http"
)

type healthHandler struct{}

func (healthHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	slog.Info("Received health check request")

	// Update response header
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")

	// Update response body and status
	res.WriteHeader(200)
	res.Write([]byte("OK"))
}
