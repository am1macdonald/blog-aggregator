-- +goose Up
CREATE TABLE feeds (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT Now(),
  updated_at TIMESTAMP NOT NULL DEFAULT Now(),
  name TEXT NOT NULL,
  url TEXT NOT NULL,
  user_id UUID REFERENCES users (id)
);

-- +goose Down
DROP TABLE feeds;
