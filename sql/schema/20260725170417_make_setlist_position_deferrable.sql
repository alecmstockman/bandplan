-- +goose Up

ALTER TABLE setlist_songs
DROP CONSTRAINT setlist_songs_setlist_id_position_key;

ALTER TABLE setlist_songs
ADD CONSTRAINT setlist_songs_position_unique
UNIQUE (setlist_id, position)
DEFERRABLE INITIALLY IMMEDIATE;

-- +goose Down
ALTER TABLE setlist_songs
DROP CONSTRAINT setlist_songs_position_unique;

ALTER TABLE setlist_songs
ADD CONSTRAINT setlist_songs_setlist_id_position_key
UNIQUE (setlist_id, position);
