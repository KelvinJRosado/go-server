package server

import (
	"log/slog"
	"net/http"
)

func indexFileHandler() http.Handler {

	slog.Info("Received fileServer request")

	// Add file handler
	currentDir := http.Dir("./internal/server/")
	return http.StripPrefix("/app", http.FileServer(currentDir))
}
