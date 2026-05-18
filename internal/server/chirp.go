package server

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"unicode/utf8"
)

func (ac *apiConfig) chirpValidateHandler() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Get input
		type input struct {
			Body string `json:"body"`
		}

		// Get body into string
		bodyData, err := io.ReadAll(req.Body)
		if err != nil {
			slog.Error("Error reading request body", "error", err)
			respondWithError(res, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Extract input into struct
		decoder := json.NewDecoder(req.Body)
		params := input{}
		err = decoder.Decode(&params)
		if err != nil {
			// If can't decode response, return 500
			slog.Error("Error decoding parameters", "error", err, "body", string(bodyData))
			respondWithError(res, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Check length of string (runes, not bytes)
		bodyLength := utf8.RuneCountInString(params.Body)

		// Return 400 error if string too long
		if bodyLength > 140 {
			respondWithError(res, http.StatusBadRequest, "Chirp is too long")
			return
		}

		// If string is valid, send success
		type success struct {
			Valid bool `json:"valid"`
		}
		suc := success{
			Valid: true,
		}
		respondWithJSON(res, http.StatusOK, suc)
	})
}
