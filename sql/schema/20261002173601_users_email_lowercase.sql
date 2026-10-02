-- +goose Up

UPDATE users
SET email = LOWER(email);

ALTER TABLE users
    ADD CONSTRAINT users_email_lowercase
    CHECK (email = LOWER(email));


-- +goose Down

ALTER TABLE users
    DROP CONSTRAINT users_email_lowercase;