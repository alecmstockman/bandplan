-- +goose Up

ALTER TABLE events
ADD COLUMN todo_list_id TEXT
    REFERENCES todo_lists(todo_list_id)
    ON DELETE SET NULL,
ADD COLUMN checklist_id TEXT
    REFERENCES checklist_lists(checklist_id)
    ON DELETE SET NULL;

CREATE UNIQUE INDEX idx_events_todo_list_id
ON events(todo_list_id)
WHERE todo_list_id IS NOT NULL;

CREATE UNIQUE INDEX idx_events_checklist_id
ON events(checklist_id)
WHERE checklist_id IS NOT NULL;

ALTER TABLE songs
ADD COLUMN todo_list_id TEXT
    REFERENCES todo_lists(todo_list_id)
    ON DELETE SET NULL,
ADD COLUMN checklist_id TEXT
    REFERENCES checklist_lists(checklist_id)
    ON DELETE SET NULL;

CREATE UNIQUE INDEX idx_songs_todo_list_id
ON songs(todo_list_id)
WHERE todo_list_id IS NOT NULL;

CREATE UNIQUE INDEX idx_songs_checklist_id
ON songs(checklist_id)
WHERE checklist_id IS NOT NULL;

ALTER TABLE setlists
ADD COLUMN todo_list_id TEXT
    REFERENCES todo_lists(todo_list_id)
    ON DELETE SET NULL,
ADD COLUMN checklist_id TEXT
    REFERENCES checklist_lists(checklist_id)
    ON DELETE SET NULL;

CREATE UNIQUE INDEX idx_setlists_todo_list_id
ON setlists(todo_list_id)
WHERE todo_list_id IS NOT NULL;

CREATE UNIQUE INDEX idx_setlists_checklist_id
ON setlists(checklist_id)
WHERE checklist_id IS NOT NULL;

ALTER TABLE todo_items
ADD COLUMN checklist_id TEXT
    REFERENCES checklist_lists(checklist_id)
    ON DELETE SET NULL;

CREATE UNIQUE INDEX idx_todo_items_checklist_id
ON todo_items(checklist_id)
WHERE checklist_id IS NOT NULL;


-- +goose Down

DROP INDEX idx_todo_items_checklist_id;

ALTER TABLE todo_items
DROP COLUMN checklist_id;

DROP INDEX idx_setlists_checklist_id;
DROP INDEX idx_setlists_todo_list_id;

ALTER TABLE setlists
DROP COLUMN checklist_id,
DROP COLUMN todo_list_id;

DROP INDEX idx_songs_checklist_id;
DROP INDEX idx_songs_todo_list_id;

ALTER TABLE songs
DROP COLUMN checklist_id,
DROP COLUMN todo_list_id;

DROP INDEX idx_events_checklist_id;
DROP INDEX idx_events_todo_list_id;

ALTER TABLE events
DROP COLUMN checklist_id,
DROP COLUMN todo_list_id;
