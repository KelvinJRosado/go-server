package server

import (
	"encoding/json"
	"io"
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
	errorMessage := struct {
		Error string `json:"error"`
	}{
		Error: msg,
	}

	// Send response
	respondWithJSON(res, code, errorMessage)
}

func respondWithInternalError(res http.ResponseWriter) {

	// Send canned message
	respondWithError(res, http.StatusInternalServerError, "Something went wrong")
}

func getInputStruct[T any](res http.ResponseWriter, req *http.Request) (T, bool) {

	// Init response
	var params T

	// Get body into string
	bodyData, err := io.ReadAll(req.Body)
	if err != nil {
		slog.Error("Error reading request body", "error", err)
		respondWithError(res, http.StatusBadRequest, "Invalid request body")
		return params, false
	}

	// Extract input into struct
	err = json.Unmarshal(bodyData, &params)
	if err != nil {
		// If can't decode response, return 500
		slog.Error("Error decoding parameters", "error", err, "body", string(bodyData))
		respondWithError(res, http.StatusBadRequest, "Invalid request body")
		return params, false
	}

	return params, true
}
