package main

import "net/http"

func HandleGetReadiness(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, 200, struct {
		Status string `json:"status"`
	}{Status: "ok"})
	return
}
