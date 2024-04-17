-- +goose Up
ALTER TABLE feed_follows
ADD last_fetched_at TIMESTAMP null; 

-- +goose Down
ALTER TABLE feed_follows
DROP COLUMN last_fetched_at; 
