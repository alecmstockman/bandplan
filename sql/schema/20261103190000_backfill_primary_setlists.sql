-- +goose Up

WITH first_setlist AS (
	SELECT DISTINCT ON (band_id)
		setlist_id
	FROM setlists s
	WHERE NOT EXISTS (
		SELECT 1
		FROM setlists primary_setlist
		WHERE primary_setlist.band_id = s.band_id
			AND primary_setlist.is_primary = TRUE
	)
	ORDER BY band_id, created_at, id
)
UPDATE setlists s
SET is_primary = TRUE
FROM first_setlist
WHERE s.setlist_id = first_setlist.setlist_id;

-- +goose Down

-- This data backfill is intentionally irreversible. A later primary-setlist
-- choice must not be discarded when rolling back application migrations.
SELECT 1;
