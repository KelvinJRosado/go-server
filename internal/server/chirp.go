package server

import (
	"log/slog"
	"net/http"
	"regexp"
	"unicode/utf8"

	"github.com/KelvinJRosado/go-server/internal/database"
	"github.com/google/uuid"
)

func (ac *apiConfig) createChirpHandler() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Get input
		type input struct {
			Body   string    `json:"body"`
			UserId uuid.UUID `json:"user_id"`
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

		// Save chirp into DB
		dbArgs := database.CreateChirpParams{
			Body:   body,
			UserID: params.UserId,
		}
		ch, err := ac.databaseQueries.CreateChirp(req.Context(), dbArgs)
		if err != nil {
			slog.Error("Could not create chirp", "error", err)
			respondWithError(res, http.StatusBadRequest, "Could not create chirp")
			return
		}

		respondWithJSON(res, http.StatusCreated, ch)
	})
}
