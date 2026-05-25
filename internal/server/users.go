package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/KelvinJRosado/go-server/internal/auth"
	"github.com/KelvinJRosado/go-server/internal/database"
	"github.com/google/uuid"
)

func (ac *apiConfig) createUserHandler() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		// Get input
		type input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		params, ok := getInputStruct[input](res, req)
		if !ok {
			// Util handles writing error response, so we can just return
			return
		}

		// Get hashed password
		hashedPw, err := auth.HashPassword(params.Password)
		if err != nil {
			slog.Error("failed to hash password", "error", err)
			respondWithInternalError(res)
			return
		}

		dbArgs := database.CreateUserParams{
			Email:          params.Email,
			HashedPassword: hashedPw,
		}

		user, err := ac.databaseQueries.CreateUser(req.Context(), dbArgs)
		if err != nil {
			slog.Error("failed to create user", "error", err)
			respondWithInternalError(res)
			return
		}

		respondWithJSON(res, http.StatusCreated, user)
	})
}

func (ac *apiConfig) loginHandler() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		// Get input
		type input struct {
			Email            string `json:"email"`
			Password         string `json:"password"`
			ExpiresInSeconds int    `json:"expires_in_seconds"`
		}

		params, ok := getInputStruct[input](res, req)
		if !ok {
			// Util handles writing error response, so we can just return
			return
		}

		// Validate expiration time, and convert to duration
		// Should be defualted to 1 hour if not set
		// Should be capped at 1 hour
		if params.ExpiresInSeconds <= 0 || params.ExpiresInSeconds > 3600 {
			params.ExpiresInSeconds = 3600
		}
		expiresIn := time.Duration(params.ExpiresInSeconds) * time.Second

		// Pull user data from DB
		user, err := ac.databaseQueries.GetUserByEmail(req.Context(), params.Email)
		if err != nil {
			slog.Error("failed to get user by email", "error", err)
			respondWithInternalError(res)
			return
		}

		// Check password
		pwMatch, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
		if err != nil {
			slog.Error("failed to check password hash", "error", err)
			respondWithInternalError(res)
			return
		}
		if !pwMatch {
			respondWithError(res, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// Generate JWT
		jwt, err := auth.MakeJWT(user.ID, ac.jwtSecret, expiresIn)
		if err != nil {
			slog.Error("failed to generate jwt", "error", err)
			respondWithInternalError(res)
			return
		}

		cleanedUser := struct {
			ID        uuid.UUID `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Email     string    `json:"email"`
			Token     string    `json:"token"`
		}{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Token:     jwt,
		}

		// Password matches, return cleaned user
		respondWithJSON(res, http.StatusOK, cleanedUser)
	})
}
