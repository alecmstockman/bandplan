-- +goose Up

ALTER TABLE messages
ADD COLUMN chat_id TEXT NOT NULL
REFERENCES chats(chat_id) ON DELETE CASCADE;

CREATE INDEX idx_messages_chat_id
ON messages(chat_id);


-- +goose Down

DROP INDEX idx_messages_chat_id;

ALTER TABLE messages
DROP COLUMN chat_id;