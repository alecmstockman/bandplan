-- +goose Up
ALTER TABLE band_members

ADD COLUMN edit_band BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN manage_members BOOLEAN NOT NULL DEFAULT false, 

ADD COLUMN create_songs BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN edit_songs BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN delete_songs BOOLEAN NOT NULL DEFAULT false,

ADD COLUMN create_setlists BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN edit_setlists BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN delete_setlists BOOLEAN NOT NULL DEFAULT false,

ADD COLUMN create_events BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN edit_events BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN delete_events BOOLEAN NOT NULL DEFAULT false,

ADD COLUMN create_todos BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN edit_todos BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN delete_todos BOOLEAN NOT NULL DEFAULT false;

UPDATE band_members
SET 
    edit_band = true,
    manage_members = true, 

    create_songs = true,
    edit_songs = true,
    delete_songs = true,

    create_setlists = true,
    edit_setlists = true,
    delete_setlists = true,

    create_events = true,
    edit_events = true,
    delete_events = true,

    create_todos = true,
    edit_todos = true,
    delete_todos = true

WHERE is_admin IS TRUE;

-- +goose Down
ALTER TABLE band_members
DROP COLUMN edit_band,
DROP COLUMN manage_members,
DROP COLUMN create_songs,
DROP COLUMN edit_songs,
DROP COLUMN delete_songs,
DROP COLUMN create_setlists,
DROP COLUMN edit_setlists,
DROP COLUMN delete_setlists,
DROP COLUMN create_events,
DROP COLUMN edit_events,
DROP COLUMN delete_events,
DROP COLUMN create_todos,
DROP COLUMN edit_todos,
DROP COLUMN delete_todos;
