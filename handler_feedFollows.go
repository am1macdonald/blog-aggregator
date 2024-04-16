package main

import (
	"errors"
	"net/http"

	"github.com/am1macdonald/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleFollowFeed(w http.ResponseWriter, r *http.Request, u *database.User) {
	body := struct {
		FeedID uuid.UUID `json:"feed_id,omitempty"`
	}{}
	err := DecodeRequest(r, &body)
	if err != nil {
		errorResponse(w, 500, err)
		return
	}
	if body.FeedID == uuid.Nil {
		errorResponse(w, 404, errors.New("invalid request"))
		return
	}
	f, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		UserID: u.ID,
		FeedID: body.FeedID,
	})
	if err != nil {
		errorResponse(w, 500, err)
		return
	}
	jsonResponse(w, 200, f)
}

func (cfg *apiConfig) handleDeleteFeedFollow(w http.ResponseWriter, r *http.Request, u *database.User) {
	id := r.PathValue("feed_follow_id")
	if id == "" {
		errorResponse(w, 404, errors.New("missing path parameter"))
		return
	}
	uuid, err := uuid.Parse(id)
	if err != nil {
		errorResponse(w, 500, err)
		return
	}
	err = cfg.DB.DeleteFeedFollow(r.Context(), uuid)
	if err != nil {
		errorResponse(w, 500, err)
		return
	}
	jsonResponse(w, 200, struct{}{})
}
