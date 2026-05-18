package server

import (
	"net/http"
)

func (ac *apiConfig) indexFileHandler() http.Handler {
	// Add file handler
	currentDir := http.Dir("./internal/server/")
	return http.StripPrefix("/app", http.FileServer(currentDir))
}
