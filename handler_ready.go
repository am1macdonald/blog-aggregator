package main

import "net/http"

func (cfg *apiConfig) handleGetReadiness(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, 200, struct {
		Status string `json:"status"`
	}{Status: "ok"})
	return
}

func (cfg *apiConfig) handleGetError(w http.ResponseWriter, r *http.Request) {
	errorResponse(w, 500, struct {
		Error string `json:"error"`
	}{Error: "internal server error"})
	return
}
