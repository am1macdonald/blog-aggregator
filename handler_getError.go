package main

import "net/http"

func (cfg *apiConfig) HandleGetError(w http.ResponseWriter, r *http.Request) {
	errorResponse(w, 500, struct {
		Error string `json:"error"`
	}{Error: "internal server error"})
	return
}
