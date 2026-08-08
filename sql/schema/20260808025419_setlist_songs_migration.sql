-- +goose Up

ALTER TABLE setlist_songs
RENAME TO setlist_items;


-- +goose Down

ALTER TABLE setlist_items
RENAME TO setlist_songs;
