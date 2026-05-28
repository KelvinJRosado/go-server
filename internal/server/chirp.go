package server

import (
	"log/slog"
	"net/http"
	"regexp"
	"unicode/utf8"

	"github.com/KelvinJRosado/go-server/internal/auth"
	"github.com/KelvinJRosado/go-server/internal/database"
	"github.com/google/uuid"
)

func (ac *apiConfig) createChirpHandler() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		// Validate token
		userToken, err := auth.GetBearerToken(req.Header)
		if err != nil {
			slog.Error("failed to get bearer token", "error", err)
			respondWithError(res, http.StatusUnauthorized, "Unauthorized")
			return
		}
		userId, err := auth.ValidateJWT(userToken, ac.jwtSecret)
		if err != nil {
			slog.Error("failed to validate jwt", "error", err)
			respondWithError(res, http.StatusUnauthorized, "Unauthorized")
			return
		}

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

		// Save chirp into DB
		dbArgs := database.CreateChirpParams{
			Body:   body,
			UserID: userId,
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

func (ac *apiConfig) getChirpsHandler() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		// Get chirps from DB
		ch, err := ac.databaseQueries.GetAllChirps(req.Context())
		if err != nil {
			slog.Error("Could not get chirps", "error", err)
			respondWithInternalError(res)
			return
		}

		respondWithJSON(res, http.StatusOK, ch)
	})
}

func (ac *apiConfig) getChirpByIdHandler() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		// Get ID from path
		chirpId, err := uuid.Parse(req.PathValue("id"))
		if err != nil {
			slog.Error("Invalid chirp ID", "error", err)
			respondWithError(res, http.StatusBadRequest, "Invalid chirp ID")
			return
		}

		// Get chirp by ID
		ch, err := ac.databaseQueries.GetChirpById(req.Context(), chirpId)
		if err != nil {
			slog.Error("Could not get chirp", "error", err)
			respondWithError(res, http.StatusNotFound, "Chirp with given ID not found")
			return
		}

		respondWithJSON(res, http.StatusOK, ch)
	})
}

func (ac *apiConfig) deleteChirpByIdHandler() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		// Get ID from path
		chirpId, err := uuid.Parse(req.PathValue("id"))
		if err != nil {
			slog.Error("Invalid chirp ID", "error", err)
			respondWithError(res, http.StatusBadRequest, "Invalid chirp ID")
			return
		}

		// Validate token
		userToken, err := auth.GetBearerToken(req.Header)
		if err != nil {
			slog.Error("failed to get bearer token", "error", err)
			respondWithError(res, http.StatusUnauthorized, "Unauthorized")
			return
		}
		userId, err := auth.ValidateJWT(userToken, ac.jwtSecret)
		if err != nil {
			slog.Error("failed to validate jwt", "error", err)
			respondWithError(res, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// Get chirp by ID
		ch, err := ac.databaseQueries.GetChirpById(req.Context(), chirpId)
		if err != nil {
			slog.Error("Could not get chirp", "error", err)
			respondWithError(res, http.StatusNotFound, "Chirp with given ID not found")
			return
		}

		// Check ownership
		if ch.UserID != userId {
			slog.Error("Can only delete chirp by owner", "user_id", ch.UserID, "token_user_id", userId)
			respondWithError(res, http.StatusForbidden, "Can only delete chirp by owner")
			return
		}

		// Delete chirp by IDs
		dbArgs := database.DeleteChirpByIdParams{
			ID:     chirpId,
			UserID: userId,
		}
		err = ac.databaseQueries.DeleteChirpById(req.Context(), dbArgs)
		if err != nil {
			slog.Error("Could not delete chirp", "error", err)
			respondWithError(res, http.StatusNotFound, "Chirp with given ID not found")
			return
		}

		respondWithJSON(res, http.StatusNoContent, nil)
	})
}
