-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS events (
    id uuid PRIMARY KEY,
    event_type VARCHAR(255) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    sent_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_events_created_at ON events (created_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_events_created_at;
DROP TABLE IF EXISTS events;
-- +goose StatementEnd
