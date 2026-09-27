-- +goose Up
ALTER TABLE users
    ADD COLUMN first_name TEXT,
    ADD COLUMN last_name TEXT;

UPDATE users
SET first_name = name 
WHERE first_name IS NULL;

ALTER TABLE users
    ALTER COLUMN first_name SET NOT NULL;


-- +goose Down
ALTER TABLE users
    DROP COLUMN first_name
    DROP COLUMN last_name;
