package server

import (
	"log/slog"
	"net/http"
)

func (ac *apiConfig) createUserHandler() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		// Get input
		type input struct {
			Email string `json:"email"`
		}

		params, ok := getInputStruct[input](res, req)
		if !ok {
			// Util handles writing error response, so we can just return
			return
		}

		user, err := ac.databaseQueries.CreateUser(req.Context(), params.Email)
		if err != nil {
			slog.Error("failed to create user", "error", err)
			respondWithInternalError(res)
			return
		}

		respondWithJSON(res, http.StatusCreated, user)
	})
}
