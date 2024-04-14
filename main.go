package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var (
	port string
)

type apiConfig struct {
}

func jsonResponse(w http.ResponseWriter, status int, payload interface{}) {
	val, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(val)
}
func errorResponse(w http.ResponseWriter, code int, msg string) {
	jsonResponse(w, code, errors.New(msg))
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Panicln("Failed to load .env")
	}
	port = os.Getenv("PORT")
}

func main() {
	mux := *http.NewServeMux()
	mux.HandleFunc("GET /v1/readiness", HandleGetReadiness)

	corsMux := middlewareCors(&mux)
	server := http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}
	log.Fatal(server.ListenAndServe())
}
