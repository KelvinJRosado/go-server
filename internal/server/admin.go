package server

import (
	"fmt"
	"log/slog"
	"net/http"
)

func (ac *apiConfig) metricsHandler() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Update response header
		res.Header().Add("Content-Type", "text/html; charset=utf-8")

		// Update response body and status
		htmlBody := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, ac.fileserverHits.Load())

		res.WriteHeader(http.StatusOK)
		res.Write(fmt.Appendf(nil, "%s", htmlBody))
	})
}

func (ac *apiConfig) adminResetHandler() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		// Block if not dev
		if ac.platform != "dev" {
			respondWithError(res, http.StatusForbidden, "Reset only available in dev")
			return
		}

		// Delete all users
		err := ac.databaseQueries.DeleteAllUsers(req.Context())
		if err != nil {
			slog.Error("failed to delete users", "error", err)
			respondWithInternalError(res)
			return
		}

		// Reset counter back to 0
		ac.fileserverHits.Store(0)

		// Update response header
		res.Header().Add("Content-Type", "text/plain; charset=utf-8")

		// Update response body and status
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("Counter was reset"))
	})
}
