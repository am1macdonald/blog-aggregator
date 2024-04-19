-- +goose Up
CREATE TABLE posts (
  id uuid PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  title VARCHAR(255) NOT NULL,
  url TEXT NOT NULL,
  description TEXT NOT NULL,
  published_at TIMESTAMP NOT NULL DEFAULT NOW(),
  feed_id uuid REFERENCES feeds(id)
);

-- +goose Down
DROP TABLE posts;
