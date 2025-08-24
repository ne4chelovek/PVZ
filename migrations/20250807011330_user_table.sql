-- +goose Up
CREATE TABLE IF NOT EXISTS users
(
    id       UUID PRIMARY KEY,
    email    TEXT UNIQUE NOT NULL,
    password TEXT        NOT NULL,
    role     TEXT        NOT NULL CHECK (role IN ('employee', 'moderator'))
);

-- +goose Down
DROP TABLE IF EXISTS users;