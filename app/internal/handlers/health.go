package handlers

import "net/http"

func init() {
	Register("GET", "/health", healthHandler)
}

func healthHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return map[string]string{"status": "ok"}, nil
}