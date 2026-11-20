-- +goose Up

CREATE TABLE reminders (
    id SERIAL PRIMARY KEY,
    reminder_id TEXT NOT NULL UNIQUE,
    
    band_id TEXT
        REFERENCES bands(band_id)
        ON DELETE CASCADE,
    user_id TEXT NOT NULL
        REFERENCES users(user_id)
        ON DELETE CASCADE,

    resource_type TEXT NOT NULL CHECK (
        resource_type IN (
            'event',
            'todo_list',
            'todo_item',
            'checklist',
            'checklist_item',
            'song',
            'setlist'
        )
    ),
    resource_id TEXT NOT NULL CHECK (length(trim(resource_id)) > 0),
    body TEXT NOT NULL CHECK (length(trim(body)) > 0),

    remind_at TIMESTAMPTZ NOT NULL,

    notification_id TEXT
        REFERENCES notifications(notification_id)
        ON DELETE SET NULL,
    push_sent_at TIMESTAMPTZ,

    attempt_count INTEGER NOT NULL DEFAULT 0
        CHECK (attempt_count >= 0),
    last_attempt_at TIMESTAMPTZ,
    last_error TEXT,

    claimed_at TIMESTAMPTZ,
    claimed_by TEXT,

    sent_at TIMESTAMPTZ,
    failed_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL REFERENCES users(user_id),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by TEXT REFERENCES users(user_id),

    CONSTRAINT reminders_claim_fields_match CHECK (
        (claimed_at IS NULL) = (claimed_by IS NULL)
    ),
    CONSTRAINT reminders_one_terminal_state CHECK (
        num_nonnulls(sent_at, failed_at, cancelled_at) <= 1
    )
);

CREATE INDEX idx_reminders_due
ON reminders(remind_at, id)
WHERE sent_at IS NULL
    AND failed_at IS NULL
    AND cancelled_at IS NULL;

CREATE INDEX idx_reminders_active_claim
ON reminders(claimed_at)
WHERE sent_at IS NULL
    AND failed_at IS NULL
    AND cancelled_at IS NULL
    AND claimed_at IS NOT NULL;

CREATE INDEX idx_reminders_user_schedule
ON reminders(user_id, remind_at DESC);

CREATE INDEX idx_reminders_resource
ON reminders(resource_type, resource_id);


-- +goose Down

DROP TABLE reminders;
