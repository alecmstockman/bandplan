-- +goose Up
ALTER TABLE setlists
ADD COLUMN is_primary BOOLEAN NOT NULL DEFAULT FALSE;

CREATE UNIQUE INDEX idx_setlists_one_primary_per_band
ON setlists (band_id)
WHERE is_primary = TRUE;

-- +goose Down
DROP INDEX idx_setlists_one_primary_per_band;

ALTER TABLE setlists
DROP COLUMN is_primary;
