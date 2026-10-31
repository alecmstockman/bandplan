-- +goose Up
ALTER TABLE band_members
ADD COLUMN create_songs BOOLEAN DEFAULT false,
ADD COLUMN edit_songs BOOLEAN DEFAULT false,
ADD COLUMN delete_songs BOOLEAN DEFAULT false,

ADD COLUMN create_setlists BOOLEAN DEFAULT false,
ADD COLUMN edit_setlists BOOLEAN DEFAULT false,
ADD COLUMN delete_setlists BOOLEAN DEFAULT false,

ADD COLUMN create_events BOOLEAN DEFAULT false,
ADD COLUMN edit_events BOOLEAN DEFAULT false,
ADD COLUMN delete_events BOOLEAN DEFAULT false,

ADD COLUMN create_todos BOOLEAN DEFAULT false,
ADD COLUMN edit_todo BOOLEAN DEFAULT false,
ADD COLUMN delete_todos BOOLEAN DEFAULT false;

UPDATE band_members
SET 
    create_song = true,
    edit_song = true,
    delete_song = true,

    create_setlists = true,
    edit_setlists = true,
    delete_setlists = true,

    create_events = true,
    edit_events = true,
    delete_events = true,

    create_todos = true,
    edit_todos = true,
    delete_todos = true, 

WHERE is_admin IS TRUE;

-- +goose Down
ALTER TABLE band_members
DROP COLUMN
