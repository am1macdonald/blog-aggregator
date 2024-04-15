package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/am1macdonald/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Name string `json:"name"`
	}{}
	err := DecodeRequest(r, &body)
	if err != nil {
		errorResponse(w, 500, err)
		return
	}
	user, err := cfg.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      body.Name,
	})
	log.Println(user)
	if err != nil {
		errorResponse(w, 500, err)
		return
	}
	jsonResponse(w, 201, user)
}

func (cfg *apiConfig) handleGetUserByApiKey(w http.ResponseWriter, r *http.Request) {
	t := r.Header.Get("Authorization")
	if t == "" || strings.Split(t, " ")[0] != "ApiKey" {
		errorResponse(w, 404, errors.New("api key required"))
		return
	}
	user, err := cfg.DB.GetUserByApiKey(context.Background(), strings.Split(t, " ")[1])
	if err != nil {
		errorResponse(w, 500, err)
		return
	}
	jsonResponse(w, 200, user)
	return
}
