-- +goose Up
CREATE TABLE IF NOT EXISTS receptions
(
    id        UUID PRIMARY KEY,
    date_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    pvz_id    UUID        NOT NULL REFERENCES pvzs (id) ON DELETE CASCADE,
    status    TEXT        NOT NULL CHECK (status IN ('in_progress', 'close'))
);

-- +goose Down
DROP TABLE IF EXISTS receptions;