-- +goose Up

DELETE FROM chat_members cm
USING chats c
WHERE cm.chat_id = c.chat_id
    AND NOT EXISTS (
        SELECT 1
        FROM band_members bm
        WHERE bm.band_id = c.band_id
            AND bm.user_id = cm.user_id
    );

-- +goose Down

-- Removed chat memberships cannot be restored safely.
SELECT 1;
