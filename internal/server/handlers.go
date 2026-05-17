package server

import "net/http"

func indexFileHandler() http.Handler {
	// Add file handler
	currentDir := http.Dir("./internal/server/")
	return http.FileServer(currentDir)
}
