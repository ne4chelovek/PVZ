-- +goose Up
CREATE TABLE IF NOT EXISTS pvzs
(
    id                UUID PRIMARY KEY,
    registration_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    city              TEXT        NOT NULL CHECK (city IN ('Москва', 'Санкт-Петербург', 'Казань'))
);

-- +goose Down
DROP TABLE IF EXISTS pvzs;