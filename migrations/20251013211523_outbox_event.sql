-- +goose Up
-- +goose StatementBegin
CREATE TABLE outbox_events (
    id UUID PRIMARY KEY,
    aggregate_id TEXT NOT NULL,
    event_type TEXT NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    processed BOOLEAN DEFAULT FALSE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS outbox_events;
-- +goose StatementEnd
