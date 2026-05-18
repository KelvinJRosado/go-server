package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func (ac *apiConfig) chirpValidateHandler() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Get input
		type input struct {
			Body string `json:"body"`
		}

		// Extract input into struct
		decoder := json.NewDecoder(req.Body)
		params := input{}
		err := decoder.Decode(&params)
		if err != nil {
			// If can't decode response, return 500
			slog.Error("Error decoding parameters", "error", err)
			res.WriteHeader(500)
			return
		}

		// TODO: Implement

		// Update response header
		res.Header().Add("Content-Type", "text/plain; charset=utf-8")

		// Update response body and status
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("OK"))
	})
}
