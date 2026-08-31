-- +goose Up
CREATE TABLE message_reactions (
    id SERIAL PRIMARY KEY,

    reaction_id TEXT UNIQUE NOT NULL,

    message_id TEXT NOT NULL
        REFERENCES messages(message_id)
        ON DELETE CASCADE,

    user_id TEXT NOT NULL
        REFERENCES users(user_id)
        ON DELETE CASCADE,

    reaction TEXT NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    UNIQUE (message_id, user_id, reaction)
);

CREATE INDEX idx_message_reactions_message_id
    ON message_reactions(message_id);

CREATE INDEX idx_message_reactions_user_id
    ON message_reactions(user_id);


-- +goose Down
DROP TABLE IF EXISTS message_reactions;
