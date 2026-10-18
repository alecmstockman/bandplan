-- +goose Up
ALTER TABLE events
ADD COLUMN end_date DATE;

UPDATE events
SET end_date = (end_time AT TIME ZONE time_zone)::DATE
WHERE end_time IS NOT NULL;

ALTER TABLE events
ADD CONSTRAINT events_end_date_not_before_event_date
CHECK (end_date IS NULL OR end_date >= event_date);

-- +goose Down
ALTER TABLE events
DROP CONSTRAINT IF EXISTS events_end_date_not_before_event_date;

ALTER TABLE events
DROP COLUMN IF EXISTS end_date;
