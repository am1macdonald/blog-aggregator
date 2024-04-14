package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/am1macdonald/blog-aggregator/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	port string
)

type apiConfig struct {
	DB *database.Queries
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

func errorResponse(w http.ResponseWriter, code int, msg interface{}) {
	log.Println(msg)
	jsonResponse(w, code, msg)
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
	db, err := sql.Open("postgres", os.Getenv("DB_CONN"))
	if err != nil {
		log.Panicf("Fatal: %v", err)
	}
	cfg := apiConfig{
		DB: database.New(db),
	}

	mux := *http.NewServeMux()

	// testing
	mux.HandleFunc("GET /v1/readiness", cfg.HandleGetReadiness)
	mux.HandleFunc("GET /v1/error", cfg.HandleGetError)

	corsMux := middlewareCors(&mux)
	server := http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}
	log.Fatal(server.ListenAndServe())
}
