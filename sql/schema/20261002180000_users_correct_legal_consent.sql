-- +goose Up

ALTER TABLE users
    ADD COLUMN terms_version TEXT,
    ADD COLUMN privacy_version TEXT,
    ALTER COLUMN legal_accepted_at DROP DEFAULT,
    ALTER COLUMN legal_accepted_at DROP NOT NULL;

UPDATE users
SET
    legal_accepted_at = NULL,
    terms_version = NULL,
    privacy_version = NULL
WHERE NOT legal_accepted;

UPDATE users
SET
    terms_version = '2026-07-14',
    privacy_version = '2026-07-14'
WHERE legal_accepted;

ALTER TABLE users
    ADD CONSTRAINT users_legal_consent_consistent CHECK (
        (
            legal_accepted
            AND legal_accepted_at IS NOT NULL
            AND NULLIF(BTRIM(terms_version), '') IS NOT NULL
            AND NULLIF(BTRIM(privacy_version), '') IS NOT NULL
        )
        OR
        (
            NOT legal_accepted
            AND legal_accepted_at IS NULL
            AND terms_version IS NULL
            AND privacy_version IS NULL
        )
    );

-- +goose Down

ALTER TABLE users
    DROP CONSTRAINT users_legal_consent_consistent,
    DROP COLUMN terms_version,
    DROP COLUMN privacy_version;

UPDATE users
SET legal_accepted_at = CURRENT_TIMESTAMP
WHERE legal_accepted_at IS NULL;

ALTER TABLE users
    ALTER COLUMN legal_accepted_at SET DEFAULT CURRENT_TIMESTAMP,
    ALTER COLUMN legal_accepted_at SET NOT NULL;
