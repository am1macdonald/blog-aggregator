-- +goose Up
CREATE TABLE feed_follows (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  user_id uuid REFERENCES users (id) NOT NULL,
  feed_id uuid REFERENCES feeds (id) NOT NULL
);

-- +goose Down
DROP TABLE feed_follows;
