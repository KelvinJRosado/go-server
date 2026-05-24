package server

import (
	"log/slog"
	"net/http"

	"github.com/KelvinJRosado/go-server/internal/auth"
	"github.com/KelvinJRosado/go-server/internal/database"
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
