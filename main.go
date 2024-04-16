package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

type apiConfig struct {
	DB *database.Queries
}

func DecodeRequest[T any](r *http.Request, dest *T) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&dest)
	if err != nil {
		log.Printf("Error decoding request: %s", err)
		return err
	}
	return nil
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
	mux.HandleFunc("GET /v1/readiness", cfg.handleGetReadiness)
	mux.HandleFunc("GET /v1/error", cfg.handleGetError)

	// users
	mux.HandleFunc("POST /v1/users", cfg.handleCreateUser)
	mux.HandleFunc("GET /v1/users", cfg.handleGetUserByApiKey)

	// feeds
	mux.HandleFunc("POST /v1/feeds", cfg.middlewareAuth(cfg.handleCreateFeed))
	mux.HandleFunc("GET /v1/feeds", cfg.handleGetFeeds)
	mux.HandleFunc("POST /v1/feed_follows", cfg.middlewareAuth(cfg.handleFollowFeed))

	corsMux := middlewareCors(&mux)
	server := http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}
	fmt.Printf("Server listening at http://localhost:%v", port)
	log.Fatal(server.ListenAndServe())
}
