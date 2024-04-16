package main

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/am1macdonald/blog-aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, *database.User)

func (cfg *apiConfig) middlewareAuth(next authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		next(w, r, &user)
	}
}
