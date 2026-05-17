package server

import "net/http"

type healthHandler struct{}

func (healthHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// Update response header
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")

	// Update response body and status
	res.WriteHeader(200)
	res.Write([]byte("OK"))
}
