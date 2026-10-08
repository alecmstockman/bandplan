-- +goose Up

INSERT INTO chat_members (chat_id, user_id)
SELECT chats.chat_id, band_members.user_id
FROM band_members
JOIN chats
    ON chats.band_id = band_members.band_id
    AND chats.is_primary = TRUE
ON CONFLICT (chat_id, user_id) DO NOTHING;

-- +goose Down

-- Backfilled memberships are indistinguishable from memberships created normally.
SELECT 1;
