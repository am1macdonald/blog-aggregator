package main

import (
	"net/http"
	"time"

	"github.com/am1macdonald/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

type Post struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	FeedID      uuid.UUID `json:"feed_id"`
}

func databasePostToPost(post database.Post) Post {
	return Post{
		ID:          post.ID,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		Title:       post.Title,
		Url:         post.Url,
		Description: post.Description,
		PublishedAt: post.PublishedAt,
		FeedID:      post.FeedID,
	}
}

func (cfg *apiConfig) handleGetPostsByUser(w http.ResponseWriter, r *http.Request, u *database.User) {
	p, err := cfg.DB.GetPostsByUser(r.Context(), u.ID)
	if err != nil {
		errorResponse(w, 500, err)
		return
	}
	posts := make([]Post, len(p))
	for _, post := range p {
		posts = append(posts, databasePostToPost(post))
	}
	jsonResponse(w, 200, struct {
		Posts []Post `json:"posts"`
	}{Posts: posts})
}
