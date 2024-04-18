package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/am1macdonald/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

func databaseFeedToFeed(feed database.Feed) Feed {
	return Feed{
		ID:        feed.ID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		Name:      feed.Name,
		Url:       feed.Url,
		UserID:    feed.UserID,
	}
}

func (cfg *apiConfig) handleCreateFeed(w http.ResponseWriter, r *http.Request, u *database.User) {
	body := struct {
		Name string `json:"name,omitempty"`
		Url  string `json:"url,omitempty"`
	}{}
	err := DecodeRequest(r, &body)
	if err != nil {
		errorResponse(w, 500, err)
		return
	}
	fmt.Println(body)
	if body.Name == "" || body.Url == "" {
		errorResponse(w, 404, errors.New("invalid request"))
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
	ff, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		UserID: u.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		errorResponse(w, 500, err)
		return
	}
	jsonResponse(w, 200, struct {
		FeedResponse database.Feed       `json:"feed"`
		FeedFollow   database.FeedFollow `json:"feed_follow"`
	}{
		FeedResponse: feed,
		FeedFollow:   ff,
	},
	)
}

func (cfg *apiConfig) handleGetFeeds(w http.ResponseWriter, r *http.Request) {
	f, err := cfg.DB.GetAllFeeds(r.Context())
	if err != nil {
		errorResponse(w, 500, err)
		return
	}
	feeds := []Feed{}
	for _, feed := range f {
		feeds = append(feeds, databaseFeedToFeed(feed))
	}
	jsonResponse(w, 200, struct {
		Feeds []Feed `json:"feeds"`
	}{Feeds: feeds})
}
