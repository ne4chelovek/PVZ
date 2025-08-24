-- +goose Up
CREATE TABLE IF NOT EXISTS products
(
    id           UUID PRIMARY KEY,
    date_time    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    type         TEXT        NOT NULL CHECK (type IN ('электроника', 'одежда', 'обувь')),
    reception_id UUID        NOT NULL REFERENCES receptions (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS products;