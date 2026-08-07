-- +goose Up

ALTER TABLE transitions
ADD COLUMN link_one TEXT,
ADD COLUMN link_two TEXT,
ADD COLUMN link_three TEXT;


-- +goose Down

ALTER TABLE transitions
DROP COLUMN link_one,
DROP COLUMN link_two,
DROP COLUMN link_three;