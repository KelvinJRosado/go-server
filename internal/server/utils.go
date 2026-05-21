package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func respondWithJSON(res http.ResponseWriter, code int, payload any) {
	// Marshal response
	data, err := json.Marshal(payload)
	if err != nil {
		slog.Info("Error marshalling JSON", "error", err)

		// Send canned internal error and return
		respondWithInternalError(res)
		return
	}

	// Write response
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(code)
	res.Write(data)
}

func respondWithError(res http.ResponseWriter, code int, msg string) {
	// Build response JSON
	type errorMessage struct {
		Error string `json:"error"`
	}
	em := errorMessage{
		Error: msg,
	}

	// Send response
	respondWithJSON(res, code, em)
}

func respondWithInternalError(res http.ResponseWriter) {

	// Send canned message
	respondWithError(res, http.StatusInternalServerError, "Something went wrong")
}
