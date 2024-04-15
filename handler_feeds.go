package main

import (
	"context"
	"net/http"
	"time"

	"github.com/am1macdonald/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleCreateFeed(w http.ResponseWriter, r *http.Request, u database.User) {
	body := struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}{}
	err := DecodeRequest(r, &body)
	if err != nil {
		errorResponse(w, 500, err)
		return
	}
	feed, err := cfg.DB.CreateFeed(context.Background(), database.CreateFeedParams{
		Name:   body.Name,
		Url:    body.Url,
		UserID: u.ID,
	})
	if err != nil {
		errorResponse(w, 500, err)
		return
	}
	jsonResponse(w, 200, struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
		Url       string    `json:"url"`
		UserID    uuid.UUID `json:"user_id"`
	}{
		ID:        feed.ID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		Name:      feed.Name,
		Url:       feed.Url,
		UserID:    feed.UserID})
}
