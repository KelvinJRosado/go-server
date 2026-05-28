package server

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

func (ac *apiConfig) upgradeUserToRedHandler() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		// Get input
		type input struct {
			Event string `json:"event"`
			Data  struct {
				UserId uuid.UUID `json:"user_id"`
			} `json:"data"`
		}

		params, ok := getInputStruct[input](res, req)
		if !ok {
			// Util handles writing error response, so we can just return
			return
		}

		// Ignore non-handled events
		if params.Event != "user.upgraded" {
			slog.Info("ignoring unhandled event", "event", params.Event)
			respondWithJSON(res, http.StatusNoContent, nil)
			return
		}

		// Check if user exists
		user, err := ac.databaseQueries.GetUserById(req.Context(), params.Data.UserId)
		if err != nil {
			slog.Error("failed to get user", "error", err)
			respondWithError(res, http.StatusNotFound, "User Not Found")
			return
		}

		err = ac.databaseQueries.UpgradeUserToRed(req.Context(), user.ID)
		if err != nil {
			slog.Error("failed to upgrade user", "error", err)
			respondWithInternalError(res)
			return
		}

		respondWithJSON(res, http.StatusNoContent, nil)
	})
}
