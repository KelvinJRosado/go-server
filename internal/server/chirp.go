package server

import (
	"net/http"
	"regexp"
	"unicode/utf8"
)

func (ac *apiConfig) createChirpHandler() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Get input
		type input struct {
			Body string `json:"body"`
		}

		// Parse input
		params, ok := getInputStruct[input](res, req)
		if !ok {
			// Util handles writing error response, so we can just return
			return
		}

		// Check length of string (runes, not bytes)
		body := params.Body
		bodyLength := utf8.RuneCountInString(body)

		// Return 400 error if string too long
		if bodyLength > 140 {
			respondWithError(res, http.StatusBadRequest, "Chirp is too long")
			return
		}

		// Filter
		blockList := []string{"kerfuffle", "sharbert", "fornax"}

		for _, word := range blockList {
			// Build case-insensitve regex search, allowing for special characters
			re := regexp.MustCompile(`(?i)(^|\s)(` + word + `)(\s|$)`)

			// Filter as needed
			body = re.ReplaceAllString(body, `${1}****${3}`)
		}

		// If string is valid, send success
		type success struct {
			CleanedBody string `json:"cleaned_body"`
		}
		suc := success{
			CleanedBody: body,
		}
		respondWithJSON(res, http.StatusOK, suc)
	})
}
