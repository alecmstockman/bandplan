-- +goose Up
CREATE TABLE todo_list_assignees(
    id SERIAL PRIMARY KEY,
    todo_list_assignee_id TEXT NOT NULL UNIQUE,
    user_id TEXT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    band_id TEXT NOT NULL REFERENCES bands(band_id) ON DELETE CASCADE,
    todo_list_id TEXT NOT NULL REFERENCES todo_lists(todo_list_id) ON DELETE CASCADE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL REFERENCES users(user_id),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by TEXT NOT NULL REFERENCES users(user_id),

    UNIQUE (todo_list_id, user_id),

    FOREIGN KEY (band_id, user_id)
    REFERENCES band_members(band_id, user_id)
    ON DELETE CASCADE
);

CREATE INDEX idx_todo_list_assignees_to_band_id
    ON todo_list_assignees (band_id);

CREATE INDEX idx_todo_list_assignees_to_user_id
    ON todo_list_assignees (user_id);

CREATE TABLE todo_item_assignees(
    id SERIAL PRIMARY KEY,
    todo_item_assignee_id TEXT NOT NULL UNIQUE,
    user_id TEXT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    band_id TEXT NOT NULL REFERENCES bands(band_id) ON DELETE CASCADE,
    item_id TEXT NOT NULL REFERENCES todo_items(item_id) ON DELETE CASCADE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL REFERENCES users(user_id),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by TEXT NOT NULL REFERENCES users(user_id),

    UNIQUE (item_id, user_id),

    FOREIGN KEY (band_id, user_id)
    REFERENCES band_members(band_id, user_id)
    ON DELETE CASCADE
);

CREATE INDEX idx_todo_item_assignees_to_band_id
    ON todo_item_assignees (band_id);

CREATE INDEX idx_todo_item_assignees_to_user_id
    ON todo_item_assignees (user_id);

CREATE TABLE checklist_assignees(
    id SERIAL PRIMARY KEY,
    checklist_assignee_id TEXT NOT NULL UNIQUE,
    user_id TEXT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    band_id TEXT NOT NULL REFERENCES bands(band_id) ON DELETE CASCADE,
    checklist_id TEXT NOT NULL REFERENCES checklist_lists(checklist_id) ON DELETE CASCADE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL REFERENCES users(user_id),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by TEXT NOT NULL REFERENCES users(user_id),

    UNIQUE (checklist_id, user_id),
    
    FOREIGN KEY (band_id, user_id)
    REFERENCES band_members(band_id, user_id)
    ON DELETE CASCADE
);

CREATE INDEX idx_checklist_assignees_to_band_id
    ON checklist_assignees (band_id);

CREATE INDEX idx_checklist_assignees_to_user_id
    ON checklist_assignees (user_id);

-- +goose Down
DROP TABLE todo_list_assignees;
DROP TABLE todo_item_assignees;
DROP TABLE checklist_assignees;
